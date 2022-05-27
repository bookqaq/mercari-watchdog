package bot

import (
	"bookq.xyz/mercari-watchdog/utils"
	Pichubot "github.com/0ojixueseno0/go-Pichubot"
)

var OperationChan chan utils.PushMsg
var Push1to2Chan chan utils.PushMsg
var Push3to4Chan chan utils.PushMsg
var Push5upChan chan utils.PushMsg

func Boot() {
	OperationChan = make(chan utils.PushMsg, 10)
	Push1to2Chan = make(chan utils.PushMsg, 10)
	Push3to4Chan = make(chan utils.PushMsg, 10)
	Push5upChan = make(chan utils.PushMsg, 10)

	Pichubot.Listeners.OnGroupMsg = append(Pichubot.Listeners.OnGroupMsg, handlerGroupMsg)
	Pichubot.Listeners.OnGroupRequest = append(Pichubot.Listeners.OnGroupRequest, handlerGroupRequest)
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
