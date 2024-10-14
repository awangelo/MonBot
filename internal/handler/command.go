package handler

import (
	"fmt"
	"log"

	"github.com/awangelo/MonBot/internal/bot"
	"github.com/awangelo/MonBot/internal/utils"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/types/events"
)

// Strings de resposta
const (
	helpMessage    = "\t\t\t\t*Comandos disponíveis:*\n\n* `!help`  -  Mostra esta mensagem\n* `!ping`  -  Responde com 'pong'\n* `!ram`  -  Mostra o uso de memória\n"
	pongMessage    = "pong 🏓"
	unknownMessage = "Comando desconhecido. Digite `!help` para ver a lista de comandos disponíveis."
)

// Formatar uso de memória
var ramMessage = func() string {
	m := utils.GetMemoryUsage()
	return fmt.Sprintf("Memória em uso: %d KB", m)
}

// HandleCommand trata comandos recebidos
func HandleCommand(client *whatsmeow.Client, msg *events.Message) {
	// Comando recebido
	command := msg.Message.GetConversation()
	log.Println("Comando recebido:", command)

	switch command {
	case "!":
		// Ignorar comando vazio
		return
	case "!help", "!h":
		// Responder com a lista de comandos disponíveis
		bot.SendMessageToGroup(client, msg, helpMessage)
	case "!ping", "!p":
		// Responder com 'pong'
		bot.SendMessageToGroup(client, msg, pongMessage)
	case "!ram", "!mem":
		// Responder com a quantidade de memória usada
		bot.SendMessageToGroup(client, msg, ramMessage())
	default:
		// Comando desconhecido/invalido
		bot.SendMessageToGroup(client, msg, unknownMessage)
		return
	}
}
