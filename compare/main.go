package compare

import (
	"fmt"

	"bookq.xyz/mercariWatchdog/utils"
	"github.com/bookqaq/goForMercari/mercarigo"
)

// implement compare, updated methods would be available in future
func Run(items []mercarigo.MercariItem, task utils.AnalysisTask) ([]mercarigo.MercariItem, error) {
	recentItems, err := utils.GetDataDB(task.Keywords[0])
	if err != nil {
		return nil, err
	}

	i, itemlen, uptime := 0, len(items), recentItems.Data[0].Updated

	for _, item := range items {
		if item.Updated < uptime {
			break
		}
		i++
	}

	if i >= itemlen {
		return nil, fmt.Errorf("items compare fail, no item update")
	}

	items = items[:i]

	if task.TargetPrice[0] >= 0 && task.TargetPrice[1] >= task.TargetPrice[0] {
		result := make([]mercarigo.MercariItem, 0)
		for _, item := range items {
			if item.Price >= task.TargetPrice[0] && item.Price <= task.TargetPrice[1] {
				result = append(result, item)
			}
		}
		if len(result) == 0 {
			result = append(result, mercarigo.MercariItem{})
		}
		items = result
	}

	return items, nil
}
