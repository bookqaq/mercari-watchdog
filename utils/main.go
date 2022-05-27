package utils

import (
	"bookq.xyz/mercari-watchdog/database"
	"bookq.xyz/mercari-watchdog/utils/analysistask"
)

func Init() {
	database.Connect()
	go analysistask.AddTaskBuffer()
}
