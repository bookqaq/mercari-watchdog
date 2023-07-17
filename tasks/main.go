package tasks

import (
	"fmt"
	"net/http"
	"time"

	"bookq.xyz/mercari-watchdog/bot"
	"bookq.xyz/mercari-watchdog/compare"
	"bookq.xyz/mercari-watchdog/models/analysisdata"
	"bookq.xyz/mercari-watchdog/models/analysistask"
	"bookq.xyz/mercari-watchdog/models/fetchdata"
	"bookq.xyz/mercari-watchdog/tools"

	"github.com/google/uuid"

	"github.com/bookqaq/mer-wrapper/common"
	wrapperv2 "github.com/bookqaq/mer-wrapper/v2"
)

const (
	TaskRoutines = 10
)

var taskChans []chan analysistask.AnalysisTask

func Boot() {
	tools.RefreshBlockedSellers()
	analysisdata.RenewAll()
	go analysistask.AddTaskBuffer()

	// create tickers for time-based tasks
	ticker_1h := time.NewTicker(1800 * time.Second)
	ticker_10m := time.NewTicker(60 * time.Second)
	ticker_5m := time.NewTicker(60 * time.Second)
	ticker_clearExpiredFetch := time.NewTicker(150 * time.Second)
	ticker_cleanImages := time.NewTicker(7200 * time.Second)

	// manage all workers in a slice
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
			go tools.RefreshBlockedSellers()
		case <-ticker_cleanImages.C:
			go cleanImages()
		}
	}
}

// tasks worker
func taskChanListener(taskInput <-chan analysistask.AnalysisTask) {
	for {
		task := <-taskInput
		taskModifier(&task)
		runTask(time.Now(), task)
	}
}

func taskModifier(task *analysistask.AnalysisTask) {
	task.MaxPage = 1
}

func runWorkflow(interval int, t time.Time) {
	// for dev locally
	//proxyUrl := "http://127.0.0.1:8889"
	//proxy, _ := url.Parse(proxyUrl)
	//tr := &http.Transport{
	//	Proxy:           http.ProxyURL(proxy),
	//	TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	//}
	//common.Client.Content = &http.Client{
	//	Transport: tr,
	//}

	common.Client.Content = http.DefaultClient

	common.Client.ClientID = uuid.NewString()

	taskResults, err := analysistask.GetAll(interval)
	if err != nil {
		fmt.Printf("error during processing workflow %s : %v", t, interval)
		return
	}

	// TODO: convert taskChans to link list with loop
	for i, taskItem := range taskResults {
		//fmt.Println("Started: ", i%TaskRoutines, taskItem.TaskID, taskItem.Keywords)
		taskChans[i%TaskRoutines] <- taskItem
	}
}

func runTask(t time.Time, task analysistask.AnalysisTask) {

	// fetch items data from mercari
	data, err := wrapperv2.Search(wrapperv2.SearchData{
		Keyword: tools.ConcatKeyword(task.Keywords),
		Limit:   30,
	})
	if err != nil {
		fmt.Printf("failed to search, taskID %v, time %v\n", task.TaskID, t.Unix())
		return
	}

	// get AnalysisData to generate message
	recentItems, err := analysisdata.GetOne(task.TaskID)
	if err != nil {
		fmt.Printf("failed to get last search data, taskID %v, time %v, %s\n", task.TaskID, t.Unix(), err)
		return
	}
	var result []wrapperv2.MercariV2Item

	// mainly v3, implement compatability about v2
	if len(task.MustMatch) <= 0 {
		result, err = compare.Run2(data, recentItems, task)
	} else {
		result, err = compare.Run3(data, recentItems, task)
	}

	if err != nil {
		fmt.Printf("failed to compare, taskID %v, time %v, %s\n", task.TaskID, t.Unix(), err)
		return
	}

	// update AnalysisData for expansion
	recentItems.Data = result
	recentItems.Time = time.Now().Unix()
	recentItems.Length = len(result)

	go bot.MercariPushMsg(recentItems, task.Owner, task.Group)

	err = analysisdata.Update(recentItems)
	if err != nil {
		fmt.Printf("failed to update AnalysisData, taskID %v, time %v, %s", recentItems.TaskID, t.Unix(), err)
		return
	}
	//fmt.Println("Pushed: ", task.TaskID, task.Keywords)
}
