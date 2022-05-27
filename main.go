package main

import (
	"bookq.xyz/mercari-watchdog/bot"
	"bookq.xyz/mercari-watchdog/tasks"
	"bookq.xyz/mercari-watchdog/webapi"
)

func main() {
	go bot.Boot()
	go webapi.Boot()
	tasks.Boot()
}
