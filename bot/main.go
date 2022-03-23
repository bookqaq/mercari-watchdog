package bot

import (
	"bookq.xyz/mercariWatchdog/utils"
	Pichubot "github.com/0ojixueseno0/go-Pichubot"
)

func Boot() {
	PushMsgChan = make(chan utils.PushMsg, 10)

	Pichubot.Listeners.OnGroupMsg = append(Pichubot.Listeners.OnGroupMsg, handlerGroupMsg)
	Pichubot.Listeners.OnGroupRequest = append(Pichubot.Listeners.OnGroupRequest, handlerGroupRequest)
	bot := Pichubot.NewBot()
	bot.Config = Pichubot.Config{
		Loglvl:   Pichubot.LOGGER_LEVEL_INFO,
		Host:     "127.0.0.1:28285",
		MasterQQ: 295589844,
		Path:     "/",
		MsgAwait: true,
	}
	go msgPushService()
	bot.Run()
}
