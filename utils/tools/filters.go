package tools

import (
	"strings"

	"github.com/bookqaq/goForMercari/mercarigo"
)

// filters
// Return items that match task.Keywords
func KeywordFilter(keywords []string, data []mercarigo.MercariItem) []mercarigo.MercariItem {
	for _, keyword := range keywords {
		tmp := make([]mercarigo.MercariItem, 0, len(data))
		for _, item := range data {
			if strings.Contains(item.ProductName, keyword) {
				tmp = append(tmp, item)
			}
			data = tmp
		}
	}
	return data
}

// Return items that price in task.TargetPrice
func PriceFilter(price [2]int, data []mercarigo.MercariItem) []mercarigo.MercariItem {
	result := make([]mercarigo.MercariItem, 0)
	if price[0] >= 0 && price[1] >= price[0] {
		for _, item := range data {
			if item.Price >= price[0] && item.Price <= price[1] {
				result = append(result, item)
			}
		}
	} else {
		return data
	}
	return result
}
