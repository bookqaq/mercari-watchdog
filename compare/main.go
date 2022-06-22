package compare

import (
	"strings"

	"bookq.xyz/mercari-watchdog/models/analysisdata"
	"bookq.xyz/mercari-watchdog/models/analysistask"
	"bookq.xyz/mercari-watchdog/tools"
	"github.com/bookqaq/goForMercari/mercarigo"
	merwrapper "github.com/bookqaq/mer-wrapper"
)

var Config = struct {
	const_Kensaku     string
	MinimumRuneLength int
	MaximumRuneLength int
	MinmumLineCount   int
	V2KeywordMatchMin float32
}{
	const_Kensaku:     "検索用",
	MinimumRuneLength: 15,
	MaximumRuneLength: 50,
	MinmumLineCount:   10,
	V2KeywordMatchMin: 0.4,
}

func Run2(data []mercarigo.MercariItem, recentData analysisdata.AnalysisData, task analysistask.AnalysisTask) ([]mercarigo.MercariItem, error) {
	uptime := recentData.Time

	i := compNewTimestamp(data, uptime)

	data = data[:i]
	data = tools.PriceFilter(task.TargetPrice, data)
	data = tools.BlockedSellerFilter(data)

	fdata := make([]mercarigo.MercariItem, 0, len(data)/4*3)
	for _, item := range data {
		desc, err := merwrapper.Client.Item(item.ProductId)
		if err != nil {
			return nil, err
		}
		if compDescriptionFilter(task.Keywords, item.ProductName, desc.Description) {
			fdata = append(fdata, item)
		}
	}
	return fdata, nil
}

func Run3(data []mercarigo.MercariItem, recentData analysisdata.AnalysisData, task analysistask.AnalysisTask) ([]mercarigo.MercariItem, error) {
	uptime := recentData.Time

	i := compNewTimestamp(data, uptime)

	data = data[:i]
	data = tools.PriceFilter(task.TargetPrice, data)
	data = tools.BlockedSellerFilter(data)

	fdata := make([]mercarigo.MercariItem, 0, len(data)/4*3)

	for _, item := range data {
		contain_flag := true
		for _, kw := range task.MustMatch {
			if !strings.Contains(item.ProductName, kw) {
				contain_flag = false
				break
			}
		}
		if contain_flag {
			fdata = append(fdata, item)
		}
	}
	return fdata, nil
}
