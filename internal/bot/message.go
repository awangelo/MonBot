package bot

import (
	"context"
	"fmt"
	"log"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"
	"google.golang.org/protobuf/proto"
)

const (
	infoMessage = "Veja os comandos disponíveis digitando !help"
)

// ReplyToMention responde à menção do bot com a mensagem de ajuda
func ReplyToMention(client *whatsmeow.Client, msg *events.Message) {
	replyMsg := &waE2E.Message{
		ExtendedTextMessage: &waE2E.ExtendedTextMessage{
			Text: proto.String(infoMessage),
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
	}
}

// SendScheduledMessage envia a mensagem de ajuda a cada intervalo de tempo
func SendScheduledMessage(client *whatsmeow.Client, groupJID types.JID, delayMins int) {
	for {
		m := &waE2E.Message{
			Conversation: proto.String(infoMessage),
		}

		_, err := client.SendMessage(context.Background(), groupJID, m)
		if err != nil {
			log.Fatalf("Erro ao enviar a mensagem programada: %v", err)
		}

		time.Sleep(time.Minute * time.Duration(delayMins))
	}
}

// SendMessageToGroup envia uma mensagem para o grupo
func SendMessageToGroup(client *whatsmeow.Client, msg *events.Message, message string) error {
	m := &waE2E.Message{
		Conversation: proto.String(message),
	}

	_, err := client.SendMessage(context.Background(), msg.Info.Chat, m)
	if err != nil {
		return fmt.Errorf("Erro ao enviar mensagem: %v", err)
	}

	log.Println("Mensagem enviada com sucesso!")
	return nil
}
