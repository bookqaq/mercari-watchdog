package main

import (
	"bookq.xyz/mercariWatchdog/bot"
	"bookq.xyz/mercariWatchdog/tasks"
	"bookq.xyz/mercariWatchdog/webapi"
)

func main() {
	go bot.Boot()
	go webapi.Boot()
	tasks.Boot()
}
