package bot

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"bookq.xyz/mercari-watchdog/models/analysisdata"
	"bookq.xyz/mercari-watchdog/models/analysistask"
	"bookq.xyz/mercari-watchdog/models/group"
	"bookq.xyz/mercari-watchdog/tools"

	Pichubot "github.com/0ojixueseno0/go-Pichubot"
)

var TIME_1H_STRING = []string{"1时", "1小时", "一小时", "60分", "3600秒"}

// group: help command handler
func handlerHelp(e Pichubot.MessageGroup) {
	if strings.EqualFold(e.Message, "/help") || strings.EqualFold(e.Message, ".help") ||
		strings.EqualFold(e.Message, "/帮助") {
		Pichubot.SendGroupMsg("指令:\n蹲煤:添加蹲煤任务\n查询:查询自己添加的任务\n删除<任务id>:删除指定id的蹲煤任务", e.GroupID)
	}
}

// group: all other commands' handler
func handlerGroupMsg(e Pichubot.MessageGroup) {
	msgarr := strings.Split(e.RawMessage, "\n")
	for i := 0; i < len(msgarr); i++ {
		msgarr[i] = strings.TrimRight(msgarr[i], "\r")
	}
	switch {
	case msgarr[0] == "蹲煤":
		msg, err := createTask(msgarr[1:], e.Sender.UserID, e.GroupID)
		if err != nil {
			OperationChan <- tools.PushMsg{Dst: e.GroupID, S: []string{fmt.Sprintf("添加失败了，这是调试用的error:%v", err)}}
			return
		}
		OperationChan <- tools.PushMsg{Dst: e.GroupID, S: []string{msg}}
	case msgarr[0] == "查询":
		res, err := analysistask.GetByQQ(e.Sender.UserID, e.GroupID)
		if err != nil {
			OperationChan <- tools.PushMsg{Dst: e.GroupID, S: []string{fmt.Sprintf("查询失败了，这是调试用的error:%v", err)}}
			return
		}
		var msgbuilder strings.Builder
		msgbuilder.Grow(1024)
		msgbuilder.WriteString("任务:")
		for _, item := range res {
			msgbuilder.WriteString("\n")
			msgbuilder.WriteString(item.FormatSimplifiedChinese())
		}
		OperationChan <- tools.PushMsg{Dst: e.GroupID, S: []string{msgbuilder.String()}}
	case strings.Index(msgarr[0], "删除") == 0:
		tmp := strings.Trim(strings.Trim(msgarr[0], "删除"), " ")
		msgarr = strings.Split(tmp, " ")
		idarr := make([]int32, len(msgarr))
		for i, item := range msgarr {
			tmp, err := strconv.Atoi(item)
			if err != nil {
				OperationChan <- tools.PushMsg{Dst: e.GroupID, S: []string{"任务编号转换失败了，请先确认输入是否是纯数字"}}
				return
			}
			idarr[i] = int32(tmp)
		}
		err := deleteTask(idarr, e.Sender.UserID)
		if err != nil {
			OperationChan <- tools.PushMsg{Dst: e.GroupID, S: []string{fmt.Sprintf("查询失败了，这是调试用的error:%v", err)}}
			return
		}
		OperationChan <- tools.PushMsg{Dst: e.GroupID, S: []string{"理论上是全部删除成功了(没给id的情况下当然也算成功)"}}
	case msgarr[0] == "/pushMsg":
		if e.Sender.UserID != Pichubot.PichuBot.Config.MasterQQ {
			return
		}
		toGroup, err := strconv.ParseInt(msgarr[1], 10, 64)
		if err != nil {
			OperationChan <- tools.PushMsg{Dst: e.GroupID, S: []string{err.Error()}}
		}
		OperationChan <- tools.PushMsg{Dst: toGroup, S: msgarr[2:]}
	}
}

// Accept Group invite that in collection GroupSettings
func handlerGroupRequest(r Pichubot.GroupRequest) {
	res, err := group.FindWhitelist(r.GroupId)
	if err != nil || !res {
		Pichubot.SetGroupInviteRequest(r.Flag, res, "可能不在白名单内")
	}
	Pichubot.SetGroupInviteRequest(r.Flag, true, "")
}

func handlerGroupLeave(r Pichubot.GroupDecrease) {
	if r.SubType != "kick_me" {
		return
	}

	if err := analysistask.DeleteByGroup(r.GroupId); err != nil {
		Pichubot.Logger.Alertf("Delete fail when removing kicked task, %s", err)
	}

	if err := analysisdata.DeleteByGroup(r.GroupId); err != nil {
		Pichubot.Logger.Alertf("Delete fail when removing kicked task, %s", err)
	}

	Pichubot.Logger.Infof("Kicked task deleted, %d", r.GroupId)
}

// Push AnalysisData to msg queue.
func MercariPushMsg(data analysisdata.AnalysisData, owner int64, group int64) {
	if data.Length <= 0 {
		return
	}

	msg := data.FormatSimplifiedChinese()
	msg = fmt.Sprintf("[CQ:at,qq=%v]\n%s", owner, msg)

	// tools.PushMsg.S must be []string
	content := tools.PushMsg{Dst: group, S: []string{msg}}
	switch data.Length / 2 {
	case 0:
		Push1to2Chan <- content
	case 1:
		Push3to4Chan <- content
	default:
		Push5upChan <- content
	}
}

// Push channel with priority
func msgPushService() {
	tick := time.NewTicker(300 * time.Millisecond)
	for {
		select {
		case push := <-OperationChan:
			pushCore(push)
		default:
			select {
			case push := <-OperationChan:
				pushCore(push)
			case push := <-Push1to2Chan:
				pushCore(push)
			default:
				select {
				case push := <-OperationChan:
					pushCore(push)
				case push := <-Push1to2Chan:
					pushCore(push)
				case push := <-Push3to4Chan:
					pushCore(push)
				default:
					select {
					case push := <-OperationChan:
						pushCore(push)
					case push := <-Push1to2Chan:
						pushCore(push)
					case push := <-Push3to4Chan:
						pushCore(push)
					case push := <-Push5upChan:
						pushCore(push)
					}
				}
			}
		}
		<-tick.C
	}
}

// the only function to send tools.PushMsg
func pushCore(push tools.PushMsg) {
	for _, item := range push.S {
		Pichubot.SendGroupMsg(item, push.Dst)
	}
}
