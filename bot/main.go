package bot

import (
	"bookq.xyz/mercari-watchdog/tools"
	Pichubot "github.com/0ojixueseno0/go-Pichubot"
)

var OperationChan chan tools.PushMsg
var Push1to2Chan chan tools.PushMsg
var Push3to4Chan chan tools.PushMsg
var Push5upChan chan tools.PushMsg

func Boot() {
	// init message chanwith Priority
	OperationChan = make(chan tools.PushMsg, 10)
	Push1to2Chan = make(chan tools.PushMsg, 10)
	Push3to4Chan = make(chan tools.PushMsg, 10)
	Push5upChan = make(chan tools.PushMsg, 10)

	// add listeners
	Pichubot.Listeners.OnGroupMsg = append(Pichubot.Listeners.OnGroupMsg, handlerHelp, handlerGroupMsg)
	Pichubot.Listeners.OnGroupRequest = append(Pichubot.Listeners.OnGroupRequest, handlerGroupRequest)
	Pichubot.Listeners.OnGroupDecrease = append(Pichubot.Listeners.OnGroupDecrease, handlerGroupLeave)

	// config connection
	bot := Pichubot.NewBot()
	bot.Config = Pichubot.Config{
		Loglvl:   Pichubot.LOGGER_LEVEL_WARNING,
		Host:     "127.0.0.1:28285",
		MasterQQ: 295589844,
		Path:     "/",
		MsgAwait: true,
	}
	go msgPushService()
	bot.Run()
}
