package bot

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"bookq.xyz/mercari-watchdog/models/analysisdata"
	"bookq.xyz/mercari-watchdog/models/analysistask"
	"bookq.xyz/mercari-watchdog/models/fetchdata"
	"bookq.xyz/mercari-watchdog/tools"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// links about fronted
const (
	merbot_tadd_link = "https://merbot.bookq.xyz/task/add/"
)

// service about creating task.
func createTask(params []string, qq int64, group int64) (string, error) {
	var result string
	switch {
	// if no params provided
	case len(params) == 0:
		result = "格式:\n蹲煤\n关键词:\n目标价格:\n搜索间隔:\n搜索页数:\n"
	// if param's length correct
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
	// if length not correct, consider it to be invalid
	default:
		return "", fmt.Errorf("可能是参数过少或者过多?")
	}

	// token generator (used in web to differ users)
	// only excute when cases above passed
	var authkey string
	tmp := strings.Split(uuid.New().String(), "-")[:2]
	authkey = fmt.Sprintf("%s%s", tmp[0], tmp[1])
	for fetchdata.IfExist(authkey) {
		tmp := strings.Split(uuid.New().String(), "-")[:2]
		authkey = fmt.Sprintf("%s%s", tmp[0], tmp[1])
	}

	fetchData := fetchdata.TaskAddFetchData{
		Override: fetchdata.FetchOverride{
			Owner: qq,
			Group: group,
		},
		Auth:   authkey,
		Expire: time.Now().Unix() + int64(600),
	}

	err := fetchdata.Insert(fetchData)
	if err != nil {
		return "", err
	}

	// concat return value
	result = fmt.Sprintf("在这个页面里也可以添加任务:\n%s%s\n%s", merbot_tadd_link, authkey, result)
	return result, nil
}

// service about deleting tasks (support deleting multiple tasks)
func deleteTask(tasks []int32, qq int64) error {
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

// qq message extractor about adding tasks
func translateParams(params []string) (analysistask.AnalysisTask, error) {
	// generate a unique taskID
	var tid int32
	tid = rand.Int31()
	for analysistask.IfExist(tid) {
		tid = rand.Int31()
	}

	// init params in tasks
	task := analysistask.AnalysisTask{
		ID:     primitive.NewObjectID(),
		TaskID: tid,
		Sort:   "created_time",
		Order:  "desc",
	}

	// get and put(if exists) four params into a map
	pmap := make(map[string]string, 4)
	for _, item := range params {
		// split by ":"
		splitindex := strings.Index(item, ":")
		if splitindex == -1 {
			return analysistask.AnalysisTask{}, fmt.Errorf("参数获取出了问题")
		}

		// get text before ":" and delete ":"
		contmp := strings.TrimLeft(item[splitindex:], ":")
		if contmp == "" {
			return analysistask.AnalysisTask{}, fmt.Errorf("参数获取出了问题")
		}

		// merge and put data
		tmp := []string{item[:splitindex], contmp}
		pmap[tmp[0]] = strings.Trim(tmp[1], " ")
	}

	if len(pmap) != 4 {
		return analysistask.AnalysisTask{}, fmt.Errorf("可能检测到了重复参数")
	}

	// get and parse targetPrice [low, high]
	tmp, ok := pmap["目标价格"]
	if !ok {
		return analysistask.AnalysisTask{}, fmt.Errorf("解析目标价格失败")
	}
	satmp := strings.Split(tmp, "，")
	if len(satmp) != 2 {
		// autofill price when price is not provided correctly
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

	// get and parse Interval
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

	// get and parse Pages
	tmp, ok = pmap["搜索页数"]
	if !ok {
		return analysistask.AnalysisTask{}, fmt.Errorf("解析搜索页数失败")
	}
	itmp, err := strconv.Atoi(tmp)
	if err != nil {
		return analysistask.AnalysisTask{}, fmt.Errorf("解析搜索页数失败")
	}
	task.MaxPage = itmp

	// write and parse keywords
	tmp, ok = pmap["关键词"]
	if !ok {
		return analysistask.AnalysisTask{}, fmt.Errorf("解析关键词失败")
	}
	// preprocess and split
	tmp = strings.Replace(tmp, " ", "，", -1)
	satmp = strings.Split(tmp, "，")

	// remove " " in satmp(result of splitted keywords)
	deleted := tools.DeleteInvalidItem(satmp, "")
	satmp = satmp[:len(satmp)-deleted]

	task.Keywords = satmp

	return task, nil
}
