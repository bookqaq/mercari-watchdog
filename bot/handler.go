package bot

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"

	"bookq.xyz/mercariWatchdog/utils"
	Pichubot "github.com/0ojixueseno0/go-Pichubot"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var TIME_1H_STRING = []string{"1时", "1小时", "60分", "3600秒"}

func handlerPrivateMsg(e Pichubot.MessagePrivate) {
	msgarr := strings.Split(e.Message, "\n")
	switch msgarr[0] {
	case "蹲煤":
		msg, err := createTask(msgarr[1:], e.Sender.UserID)
		if err != nil {
			Pichubot.SendPrivateMsg(fmt.Sprintf("添加任务失败了，这是调试用的error:%v", err), e.Sender.UserID)
		}
		Pichubot.SendPrivateMsg(msg, e.Sender.UserID)
	case "查询":
		res, err := utils.GetTasksByQQ(e.Sender.UserID)
		if err != nil {
			Pichubot.SendPrivateMsg(fmt.Sprintf("查询失败了，这是调试用的error:%v", err), e.Sender.UserID)
		}
		msg := "任务:"
		for _, item := range res {
			msg += "\n"
			msg += item.FormatSimplifiedChinese()
		}
		Pichubot.SendPrivateMsg(msg, e.Sender.UserID)
	default:
		Pichubot.SendPrivateMsg("程序有任何问题，请加好友问我295589844，未来会开留言功能向号主转发", e.Sender.UserID)
	}
}

func createTask(params []string, id int64) (string, error) {
	var result string
	switch {
	case len(params) == 0:
		result = "格式(使用中文逗号分隔一行中的内容):\n" + "蹲煤\n" +
			"关键词:\n" + "目标价格:\n" + "搜索间隔:\n" + "搜索页数:\n" +
			"注:关键词上限为5个,有多个关键词时会进行逐级筛选，目标价格中最低价为负数时视为任意价格，搜索间隔目前只有10分钟和1小时，每搜索页中有30个结果\n" +
			"以下是举例:\n" + "蹲煤\n" + "关键词:sasakure, lasah\n" + "搜索间隔:1小时" + "目标价格:100，500\n" + "搜索页数:3"

	case len(params) == 4:
		task, err := translateParams(params)
		if err != nil {
			return "", err
		}
		task.Owner = id
		err = utils.InsertDataDB(task)
		if err != nil {
			return "", err
		}
		utils.AddTaskChannel <- task
		result = "添加完成请求已提交，完成情况请通过查询进行查看"
	default:
		return "", fmt.Errorf("可能是参数过少或者过多?")
	}
	return result, nil
}

func translateParams(params []string) (utils.AnalysisTask, error) {
	pmap := make(map[string]string, 4)
	for _, item := range params {
		tmp := strings.Split(item, ":")
		if len(tmp) != 2 {
			return utils.AnalysisTask{}, fmt.Errorf("参数获取出了问题")
		}
		pmap[tmp[0]] = strings.Trim(tmp[1], " ")
	}
	if len(pmap) != 4 {
		return utils.AnalysisTask{}, fmt.Errorf("可能检测到了重复参数")
	}

	task := utils.AnalysisTask{}
	task.ID = primitive.NewObjectID()
	task.TaskID = rand.Int31()
	task.Sort = "created_time"
	task.Order = "desc"

	tmp, ok := pmap["目标价格"]
	if !ok {
		return utils.AnalysisTask{}, fmt.Errorf("解析目标价格失败")
	}
	satmp := strings.Split(tmp, "，")
	itmp, err := strconv.Atoi(satmp[0])
	if err != nil {
		return utils.AnalysisTask{}, fmt.Errorf("解析目标价格失败")
	}
	task.TargetPrice[0] = itmp
	itmp, err = strconv.Atoi(satmp[1])
	if err != nil {
		return utils.AnalysisTask{}, fmt.Errorf("解析目标价格失败")
	}
	task.TargetPrice[1] = itmp

	tmp, ok = pmap["搜索间隔"]
	if !ok {
		return utils.AnalysisTask{}, fmt.Errorf("解析时间间隔失败")
	}
	switch {
	case tmp == "10分" || tmp == "10分钟":
		task.Interval = 600
	default:
		task.Interval = 3600
	}

	tmp, ok = pmap["搜索页数"]
	if !ok {
		return utils.AnalysisTask{}, fmt.Errorf("解析搜索页数失败")
	}
	itmp, err = strconv.Atoi(tmp)
	if err != nil {
		return utils.AnalysisTask{}, fmt.Errorf("解析搜索页数失败")
	}
	task.MaxPage = itmp

	tmp, ok = pmap["关键词"]
	if !ok {
		return utils.AnalysisTask{}, fmt.Errorf("解析关键词失败")
	}
	satmp = strings.Split(tmp, "，")
	task.Keywords = satmp

	return task, nil
}
