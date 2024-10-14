package main

import (
	"log"

	"github.com/awangelo/MonBot/internal/bot"
	"github.com/awangelo/MonBot/internal/config"
	"github.com/awangelo/MonBot/internal/handler"
	"github.com/awangelo/MonBot/internal/utils"
)

func main() {
	// Obter variáveis de ambiente
	client, groupJID, delayMins, err := config.InitBot()
	if err != nil {
		log.Fatalf("Erro ao inicializar o bot: %v", err)
	}
	// Inicializa o handler
	handler.InitHandler(client, groupJID)

	// Login do bot
	bot.Login(client)
	defer client.Disconnect()

	// Listar grupos
	bot.ListGroups(client)

	// Enviar mensagem automatica
	go bot.SendScheduledMessage(client, groupJID, delayMins)

	utils.PreventExit()
}
