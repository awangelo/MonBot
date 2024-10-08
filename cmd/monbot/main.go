package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

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

	utils.ListGroups(client)
	bot.SendScheduledMessage(client, groupJID, delayMins)

	// Listen to Ctrl+C (you can also do something else that prevents the program from exiting)
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	client.Disconnect()
}
