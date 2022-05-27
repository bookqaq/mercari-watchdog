package compare

import (
	"bookq.xyz/mercari-watchdog/utils/analysisdata"
	"bookq.xyz/mercari-watchdog/utils/analysistask"
	"bookq.xyz/mercari-watchdog/utils/tools"
	"github.com/bookqaq/goForMercari/mercarigo"
	merwrapper "github.com/bookqaq/mer-wrapper"
)

var Config = struct {
	const_Kensaku     string
	MinimumRuneLength int
	MinmumLineCount   int
	MaximumLineCount  int
}{
	const_Kensaku:     "検索用",
	MinimumRuneLength: 15,
	MinmumLineCount:   10,
	MaximumLineCount:  100,
}

func Run3(data []mercarigo.MercariItem, recentData analysisdata.AnalysisData, task analysistask.AnalysisTask) ([]mercarigo.MercariItem, error) {
	uptime := recentData.Time

	i := compNewTimestamp(data, uptime)

	data = data[:i]
	data = tools.PriceFilter(task.TargetPrice, data)
	fdata := make([]mercarigo.MercariItem, 0, len(data)/4*3)
	for _, item := range data {
		desc, err := merwrapper.Client.Item(item.ProductId)
		if err != nil {
			return nil, err
		}
		if compDescriptionFilter(task.MustMatch, item.ProductName, desc.Description) {
			fdata = append(fdata, item)
		}
	}
	return fdata, nil
}
