package bot

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"bookq.xyz/mercari-watchdog/utils"
	"bookq.xyz/mercari-watchdog/utils/analysisdata"
	"bookq.xyz/mercari-watchdog/utils/analysistask"
	"bookq.xyz/mercari-watchdog/utils/group"

	Pichubot "github.com/0ojixueseno0/go-Pichubot"
)

var TIME_1H_STRING = []string{"1时", "1小时", "一小时", "60分", "3600秒"}

func handlerGroupMsg(e Pichubot.MessageGroup) {
	msgarr := strings.Split(e.RawMessage, "\n")
	for i := 0; i < len(msgarr); i++ {
		msgarr[i] = strings.TrimRight(msgarr[i], "\r")
	}
	switch {
	case msgarr[0] == "蹲煤":
		msg, err := createTask(msgarr[1:], e.Sender.UserID, e.GroupID)
		if err != nil {
			OperationChan <- utils.PushMsg{Dst: e.GroupID, S: []string{fmt.Sprintf("添加失败了，这是调试用的error:%v", err)}}
			return
		}
		OperationChan <- utils.PushMsg{Dst: e.GroupID, S: []string{msg}}
	case msgarr[0] == "查询":
		res, err := analysistask.GetByQQ(e.Sender.UserID, e.GroupID)
		if err != nil {
			OperationChan <- utils.PushMsg{Dst: e.GroupID, S: []string{fmt.Sprintf("查询失败了，这是调试用的error:%v", err)}}
			return
		}
		var msgbuilder strings.Builder
		msgbuilder.Grow(1024)
		msgbuilder.WriteString("任务:")
		for _, item := range res {
			msgbuilder.WriteString("\n")
			msgbuilder.WriteString(item.FormatSimplifiedChinese())
		}
		OperationChan <- utils.PushMsg{Dst: e.GroupID, S: []string{msgbuilder.String()}}
	case strings.Index(msgarr[0], "删除") == 0:
		tmp := strings.Trim(msgarr[0], "删除")
		msgarr = strings.Split(tmp, " ")
		idarr := make([]int32, len(msgarr))
		for i, item := range msgarr {
			tmp, err := strconv.Atoi(item)
			if err != nil {
				OperationChan <- utils.PushMsg{Dst: e.GroupID, S: []string{"任务编号转换失败了，请先确认输入是否是纯数字"}}
				return
			}
			idarr[i] = int32(tmp)
		}
		err := deleteTask(idarr, e.Sender.UserID)
		if err != nil {
			OperationChan <- utils.PushMsg{Dst: e.GroupID, S: []string{fmt.Sprintf("查询失败了，这是调试用的error:%v", err)}}
			return
		}
		OperationChan <- utils.PushMsg{Dst: e.GroupID, S: []string{"理论上是全部删除成功了(没给id的情况下当然也算成功)"}}
	case msgarr[0] == "/pushMsg":
		if e.Sender.UserID != Pichubot.PichuBot.Config.MasterQQ {
			return
		}
		toGroup, err := strconv.ParseInt(msgarr[1], 10, 64)
		if err != nil {
			OperationChan <- utils.PushMsg{Dst: e.GroupID, S: []string{err.Error()}}
		}
		OperationChan <- utils.PushMsg{Dst: toGroup, S: msgarr[2:]}
	}
}

// Accept Group invite that in collection GroupWhitelist
func handlerGroupRequest(r Pichubot.GroupRequest) {
	res, err := group.FindWhitelist(r.GroupId)
	if err != nil || !res {
		Pichubot.SetGroupInviteRequest(r.Flag, res, "可能不在白名单内")
	}
	Pichubot.SetGroupInviteRequest(r.Flag, true, "")
}

// Push AnalysisData to msg queue.
func MercariPushMsg(data analysisdata.AnalysisData, owner int64, group int64) {
	if data.Length <= 0 {
		return
	}

	msg := data.FormatSimplifiedChinese()
	msg = fmt.Sprintf("[CQ:at,qq=%v]\n%s", owner, msg)

	// utils.PushMsg.S must be []string
	content := utils.PushMsg{Dst: group, S: []string{msg}}
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
	tick := time.Tick(3 * time.Second)
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
		<-tick
	}
}

func pushCore(push utils.PushMsg) {
	for _, item := range push.S {
		Pichubot.SendGroupMsg(item, push.Dst)
	}
}
