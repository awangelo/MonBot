package config

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.mau.fi/whatsmeow/types"
)

// InitBot inicializa o bot com as variáveis de ambiente e banco de dados
func InitBot() (*whatsmeow.Client, types.JID, int, error) {
	var client *whatsmeow.Client
	var delayMins int
	var groupJID types.JID
	// Container envolve um storage sqlite
	container, err := sqlstore.New("sqlite3", "file:whatsmeow.db?_foreign_keys=true", nil)
	if err != nil {
		return client, groupJID, delayMins, fmt.Errorf("Erro ao criar container SQL: %w", err)
	}
	// Client que ira interagir com a API do WhatsApp Web
	client = createClient(container)
	groupJID, delayMins, err = loadFromEnv()
	if err != nil {
		return client, groupJID, delayMins, fmt.Errorf("Erro ao carregar variáveis de ambiente: %w", err)
	}
	return client, groupJID, delayMins, nil
}

func createClient(container *sqlstore.Container) *whatsmeow.Client {
	// Criar um novo dispositivo unico e client
	device, err := container.GetFirstDevice()
	if err != nil {
		log.Fatalf("Erro ao obter o dispositivo: %v", err)
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
		return groupJID, 0, fmt.Errorf("GROUP_JID nao definido no arquivo .env")
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
		return groupJID, 0, fmt.Errorf("Erro ao converter DELAY para inteiro: %w", err)
	}
	return groupJID, delayMins, nil
}
