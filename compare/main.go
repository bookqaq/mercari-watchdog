package compare

import (
	"fmt"

	"bookq.xyz/mercariWatchdog/utils"
	"github.com/bookqaq/goForMercari/mercarigo"
	merwrapper "github.com/bookqaq/mer-wrapper"
)

var Config = struct {
	const_V2Kensaku     string
	V2MinimumRuneLength int
	V2MinmumLineCount   int
	V2KeywordMatchMin   float32
}{
	const_V2Kensaku:     "検索用",
	V2MinimumRuneLength: 14,
	V2MinmumLineCount:   10,
	V2KeywordMatchMin:   0.4,
}

// implement compare, updated methods would be available in future
func Run(data []mercarigo.MercariItem, recentData utils.AnalysisData, task utils.AnalysisTask) ([]mercarigo.MercariItem, error) {
	data = utils.KeywordFilter(task, data)

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

func Run2(data []mercarigo.MercariItem, recentData utils.AnalysisData, task utils.AnalysisTask) ([]mercarigo.MercariItem, error) {
	uptime := recentData.Time

	i := compNewTimestamp(data, uptime)

	data = data[:i]
	data = utils.PriceFilter(task, data)

	//fdata := make([]mercarigo.MercariItem, 0, len(data)/4*3)
	for i, item := range data {
		desc, err := merwrapper.Client.Item(item.ProductId)
		if err != nil {
			return nil, err
		}
		if compDescriptionFilter(task.Keywords, item.ProductName, desc.Description) {
			//fdata = append(fdata, item)
			data[i].ProductName += " 该项目通过了filter"
		}
	}
	return data, nil
}
