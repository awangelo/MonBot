package main

import (
	"log"

	"github.com/awangelo/MonBot/internal/bot"
	"github.com/awangelo/MonBot/internal/handler"
	"github.com/awangelo/MonBot/internal/utils"
)

func main() {
	client, groupJID, delayMins, err := bot.InitBot()
	if err != nil {
		log.Fatalf("Error creating bot: %v", err)
	}
	handler.InitHandler(client)
	handler.LoginByQr(client)
	defer client.Disconnect()

	utils.ListGroups(client)
	bot.SendScheduledMessage(client, groupJID, delayMins)

	utils.PreventExit()
}
