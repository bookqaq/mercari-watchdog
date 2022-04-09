package tasks

import (
	"fmt"
	"time"

	"bookq.xyz/mercariWatchdog/bot"
	"bookq.xyz/mercariWatchdog/compare"
	"bookq.xyz/mercariWatchdog/utils"
	"github.com/bookqaq/goForMercari/mercarigo"
	merwrapper "github.com/bookqaq/mer-wrapper"
	"github.com/google/uuid"
)

const (
	TaskRoutines = 5
)

var taskChans []chan utils.AnalysisTask

func Boot() {
	utils.Init()

	debug_ticker := time.NewTicker(18400 * time.Second)
	ticker_10m := time.NewTicker(600 * time.Second)

	taskChans = make([]chan utils.AnalysisTask, TaskRoutines)
	for i := 0; i < TaskRoutines; i++ {
		taskChans[i] = make(chan utils.AnalysisTask, 5)
		go taskChanListener(taskChans[i])
	}

	tickCounter := 0
	maxCounter := false

	for {
		select {
		case t := <-ticker_10m.C:
			tickCounter++
			go runWorkflow(600, t)
			if 1 <= (tickCounter / 6) {
				go runWorkflow(3600, t)
				maxCounter = true
			}
			if maxCounter {
				tickCounter = 0
				maxCounter = false
			}
		case t := <-debug_ticker.C:
			fmt.Printf("占位计时器. %v\n", t.Unix())
		}
	}
}

func taskChanListener(taskInput <-chan utils.AnalysisTask) {
	for {
		task := <-taskInput
		runTask(time.Now(), task)
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
	taskResults, err := utils.GetAllTasks(interval)
	if err != nil {
		fmt.Printf("error during processing workflow %s : %v", t, interval)
		return
	}

	for i, taskItem := range taskResults {
		taskChans[i%5] <- taskItem
	}
}

func runTask(t time.Time, task utils.AnalysisTask) {
	//fmt.Printf("debug: task %v run\n", task.TaskID)
	data, err := mercarigo.Mercari_search(utils.ConcatKeyword(task.Keywords), task.Sort, task.Order, "on_sale", 30, task.MaxPage)
	if err != nil {
		fmt.Printf("failed to search, taskID %v, time %v\n", task.TaskID, t.Unix())
		return
	}

	recentItems, err := utils.GetDataDB(task.TaskID)
	if err != nil {
		fmt.Printf("failed to get last search data, taskID %v, time %v, %s\n", task.TaskID, t.Unix(), err)
		return
	}
	result, err := compare.Run2(data, recentItems, task)
	if err != nil {
		fmt.Printf("failed to compare, taskID %v, time %v, %s\n", task.TaskID, t.Unix(), err)
		return
	}

	recentItems.Data = result
	recentItems.Time = time.Now().Unix()
	recentItems.Length = len(result)

	go bot.MercariPushMsg(recentItems, task.Owner, task.Group)

	err = utils.UpdateDataDB(recentItems)
	if err != nil {
		fmt.Printf("failed to update AnalysisData, taskID %v, time %v, %s", recentItems.TaskID, t.Unix(), err)
		return
	}
}
