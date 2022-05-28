package main

import (
	"bookq.xyz/mercari-watchdog/bot"
	"bookq.xyz/mercari-watchdog/database"
	"bookq.xyz/mercari-watchdog/tasks"
	"bookq.xyz/mercari-watchdog/webapi"
)

func main() {
	database.Connect()
	go bot.Boot()
	go webapi.Boot()
	tasks.Boot()
}
