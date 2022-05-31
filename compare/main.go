package compare

import (
	"bookq.xyz/mercari-watchdog/datatype/analysisdata"
	"bookq.xyz/mercari-watchdog/datatype/analysistask"
	"bookq.xyz/mercari-watchdog/tools"
	"github.com/bookqaq/goForMercari/mercarigo"
	merwrapper "github.com/bookqaq/mer-wrapper"
)

var Config = struct {
	const_Kensaku     string
	MinimumRuneLength int
	MaximumRuneLength int
	MinmumLineCount   int
}{
	const_Kensaku:     "検索用",
	MinimumRuneLength: 15,
	MaximumRuneLength: 50,
	MinmumLineCount:   10,
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
