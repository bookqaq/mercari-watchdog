package tasks

import (
	"fmt"
	"time"

	"bookq.xyz/mercariWatchdog/compare"
	"bookq.xyz/mercariWatchdog/utils"
	"github.com/bookqaq/goForMercari/mercarigo"
)

func Boot() {
	debug_ticker_1s := time.NewTicker(1 * time.Second)
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
		case t := <-debug_ticker_1s.C:
			runWorkflow(1, t)
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
	data, err := mercarigo.Mercari_search(task.Keywords[0], task.Sort, task.Order, "", 30, task.MaxPage)
	if err != nil {
		fmt.Printf("failed to search, taskID %v, time %v", task.TaskID, t.Unix())
		return
	}

	data = utils.KeywordFilter(task, data)

	result, err := compare.Run(data, task)
	if err != nil {
		fmt.Printf("failed to compare, taskID %v, time %v, error:\n%s", task.TaskID, t.Unix(), err)
		return
	}

	fmt.Println(result)
}
