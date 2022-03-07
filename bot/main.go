package bot

import (
	Pichubot "github.com/0ojixueseno0/go-Pichubot"
)

func Boot() {
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
	bot.Run()
}
