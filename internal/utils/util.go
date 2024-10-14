package utils

import (
	"log"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"go.mau.fi/whatsmeow/types/events"
)

// LogMessageEvent loga mensagens recebidas
func LogMessageEvent(msg *events.Message) {
	sender := msg.Info.Sender.String()

	// Verificar diferentes tipos de mensagens
	var text string
	if msg.Message.Conversation != nil {
		text = msg.Message.GetConversation()
	} else if msg.Message.ExtendedTextMessage != nil {
		text = msg.Message.ExtendedTextMessage.GetText()
	} else {
		// Outros tipos de mensagens (mídia, etc.)
		text = "media message"
	}

	log.Printf("Mensagem recebida no grupo de %s: %s\n", sender, text)
}

// PreventExit previne a saída do programa
func PreventExit() {
	// Listen to Ctrl+C (you can also do something else that prevents the program from exiting)
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
}

// GetMemoryUsage retorna a quantidade de memória usada em KB
func GetMemoryUsage() uint64 {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	// Converter para KB
	return m.Alloc / 1024
}
