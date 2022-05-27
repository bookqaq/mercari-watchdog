package tasks

import (
	"fmt"
	"time"

	"bookq.xyz/mercari-watchdog/bot"
	"bookq.xyz/mercari-watchdog/compare"
	"bookq.xyz/mercari-watchdog/utils"
	"bookq.xyz/mercari-watchdog/utils/analysisdata"
	"bookq.xyz/mercari-watchdog/utils/analysistask"
	"bookq.xyz/mercari-watchdog/utils/fetchdata"
	"github.com/bookqaq/goForMercari/mercarigo"
	merwrapper "github.com/bookqaq/mer-wrapper"
	"github.com/google/uuid"
)

const (
	TaskRoutines = 5
)

var taskChans []chan analysistask.AnalysisTask

func Boot() {
	utils.Init()

	ticker_10m := time.NewTicker(600 * time.Second)
	ticker_1h := time.NewTicker(3600 * time.Second)
	ticker_5m := time.NewTicker(300 * time.Second)
	ticker_clearExpiredFetch := time.NewTicker(150 * time.Second)

	taskChans = make([]chan analysistask.AnalysisTask, TaskRoutines)
	for i := 0; i < TaskRoutines; i++ {
		taskChans[i] = make(chan analysistask.AnalysisTask, 5)
		go taskChanListener(taskChans[i])
	}

	// Run tasks when Ticker tick.
	for {
		select {
		case t := <-ticker_1h.C:
			go runWorkflow(3600, t)
		case t := <-ticker_10m.C:
			go runWorkflow(600, t)
		case t := <-ticker_5m.C:
			go runWorkflow(300, t)
		case <-ticker_clearExpiredFetch.C:
			go fetchdata.ClearExpired()
		}
	}
}

func taskChanListener(taskInput <-chan analysistask.AnalysisTask) {
	for {
		task := <-taskInput
		runTask(time.Now(), task)
		time.Sleep(6 * time.Second)
	}
}

func runWorkflow(interval int, t time.Time) {
	//proxyUrl := "http://127.0.0.1:12355"
	//proxy, _ := url.Parse(proxyUrl)
	//tr := &http.Transport{
	//	Proxy:           http.ProxyURL(proxy),
	//	TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	//}
	//merwrapper.Client.Content = &http.Client{
	//	Transport: tr,
	//}

	merwrapper.Client.ClientID = uuid.NewString()
	taskResults, err := analysistask.GetAll(interval)
	if err != nil {
		fmt.Printf("error during processing workflow %s : %v", t, interval)
		return
	}

	for i, taskItem := range taskResults {
		taskChans[i%5] <- taskItem
	}
}

func runTask(t time.Time, task analysistask.AnalysisTask) {
	//fmt.Printf("debug: task %v run\n", task.TaskID)
	data, err := mercarigo.Mercari_search(utils.ConcatKeyword(task.Keywords), task.Sort, task.Order, "on_sale", 30, task.MaxPage)
	if err != nil {
		fmt.Printf("failed to search, taskID %v, time %v\n", task.TaskID, t.Unix())
		return
	}

	recentItems, err := analysisdata.GetOne(task.TaskID)
	if err != nil {
		fmt.Printf("failed to get last search data, taskID %v, time %v, %s\n", task.TaskID, t.Unix(), err)
		return
	}
	result, err := compare.Run3(data, recentItems, task)
	if err != nil {
		fmt.Printf("failed to compare, taskID %v, time %v, %s\n", task.TaskID, t.Unix(), err)
		return
	}

	//if !reflect.DeepEqual(recentItems.Keywords, task.Keywords) {
	//	fmt.Printf("Found keyword error in %d, %v -> %v", recentItems.TaskID, task.Keywords, recentItems.Keywords)
	//	recentItems.Keywords = task.Keywords
	//}

	recentItems.Data = result
	recentItems.Time = time.Now().Unix()
	recentItems.Length = len(result)

	go bot.MercariPushMsg(recentItems, task.Owner, task.Group)

	err = analysisdata.Update(recentItems)
	if err != nil {
		fmt.Printf("failed to update AnalysisData, taskID %v, time %v, %s", recentItems.TaskID, t.Unix(), err)
		return
	}
}
