package utils

import (
	"time"

	"bookq.xyz/mercari-watchdog/database"
	"bookq.xyz/mercari-watchdog/utils/analysistask"
	"bookq.xyz/mercari-watchdog/utils/fetchdata"
)

func Init() {
	database.Connect()
	go analysistask.AddTaskBuffer()
	go fetchdata.TickClearExpired(120 * time.Second)
}
