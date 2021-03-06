package compare

import (
	"strings"

	"bookq.xyz/mercari-watchdog/models/analysisdata"
	"bookq.xyz/mercari-watchdog/models/analysistask"
	"bookq.xyz/mercari-watchdog/tools"
	merwrapper "github.com/bookqaq/mer-wrapper"
	wrapperv1 "github.com/bookqaq/mer-wrapper/v1"
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

// Legacy compare method that use keyword-match threshold
func Run2(data []wrapperv1.MercariItem, recentData analysisdata.AnalysisData, task analysistask.AnalysisTask) ([]wrapperv1.MercariItem, error) {
	uptime := recentData.Time

	i := compNewTimestamp(data, uptime)

	data = data[:i]
	data = tools.PriceFilter(task.TargetPrice, data)
	data = tools.BlockedSellerFilter(data)

	fdata := make([]wrapperv1.MercariItem, 0, len(data)/4*3)
	for _, item := range data {
		desc, err := merwrapper.Item(item.ProductId)
		if err != nil {
			return nil, err
		}
		if compDescriptionFilter(task.Keywords, item.ProductName, desc.Description) {
			fdata = append(fdata, item)
		}
	}
	return fdata, nil
}

// CompareV3 compare method, math exactly in task.MustMatch
func Run3(data []wrapperv1.MercariItem, recentData analysisdata.AnalysisData, task analysistask.AnalysisTask) ([]wrapperv1.MercariItem, error) {
	uptime := recentData.Time

	i := compNewTimestamp(data, uptime)

	data = data[:i]
	data = tools.PriceFilter(task.TargetPrice, data)
	data = tools.BlockedSellerFilter(data)

	fdata := make([]wrapperv1.MercariItem, 0, len(data)/4*3)

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
