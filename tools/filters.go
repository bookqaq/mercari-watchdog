package tools

import (
	"errors"
	"strconv"
	"strings"
	"sync"

	"bookq.xyz/mercari-watchdog/models/blacklist"
	wrapperv2 "github.com/bookqaq/mer-wrapper/v2"
)

var blockedSellers map[int64]blacklist.BlockedSeller
var lock sync.Mutex

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
func KeywordFilter(keywords []string, data []wrapperv2.MercariV2Item) []wrapperv2.MercariV2Item {
	ans, lenKeyword := make([]wrapperv2.MercariV2Item, 0, len(data)), len(keywords)
	for _, d := range data {
		matched := 0
		title_sp := strings.Split(StringMultipleReplacer(d.ProductName, []rune{'\u3000', '\xa0', '、', '/'}, ' '), " ")
		for _, keyword := range keywords {
			findMatch := false
			for _, titleWord := range title_sp {
				if keywordPercentage(titleWord, keyword) >= thresholdSingleKeyword {
					findMatch = true
				}
			}

			if findMatch {
				matched++
			}
		}
		if matched >= lenKeyword {
			ans = append(ans, d)
		}
	}
	return ans
}

// Return items that price in task.TargetPrice
func PriceFilter(price [2]int, data []wrapperv2.MercariV2Item) []wrapperv2.MercariV2Item {
	result := make([]wrapperv2.MercariV2Item, 0, len(data))
	if price[0] >= 0 && price[1] >= price[0] {
		for _, item := range data {
			itemPrice, _ := strconv.ParseInt(item.Price, 10, 64)
			priceInt := int(itemPrice)
			if priceInt >= price[0] && priceInt <= price[1] {
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
func BlockedSellerFilter(data []wrapperv2.MercariV2Item) []wrapperv2.MercariV2Item {
	if blockedSellers == nil {
		panic(errors.New("BlockedSeller Must be a map, not nil"))
	}

	result := make([]wrapperv2.MercariV2Item, 0, len(data))
	for _, item := range data {
		seller, _ := strconv.ParseInt(item.Seller, 10, 64)
		if data, ok := blockedSellers[seller]; !ok || blockedSellersIfInclude(data.Reason) {
			result = append(result, item)
		}
	}
	return result
}

// filter tools
func keywordPercentage(s, compareTo string) float64 {
	return float64(LongestCommon([]rune(s), []rune(compareTo))) / float64(len([]rune(compareTo)))
}
