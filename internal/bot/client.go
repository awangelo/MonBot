package bot

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"
	"google.golang.org/protobuf/proto"
)

// Declaração da variável global
var groupJID types.JID

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
		sendMessageToGroup(client,
			`Comandos disponíveis:
		!help: Mostra esta mensagem
		!ping: Responde com 'pong'`)
	case "!ping", "!p":
		// Responder com 'pong'
		sendMessageToGroup(client, "pong 🏓")
	default:
		// Comando desconhecido
		sendMessageToGroup(client, "Comando desconhecido")
	}
}

func ReplyToMention(client *whatsmeow.Client, msg *events.Message) {
	replyMsg := &waE2E.Message{
		ExtendedTextMessage: &waE2E.ExtendedTextMessage{
			Text: proto.String("Veja os comandos disponíveis digitando !help"),
			ContextInfo: &waE2E.ContextInfo{
				StanzaID:      &msg.Info.ID,
				Participant:   proto.String(msg.Info.Sender.String()),
				QuotedMessage: msg.Message,
				MentionedJID:  []string{msg.Info.Sender.String()},
			},
		},
	}

	_, err := client.SendMessage(context.Background(), msg.Info.Chat, replyMsg)
	if err != nil {
		log.Printf("Erro ao responder à menção: %v", err)
	} else {
		fmt.Println("Resposta enviada com sucesso!")
	}
}

func SendScheduledMessage(client *whatsmeow.Client, delayMins int) {
	for {
		err := sendMessageToGroup(client, fmt.Sprintf("Mensagem"))
		if err != nil {
			log.Fatalf("Erro ao enviar mensagem: %v", err)
		}
		time.Sleep(time.Minute * time.Duration(delayMins))
	}
}

func sendMessageToGroup(client *whatsmeow.Client, message string) error {
	msg := &waE2E.Message{
		Conversation: proto.String(message),
	}

	_, err := client.SendMessage(context.Background(), groupJID, msg)
	if err != nil {
		return fmt.Errorf("Erro ao enviar mensagem: %v", err)
	}

	log.Println("Mensagem enviada com sucesso!")
	return nil
}

func InitBot() (*whatsmeow.Client, int, error) {
	var client *whatsmeow.Client
	var delayMins int
	// Container envolve um storage sqlite
	container, err := sqlstore.New("sqlite3", "file:whatsmeow.db?_foreign_keys=true", nil)
	if err != nil {
		return client, delayMins, fmt.Errorf("Error creating container: %w", err)
	}
	// Client que ira interagir com a WhatsApp web API
	client = createClient(container)
	groupJID, delayMins, err = loadFromEnv()
	if err != nil {
		return client, delayMins, fmt.Errorf("Error creating new bot: %w", err)
	}
	return client, delayMins, nil
}

func createClient(container *sqlstore.Container) *whatsmeow.Client {
	// Criar um novo dispositivo unico e client
	device, err := container.GetFirstDevice()
	if err != nil {
		log.Fatalf("Error getting first device: %v", err)
	}
	return whatsmeow.NewClient(device, nil)
}

func loadFromEnv() (types.JID, int, error) {
	var groupJID types.JID
	var delayMins int
	if err := godotenv.Load(); err != nil {
		return groupJID, 0, err
	}
	// Jid do grupo deve ser guardado em .env
	groupJIDString := os.Getenv("GROUP_JID")
	if groupJIDString == "" {
		return groupJID, 0, fmt.Errorf("GROUP_JID is not defined in .env")
	}
	// Parse do JID
	groupJID, err := types.ParseJID(groupJIDString)
	if err != nil {
		return groupJID, 0, err
	}
	// Delay de minutos entre mensagens automaticas
	delayMinsStr := os.Getenv("DELAY")
	delayMins, err = strconv.Atoi(delayMinsStr)
	if err != nil {
		return groupJID, 0, fmt.Errorf("Invalid DELAY value: %v", err)
	}
	return groupJID, delayMins, nil
}
