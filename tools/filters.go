package tools

import (
	"errors"
	"reflect"
	"strings"
	"sync"

	"bookq.xyz/mercari-watchdog/models/blacklist"
	wrapperv1 "github.com/bookqaq/mer-wrapper/v1"
)

var blockedSellers map[int64]blacklist.BlockedSeller
var lock *sync.Mutex

//
func RefreshBlockedSellers() {
	res, err := blacklist.BlockedSellerGetAll()
	if err != nil {
		panic(err)
	}

	blockMap_tmp := make(map[int64]blacklist.BlockedSeller, len(res))
	for _, seller := range res {
		blockMap_tmp[seller.UserID] = seller
	}
	lock.Lock()
	blockedSellers = blockMap_tmp
	lock.Unlock()
}

// filters
// Return items that match task.Keywords
func KeywordFilter(keywords []string, data []wrapperv1.MercariItem) []wrapperv1.MercariItem {
	for _, keyword := range keywords {
		tmp := make([]wrapperv1.MercariItem, 0, len(data))
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
func PriceFilter(price [2]int, data []wrapperv1.MercariItem) []wrapperv1.MercariItem {
	result := make([]wrapperv1.MercariItem, 0, len(data))
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

func blockedSellersIfInclude(reason string) bool {
	var res bool
	switch reason {
	case "虚假标价":
		res = true
	case "圈外检索词":
		res = false
	case "出售假谷":
		res = false
	default:
		res = true
	}
	return res
}

// return items that seller not in blacklist
func BlockedSellerFilter(data []wrapperv1.MercariItem) []wrapperv1.MercariItem {
	if blockedSellers == nil {
		panic(errors.New("BlockedSeller Must be a map, not nil"))
	}

	result := make([]wrapperv1.MercariItem, 0, len(data))
	for _, item := range data {
		if data, ok := blockedSellers[item.Seller.Id]; !ok || blockedSellersIfInclude(data.Reason) {
			result = append(result, item)
		}
	}
	return result
}

// Delete item in ordered array src that filter when reflect.DeepEqual(item, value) != true, return lenght deleted
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
