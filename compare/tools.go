package compare

import (
	"strconv"
	"strings"
	"unicode/utf8"

	"bookq.xyz/mercari-watchdog/tools"
	wrapperv2 "github.com/bookqaq/mer-wrapper/v2"
)

func compNewTimestamp(data []wrapperv2.MercariV2Item, uptime int64) int {
	i := 0
	for _, item := range data {
		updated, _ := strconv.ParseInt(item.Updated, 10, 64)
		if updated < uptime {
			break
		}
		i++
	}
	return i
}

// service: format item description and judge
func compDescriptionFilter(keywords []string, title string, description string) bool {
	// loads of runes used by yhm to split lines/words and split those into array
	descrpition_arr := strings.Split(tools.StringMultipleReplacer(description, []rune{'\n', '\u3000', '\xa0', '\\', '、', '/'}, ' '), " ")
	del_count := tools.DeleteInvalidItem(descrpition_arr, "")
	descrpition_arr = descrpition_arr[:len(descrpition_arr)-del_count]

	// get position of "検索用"
	var word_mark [][2]int
	for i, item := range descrpition_arr {
		if strings.Contains(item, Config.const_Kensaku) {
			tmp := getKnownKensaku(descrpition_arr, i)
			if len(tmp) >= Config.MinmumLineCount {
				word_mark = append(word_mark, tmp)
			}
			break
		}
	}

	// bet kensaku in description
	if len(word_mark) <= 0 {
		word_mark = betKensaku(descrpition_arr)
	}

	// allow item that find no kensaku words
	if len(word_mark) <= 0 {
		return true
	}

	// delete words if exists
	for i := len(word_mark) - 1; i >= 0; i-- {
		cutKnownKensaku(descrpition_arr, word_mark[i])
		descrpition_arr = descrpition_arr[:len(descrpition_arr)-(word_mark[i][1]-word_mark[i][0])]
	}

	// calculate percentage of keyword contains
	contain_count := 0
	for _, item := range keywords {
		if strings.Contains(title, item) {
			contain_count++
		}
	}

	// simplify judges about return value (forced by gopls)
	return float32(contain_count)/float32(len(keywords)) >= Config.V2KeywordMatchMin
}

// move words that after kensaku words forward, need to slice manually
func cutKnownKensaku(arr []string, pos [2]int) {
	for i, j := pos[0], pos[1]; i < pos[1] && j < len(arr); i, j = i+1, j+1 {
		arr[i] = arr[j]
	}
}

// find kensaku after start
func getKnownKensaku(arr []string, start int) [2]int {
	var mark [2]int
	mark[0] = start
	for i := start; i < len(arr); i++ {
		if linelen := utf8.RuneCount([]byte(arr[i])); linelen <= 0 || linelen > Config.MinimumRuneLength {
			break
		}
		mark[1] = i
	}
	return mark
}

// bet the start position of kensaku and find its end
func betKensaku(arr []string) [][2]int {
	mark_storage := make([][2]int, 0, 2)
	conlen := len(arr)

	for i := 0; i < conlen; i++ {
		var mark [2]int
		if linelen := utf8.RuneCount([]byte(arr[i])); (linelen > 0 && linelen < Config.MinimumRuneLength) || linelen > Config.MaximumRuneLength {
			mark[0] = i
			for i++; i < conlen; i++ {
				if linelen := utf8.RuneCount([]byte(arr[i])); (linelen > 0 && linelen < Config.MinimumRuneLength) || linelen > Config.MaximumRuneLength {
					mark[1] = i
				} else {
					break
				}
			}
			if mark[1]-mark[0] >= Config.MinmumLineCount {
				mark_storage = append(mark_storage, mark)
			}
		}
	}
	return mark_storage
}
