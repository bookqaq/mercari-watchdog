package tools

import (
	"errors"
	"strings"

	"bookq.xyz/mercari-watchdog/models/blacklist"
	wrapperv1 "github.com/bookqaq/mer-wrapper/v1"
)

var blockedSellers map[int64]struct{}

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
func KeywordFilter(keywords []string, data []wrapperv1.MercariItem) []wrapperv1.MercariItem {
	ans, lenKeyword := make([]wrapperv1.MercariItem, 0, len(data)), len(keywords)
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

// return items that seller not in blacklist
func BlockedSellerFilter(data []wrapperv1.MercariItem) []wrapperv1.MercariItem {
	if blockedSellers == nil {
		panic(errors.New("BlockedSeller Must be a map, not nil"))
	}

	result := make([]wrapperv1.MercariItem, 0, len(data))
	for _, item := range data {
		if _, ok := blockedSellers[item.Seller.Id]; !ok {
			result = append(result, item)
		}
	}
	return result
}

// filter tools
func keywordPercentage(s, compareTo string) float64 {
	return float64(LongestCommon([]rune(s), []rune(compareTo))) / float64(len([]rune(compareTo)))
}
