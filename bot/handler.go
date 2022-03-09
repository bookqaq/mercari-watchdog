package bot

import (
	"fmt"
	"strconv"
	"strings"

	"bookq.xyz/mercariWatchdog/utils"
	Pichubot "github.com/0ojixueseno0/go-Pichubot"
)

var TIME_1H_STRING = []string{"1时", "1小时", "60分", "3600秒"}
var PushMsgChan chan utils.PushMsg

func handlerGroupMsg(e Pichubot.MessageGroup) {
	msgarr := strings.Split(e.Message, "\n")
	switch {
	case msgarr[0] == "蹲煤":
		msg, err := createTask(msgarr[1:], e.Sender.UserID, e.GroupID)
		if err != nil {
			PushMsgChan <- utils.PushMsg{Dst: e.GroupID, S: fmt.Sprintf("查询失败了，这是调试用的error:%v", err)}
			return
		}
		PushMsgChan <- utils.PushMsg{Dst: e.GroupID, S: msg}
	case msgarr[0] == "查询":
		res, err := utils.GetTasksByQQ(e.Sender.UserID)
		if err != nil {
			PushMsgChan <- utils.PushMsg{Dst: e.GroupID, S: fmt.Sprintf("查询失败了，这是调试用的error:%v", err)}
			return
		}
		msg := "任务:"
		for _, item := range res {
			msg += "\n"
			msg += item.FormatSimplifiedChinese()
		}
		PushMsgChan <- utils.PushMsg{Dst: e.GroupID, S: msg}
	case strings.Index(msgarr[0], "删除") == 0:
		msgarr = strings.Split(msgarr[0], " ")
		msgarr = msgarr[1:]
		idarr := make([]int32, len(msgarr))
		for i, item := range msgarr {
			tmp, err := strconv.Atoi(item)
			if err != nil {
				PushMsgChan <- utils.PushMsg{Dst: e.GroupID, S: "任务编号转换失败了，请先确认输入是否是纯数字"}
				return
			}
			idarr[i] = int32(tmp)
		}
		err := deleteTask(idarr)
		if err != nil {
			PushMsgChan <- utils.PushMsg{Dst: e.GroupID, S: fmt.Sprintf("查询失败了，这是调试用的error:%v", err)}
			return
		}
		PushMsgChan <- utils.PushMsg{Dst: e.GroupID, S: "理论上是全部删除成功了(没给id的情况下当然也算成功)"}
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
	if data.Length <= 0 {
		return
	}
	msgarr := data.FormatSimplifiedChinese()
	msgarr[0] = fmt.Sprintf("[CQ:at,qq=%v]\n", owner) + msgarr[0]
	for i, item := range msgarr {
		PushMsgChan <- utils.PushMsg{Dst: group, S: item}
		if i > 10 {
			break
		}
	}
}

func msgPushService() {
	for {
		push := <-PushMsgChan
		Pichubot.SendGroupMsg(push.S, push.Dst)
	}
}
