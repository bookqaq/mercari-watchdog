package tasks

import (
	"fmt"
	"time"

	"bookq.xyz/mercariWatchdog/utils"
)

func Boot() {
	debug_ticker_1s := time.NewTicker(1 * time.Second)
	ticker_10m := time.NewTicker(600 * time.Second)
	ticker_1h := time.NewTicker(3600 * time.Second)

	for {
		var err error
		select {
		case t := <-ticker_10m.C:
			err = runWorkflow(600, t)
		case t := <-ticker_1h.C:
			err = runWorkflow(3600, t)
		case t := <-debug_ticker_1s.C:
			err = runWorkflow(1, t)
		}
		if err != nil {
			fmt.Println(err)
		}
	}
}

func runWorkflow(interval int, t time.Time) error {
	taskResults, err := utils.SearchAllTasks(interval)
	if err != nil {
		return fmt.Errorf("error during processing workflow %s : %v", t, interval)
	}

	for i, taskItem := range taskResults {
		err = runTask(taskItem)
		if err != nil {
			return fmt.Errorf("error running task %v in workflow %s : %v", i, t, interval)
		}
	}

	return nil
}

func runTask(tasks utils.AnalysisTask) error {
	//result, err := mercarigo.Mercari_search()
	return nil
}
