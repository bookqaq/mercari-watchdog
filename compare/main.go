package compare

import (
	"fmt"

	"bookq.xyz/mercariWatchdog/utils"
	"github.com/bookqaq/goForMercari/mercarigo"
)

// implement compare, updated methods would be available in future
func Run(data []mercarigo.MercariItem, recentData utils.AnalysisData, task utils.AnalysisTask) ([]mercarigo.MercariItem, error) {
	i, itemlen, uptime := 0, len(data), recentData.Time

	for _, item := range data {
		if item.Updated < uptime {
			break
		}
		i++
	}

	if i >= itemlen {
		return nil, fmt.Errorf("items compare fail, no item update")
	}

	data = data[:i]

	data = utils.PriceFilter(task, data)

	return data, nil
}
