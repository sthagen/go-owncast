package chat

import (
	"errors"
	"time"

	"github.com/owncast/owncast/models"
)

// Setup sets up the chat server.
func Setup(listener models.ChatListener) {
	setupPersistence()

	clients := make(map[string]*Client)
	addCh := make(chan *Client)
	delCh := make(chan *Client)
	sendAllCh := make(chan models.ChatEvent)
	pingCh := make(chan models.PingMessage)
	doneCh := make(chan bool)
	errCh := make(chan error)

	_server = &server{
		clients,
		"/entry", //hardcoded due to the UI requiring this and it is not configurable
		listener,
		addCh,
		delCh,
		sendAllCh,
		pingCh,
		doneCh,
		errCh,
	}
}

// Start starts the chat server.
func Start() error {
	if _server == nil {
		return errors.New("chat server is nil")
	}

	ticker := time.NewTicker(30 * time.Second)
	go func() {
		for range ticker.C {
			_server.ping()
		}
	}()

	_server.Listen()

	return errors.New("chat server failed to start")
}

// SendMessage sends a message to all.
func SendMessage(message models.ChatEvent) {
	if _server == nil {
		return
	}

	_server.SendToAll(message)
}

// GetMessages gets all of the messages.
func GetMessages() []models.ChatEvent {
	if _server == nil {
		return []models.ChatEvent{}
	}

	return getChatHistory()
}

func GetModerationChatMessages() []models.ChatEvent {
	return getChatModerationHistory()
}

func GetClient(clientID string) *Client {
	for _, client := range _server.Clients {
		if client.ClientID == clientID {
			return client
		}
	}
	return nil
}
