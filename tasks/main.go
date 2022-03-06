package tasks

import (
	"fmt"
	"time"

	"bookq.xyz/mercariWatchdog/compare"
	"bookq.xyz/mercariWatchdog/utils"
	"github.com/bookqaq/goForMercari/mercarigo"
)

func Boot() {
	utils.Init()

	debug_ticker := time.NewTicker(18400 * time.Second)
	ticker_10m := time.NewTicker(600 * time.Second)

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
			//runWorkflow(10, t)
		}
	}
}

func runWorkflow(interval int, t time.Time) {
	taskResults, err := utils.GetAllTasks(interval)
	if err != nil {
		fmt.Printf("error during processing workflow %s : %v", t, interval)
		return
	}

	for i, taskItem := range taskResults {
		go runTask(i, t, taskItem)
	}
}

func runTask(i int, t time.Time, task utils.AnalysisTask) {
	fmt.Printf("debug: task %v run\n", task.TaskID)
	data, err := mercarigo.Mercari_search(task.Keywords[0], task.Sort, task.Order, "", 30, task.MaxPage)
	if err != nil {
		fmt.Printf("failed to search, taskID %v, time %v\n", task.TaskID, t.Unix())
		return
	}

	fmt.Printf("debug: result found\n")

	data = utils.KeywordFilter(task, data)

	fmt.Printf("debug: filtered data %v\n", data)

	recentItems, err := utils.GetDataDB(task.TaskID)
	if err != nil {
		fmt.Printf("failed to get last search data, taskID %v, time %v, %s\n", task.TaskID, t.Unix(), err)
		return
	}
	result, err := compare.Run(data, recentItems, task)
	if err != nil {
		fmt.Printf("failed to compare, taskID %v, time %v, %s\n", task.TaskID, t.Unix(), err)
		return
	}

	recentItems.Data = result
	recentItems.Time = time.Now().Unix()
	recentItems.Length = len(result)

	err = utils.UpdateDataDB(recentItems)
	if err != nil {
		fmt.Printf("failed to update AnalysisData, taskID %v, time %v, %s", recentItems.TaskID, t.Unix(), err)
		return
	}
}
