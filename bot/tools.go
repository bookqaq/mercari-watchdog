package bot

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"bookq.xyz/mercari-watchdog/datatype/analysisdata"
	"bookq.xyz/mercari-watchdog/datatype/analysistask"
	"bookq.xyz/mercari-watchdog/datatype/fetchdata"
	"bookq.xyz/mercari-watchdog/tools"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	merbot_tadd_link = "https://merbot.bookq.xyz/task/add/"
)

func createTask(params []string, qq int64, group int64) (string, error) {
	var result string
	switch {
	case len(params) == 0:
		result = "格式:\n" + "蹲煤\n" +
			"关键词:\n" + "目标价格:\n" + "搜索间隔:\n" + "搜索页数:\n" +
			"以下是举例:\n" + "蹲煤\n" + "关键词:プロセカ グリッター缶バッジ\n" + "搜索间隔:1小时\n" + "目标价格:100，500\n" + "搜索页数:3"

	case len(params) == 4:
		task, err := translateParams(params)
		if err != nil {
			return "", err
		}
		task.Owner = qq
		task.Group = group
		task.MustMatch = task.Keywords
		data := analysisdata.AnalysisData{
			ID:       primitive.NewObjectID(),
			Keywords: task.Keywords,
			TaskID:   task.TaskID,
			Length:   0,
			Data:     nil,
		}
		err = analysisdata.Insert(data)
		if err != nil {
			return "", err
		}
		analysistask.AddTaskChannel <- task
		result = "添加完成请求已提交，完成情况请通过查询进行查看"
	default:
		return "", fmt.Errorf("可能是参数过少或者过多?")
	}

	var authkey string
	authkey = strings.ReplaceAll(uuid.New().String(), "-", "")
	for fetchdata.IfExist(authkey) {
		authkey = strings.ReplaceAll(uuid.New().String(), "-", "")
	}

	fetchData := fetchdata.TaskAddFetchData{
		Override: fetchdata.FetchOverride{
			Owner: qq,
		},
		Auth:   authkey,
		Expire: time.Now().Unix() + int64(600),
	}

	err := fetchdata.Insert(fetchData)
	if err != nil {
		return "", err
	}
	result = fmt.Sprintf("在这个页面里也可以添加任务:%s%s\n%s", merbot_tadd_link, authkey, result)
	return result, nil
}

func deleteTask(tasks []int32, qq int64) error { //未来会添加信息所属的验证
	for _, item := range tasks {
		err := analysistask.Delete(item, qq)
		if err != nil {
			return err
		}
		err = analysisdata.Delete(item)
		if err != nil {
			return err
		}
	}
	return nil
}

func translateParams(params []string) (analysistask.AnalysisTask, error) {
	var tid int32
	tid = rand.Int31()
	for analysistask.IfExist(tid) {
		tid = rand.Int31()
	}
	task := analysistask.AnalysisTask{}
	task.ID = primitive.NewObjectID()
	task.TaskID = tid
	task.Sort = "created_time"
	task.Order = "desc"

	pmap := make(map[string]string, 4)
	for _, item := range params {
		splitindex := strings.Index(item, ":")
		if splitindex == -1 {
			return analysistask.AnalysisTask{}, fmt.Errorf("参数获取出了问题")
		}
		contmp := strings.TrimLeft(item[splitindex:], ":")
		if contmp == "" && strings.Index(item, "目标价格") != 0 {
			return analysistask.AnalysisTask{}, fmt.Errorf("参数获取出了问题")
		}
		tmp := []string{item[:splitindex], contmp}
		pmap[tmp[0]] = strings.Trim(tmp[1], " ")
	}
	if len(pmap) != 4 {
		return analysistask.AnalysisTask{}, fmt.Errorf("可能检测到了重复参数")
	}

	tmp, ok := pmap["目标价格"]
	if !ok {
		return analysistask.AnalysisTask{}, fmt.Errorf("解析目标价格失败")
	}
	satmp := strings.Split(tmp, "，")
	if len(satmp) != 2 {
		task.TargetPrice[0], task.TargetPrice[1] = -1, 0
	} else {
		itmp, err := strconv.Atoi(satmp[0])
		if err != nil {
			return analysistask.AnalysisTask{}, fmt.Errorf("解析目标价格失败")
		}
		task.TargetPrice[0] = itmp
		itmp, err = strconv.Atoi(satmp[1])
		if err != nil {
			return analysistask.AnalysisTask{}, fmt.Errorf("解析目标价格失败")
		}
		task.TargetPrice[1] = itmp
	}

	tmp, ok = pmap["搜索间隔"]
	if !ok {
		return analysistask.AnalysisTask{}, fmt.Errorf("解析时间间隔失败")
	}
	switch {
	case tmp == "10分" || tmp == "10分钟":
		task.Interval = 600
	default:
		task.Interval = 3600
	}

	tmp, ok = pmap["搜索页数"]
	if !ok {
		return analysistask.AnalysisTask{}, fmt.Errorf("解析搜索页数失败")
	}
	itmp, err := strconv.Atoi(tmp)
	if err != nil {
		return analysistask.AnalysisTask{}, fmt.Errorf("解析搜索页数失败")
	}
	task.MaxPage = itmp

	tmp, ok = pmap["关键词"]
	if !ok {
		return analysistask.AnalysisTask{}, fmt.Errorf("解析关键词失败")
	}
	tmp = strings.Replace(tmp, " ", "，", -1)
	satmp = strings.Split(tmp, "，")
	deleted := tools.DeleteInvalidItem(satmp, "")
	satmp = satmp[:len(satmp)-deleted]
	task.Keywords = satmp

	return task, nil
}
