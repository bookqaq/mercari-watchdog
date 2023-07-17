package analysisdata

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"bookq.xyz/mercari-watchdog/tools"
	"github.com/bookqaq/mer-wrapper/common"
	wrapperv2 "github.com/bookqaq/mer-wrapper/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var wd, _ = os.Getwd()

type AnalysisData struct {
	ID       primitive.ObjectID        `bson:"_id"`
	Keywords []string                  `bson:"keyword"`
	TaskID   int32                     `bson:"taskID"` // a primary-key-alike value
	Time     int64                     `bson:"time"`   // unix time
	Length   int                       `bson:"length"` // length of data
	Data     []wrapperv2.MercariV2Item `bson:"data"`
}

// privleged one display items in mulitple messages
func (d *AnalysisData) PrivlegedFormatSimplifiedChinese() []string {
	location, _ := time.LoadLocation("Asia/Shanghai")
	res := make([]string, 1, 6)
	res[0] = fmt.Sprintf("任务ID:%v\n关键词:%s\n时间:%s\n蹲到符合要求的结果数为%v",
		d.TaskID, tools.ConcatKeyword(d.Keywords), time.Unix(d.Time, 0).In(location).Format("2006-01-02 15:04:05"), d.Length)

	if d.Length > 0 {
		for _, item := range d.Data {
			updated, _ := strconv.ParseInt(item.Updated, 10, 64)
			filepath, err := saveWebpImage(item.ImageURL[0])
			if err != nil {
				log.Printf("发送信息时出现错误: %v", err)
			}

			tmp := fmt.Sprintf("[CQ:image,file=file://%s/%s]名称:%s\n价格:%vyen\n更新时间:%s\n链接:%s",
				wd, filepath, item.ProductName, item.Price,
				time.Unix(updated, 0).In(location).Format("2006-01-02 15:04:05"), item.ProductId)
			res = append(res, tmp)
		}
	}
	return res
}

// normal one display items in one message
func (d *AnalysisData) FormatSimplifiedChinese() string {
	location, _ := time.LoadLocation("Asia/Shanghai")
	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("任务ID:%v\n关键词:%s\n时间:%s\n",
		d.TaskID, tools.ConcatKeyword(d.Keywords), time.Unix(d.Time, 0).In(location).Format("2006-01-02 15:04:05")))

	if d.Length > 0 {
		for _, item := range d.Data {
			filepath, err := saveWebpImage(item.ImageURL[0])
			if err != nil {
				log.Printf("发送信息时出现错误: %v", err)
			}

			fmt.Println(filepath, PathExists(wd+"/"+filepath))

			builder.WriteString(fmt.Sprintf("\n[CQ:image,file=file://%s/%s]\n名称:%s\n价格:%vyen\n链接:%s",
				wd, filepath, item.ProductName, item.Price, item.ProductId))
		}
	}
	return builder.String()
}

// The layout string used by the Parse function and Format method
// shows by example how the reference time should be represented.
// We stress that one must show how the reference time is formatted,
// not a time of the user's choosing. Thus each layout string is a
// representation of the time stamp,
//	Jan 2 15:04:05 2006 MST
// An easy way to remember this value is that it holds, when presented
// in this order, the values (lined up with the elements above):
//	  1 2  3  4  5    6  -7

func saveWebpImage(url string) (string, error) {
	filenameStart := strings.LastIndex(url, "/")
	if filenameStart == -1 || filenameStart+1 >= len(url) {
		return "", errors.New("cant find file name")
	}

	questionMarkPos := strings.LastIndex(url, "?")
	if questionMarkPos == -1 || questionMarkPos < filenameStart {
		questionMarkPos = len(url)
	}

	filename := fmt.Sprintf("%s/%s", "files", url[filenameStart+1:questionMarkPos])

	if PathExists(filename) {
		return filename, nil
	}

	res, err := common.Client.Content.Get(url)
	if err != nil {
		return "", fmt.Errorf("request get failed: %w", err)
	}
	defer res.Body.Close()

	fp, err := os.Create(filename)
	if err != nil {
		return "", err
	}
	defer fp.Close()

	if _, err := io.Copy(fp, res.Body); err != nil {
		return "", err
	}

	return filename, nil
}

func PathExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil || errors.Is(err, os.ErrExist)
}
