package bot

import (
	"fmt"
	"strconv"
	"strings"

	"bookq.xyz/mercariWatchdog/utils"
	Pichubot "github.com/0ojixueseno0/go-Pichubot"
)

var TIME_1H_STRING = []string{"1时", "1小时", "60分", "3600秒"}

func handlerGroupMsg(e Pichubot.MessageGroup) {
	msgarr := strings.Split(e.Message, "\n")
	switch {
	case msgarr[0] == "/蹲煤":
		msg, err := createTask(msgarr[1:], e.Sender.UserID, e.GroupID)
		if err != nil {
			Pichubot.SendGroupMsg(fmt.Sprintf("添加任务失败了，这是调试用的error:%v", err), e.GroupID)
			return
		}
		Pichubot.SendGroupMsg(msg, e.GroupID)
	case msgarr[0] == "/查询":
		res, err := utils.GetTasksByQQ(e.Sender.UserID)
		if err != nil {
			Pichubot.SendGroupMsg(fmt.Sprintf("查询失败了，这是调试用的error:%v", err), e.GroupID)
			return
		}
		msg := "任务:"
		for _, item := range res {
			msg += "\n"
			msg += item.FormatSimplifiedChinese()
		}
		Pichubot.SendGroupMsg(msg, e.GroupID)
	case strings.Index(msgarr[0], "/删除") == 0:
		msgarr = strings.Split(msgarr[0], " ")
		msgarr = msgarr[1:]
		idarr := make([]int32, len(msgarr))
		for i, item := range msgarr {
			tmp, err := strconv.Atoi(item)
			if err != nil {
				Pichubot.SendGroupMsg("任务编号转换失败了，请先确认输入是否是纯数字", e.GroupID)
				return
			}
			idarr[i] = int32(tmp)
		}
		err := deleteTask(idarr)
		if err != nil {
			Pichubot.SendGroupMsg(fmt.Sprintf("删除失败了，这是调试用的error:%v", err), e.GroupID)
			return
		}
		Pichubot.SendGroupMsg("理论上是全部删除成功了(没给id的情况下当然也算成功)", e.GroupID)
	default:
		Pichubot.SendGroupMsg("程序有任何问题，请加好友问我295589844，未来会开留言功能向号主转发", e.GroupID)
	}
}

func handlerGroupRequest(r Pichubot.GroupRequest) {
	res, err := utils.FindWhitelist(r.GroupId)
	if err != nil || !res {
		Pichubot.SetGroupInviteRequest(r.Flag, res, "可能不在白名单内")
	}
	Pichubot.SetGroupInviteRequest(r.Flag, true, "")
}

func MercariPushMsg(data utils.AnalysisData, owner int64, group int64) {
	msgarr := data.FormatSimplifiedChinese()
	msgarr[0] = fmt.Sprintf("[CQ:at,qq=%v]\n", owner) + msgarr[0]
	for i, item := range msgarr {
		Pichubot.SendGroupMsg(item, group)
		if i > 5 {
			Pichubot.SendGroupMsg("新数据太多了，我没有使用经验怕封号，其他的新数据还请自己到mercari看，未来肯定会实现更好的筛选算法", group)
			break
		}
	}
}
