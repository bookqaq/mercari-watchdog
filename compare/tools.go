package compare

import (
	"reflect"
	"strings"
	"unicode/utf8"

	"github.com/bookqaq/goForMercari/mercarigo"
)

func compNewTimestamp(data []mercarigo.MercariItem, uptime int64) int {
	i := 0
	for _, item := range data {
		if item.Updated < uptime {
			break
		}
		i++
	}
	return i
}

func compDescriptionFilter(keywords []string, title string, description string) bool {
	if len(keywords) == 1 {
		return true
	}
	descrpition_arr := strings.Split(strings.ReplaceAll(description, "\n", " "), " ")
	del_count := deleteInvalidItem(descrpition_arr, "")
	descrpition_arr = descrpition_arr[:len(descrpition_arr)-del_count]

	var word_mark [][2]int
	for i, item := range descrpition_arr {
		if strings.Contains(item, Config.const_V2Kensaku) {
			tmp := getKnownKensaku(descrpition_arr, i)
			if len(tmp) >= Config.V2MinmumLineCount {
				word_mark = append(word_mark, tmp)
			}
			break
		}
	}

	if len(word_mark) <= 0 {
		word_mark = betKensaku(descrpition_arr)
	}

	if len(word_mark) <= 0 {
		return true
	}

	for i := len(word_mark) - 1; i >= 0; i-- {
		cutKnownKensaku(descrpition_arr, word_mark[i])
		descrpition_arr = descrpition_arr[:len(descrpition_arr)-(word_mark[i][1]-word_mark[i][0])]
	}

	contain_count := 0
	for _, item := range keywords {
		if strings.Contains(title, item) {
			contain_count++
		}
	}

	if float32(contain_count)/float32(len(keywords)) >= Config.V2KeywordMatchMin {
		return true
	}

	return false
}

func deleteInvalidItem[T any](src []T, value T) int {
	deleted, formerpt, length := 0, 0, len(src)
	for i := 0; i < length; i++ {
		for ; reflect.DeepEqual(src[i], value); i++ {
			deleted++
		}
		src[formerpt] = src[i]
		formerpt++
	}
	return deleted
}

func cutKnownKensaku(arr []string, pos [2]int) {
	for i, j := pos[0], pos[1]; i < pos[1] && j < len(arr); i, j = i+1, j+1 {
		arr[i] = arr[j]
	}
}

func getKnownKensaku(arr []string, start int) [2]int {
	var mark [2]int
	mark[0] = start
	for i := start; i < len(arr); i++ {
		if linelen := utf8.RuneCount([]byte(arr[i])); linelen <= 0 || linelen > Config.V2MinimumRuneLength {
			break
		}
		mark[1] = i
	}
	return mark
}

func betKensaku(arr []string) [][2]int {
	mark_storage := make([][2]int, 0)
	conlen := len(arr)

	for i := 0; i < conlen; i++ {
		var mark [2]int
		if linelen := utf8.RuneCount([]byte(arr[i])); linelen > 0 && linelen < Config.V2MinimumRuneLength {
			mark[0] = i
			for i++; i < conlen; i++ {
				if linelen := utf8.RuneCount([]byte(arr[i])); linelen > 0 && linelen < Config.V2MinimumRuneLength {
					mark[1] = i
				} else {
					break
				}
			}
			if mark[1]-mark[0] > Config.V2MinmumLineCount {
				mark_storage = append(mark_storage, mark)
			}
		}
	}
	return mark_storage
}