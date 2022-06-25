package tools

import (
	"errors"
	"reflect"
	"strings"

	"bookq.xyz/mercari-watchdog/models/blacklist"
	"github.com/bookqaq/goForMercari/mercarigo"
)

var blockedSellers map[int64]struct{}

//
func RefreshBlockedSellers() {
	res, err := blacklist.BlockedSellerGetAll()
	if err != nil {
		panic(err)
	}

	blockMap_tmp := make(map[int64]struct{}, len(res))
	for _, seller := range res {
		blockMap_tmp[seller.UserID] = struct{}{}
	}
	blockedSellers = blockMap_tmp
}

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
	result := make([]mercarigo.MercariItem, 0, len(data))
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

// return items that seller not in blacklist
func BlockedSellerFilter(data []mercarigo.MercariItem) []mercarigo.MercariItem {
	if blockedSellers == nil {
		panic(errors.New("BlockedSeller Must be a map, not nil"))
	}

	result := make([]mercarigo.MercariItem, 0, len(data))
	for _, item := range data {
		if _, ok := blockedSellers[item.Seller.Id]; !ok {
			result = append(result, item)
		}
	}
	return result
}

// Delete item in ordered array src that return true in reflect.DeepEqual(item, value), return lenght deleted
func DeleteInvalidItem[T any](src []T, value T) int {
	deleted, formerpt, length := 0, 0, len(src)
	for i := 0; i < length; i++ {
		if reflect.DeepEqual(src[i], value) {
			deleted++
		} else {
			src[formerpt] = src[i]
			formerpt++
		}
	}
	return deleted
}
