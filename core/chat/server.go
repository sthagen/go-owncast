package chat

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
	"golang.org/x/net/websocket"

	"github.com/owncast/owncast/core/data"
	"github.com/owncast/owncast/core/webhooks"
	"github.com/owncast/owncast/models"
)

var (
	_server *server
)

var l = &sync.RWMutex{}

// Server represents the server which handles the chat.
type server struct {
	Clients map[string]*Client

	pattern  string
	listener models.ChatListener

	addCh     chan *Client
	delCh     chan *Client
	sendAllCh chan models.ChatEvent
	pingCh    chan models.PingMessage
	doneCh    chan bool
	errCh     chan error
}

// Add adds a client to the server.
func (s *server) add(c *Client) {
	s.addCh <- c
}

// Remove removes a client from the server.
func (s *server) remove(c *Client) {
	s.delCh <- c
}

// SendToAll sends a message to all of the connected clients.
func (s *server) SendToAll(msg models.ChatEvent) {
	s.sendAllCh <- msg
}

// Err handles an error.
func (s *server) err(err error) {
	s.errCh <- err
}

func (s *server) sendAll(msg models.ChatEvent) {
	l.RLock()
	for _, c := range s.Clients {
		c.write(msg)
	}
	l.RUnlock()
}

func (s *server) ping() {
	ping := models.PingMessage{MessageType: models.PING}

	l.RLock()
	for _, c := range s.Clients {
		c.pingch <- ping
	}
	l.RUnlock()
}

func (s *server) usernameChanged(msg models.NameChangeEvent) {
	l.RLock()
	for _, c := range s.Clients {
		c.usernameChangeChannel <- msg
	}
	l.RUnlock()

	go webhooks.SendChatEventUsernameChanged(msg)
}

func (s *server) userJoined(msg models.UserJoinedEvent) {
	l.RLock()
	if s.listener.IsStreamConnected() {
		for _, c := range s.Clients {
			c.userJoinedChannel <- msg
		}
	}
	l.RUnlock()

	go webhooks.SendChatEventUserJoined(msg)
}

func (s *server) onConnection(ws *websocket.Conn) {
	client := NewClient(ws)

	defer func() {
		s.removeClient(client)

		if err := ws.Close(); err != nil {
			log.Debugln(err)
			//s.errCh <- err
		}
	}()

	s.add(client)
	client.listen()
}

// Listen and serve.
// It serves client connection and broadcast request.
func (s *server) Listen() {
	http.Handle(s.pattern, websocket.Handler(s.onConnection))

	log.Tracef("Starting the websocket listener on: %s", s.pattern)

	for {
		select {
		// add new a client
		case c := <-s.addCh:
			l.Lock()
			s.Clients[c.socketID] = c

			if !c.Ignore {
				s.listener.ClientAdded(c.GetViewerClientFromChatClient())
				s.sendWelcomeMessageToClient(c)
			}
			l.Unlock()

		// remove a client
		case c := <-s.delCh:
			s.removeClient(c)
		case msg := <-s.sendAllCh:
			if data.GetChatDisabled() {
				break
			}

			if !msg.Empty() {
				// set defaults before sending msg to anywhere
				msg.SetDefaults()

				s.listener.MessageSent(msg)
				s.sendAll(msg)

				// Store in the message history
				if !msg.Ephemeral {
					addMessage(msg)
				}

				// Send webhooks
				go webhooks.SendChatEvent(msg)
			}
		case ping := <-s.pingCh:
			fmt.Println("PING?", ping)

		case err := <-s.errCh:
			log.Trace("Error: ", err.Error())

		case <-s.doneCh:
			return
		}
	}
}

func (s *server) removeClient(c *Client) {
	l.Lock()
	if _, ok := s.Clients[c.socketID]; ok {
		delete(s.Clients, c.socketID)

		s.listener.ClientRemoved(c.socketID)
		log.Tracef("The client was connected for %s and sent %d messages (%s)", time.Since(c.ConnectedAt), c.MessageCount, c.ClientID)
	}
	l.Unlock()
}

func (s *server) sendWelcomeMessageToClient(c *Client) {
	go func() {
		// Add an artificial delay so people notice this message come in.
		time.Sleep(7 * time.Second)

		welcomeMessage := data.GetServerWelcomeMessage()
		if welcomeMessage != "" {
			initialMessage := models.ChatEvent{ClientID: "owncast-server", Author: data.GetServerName(), Body: welcomeMessage, ID: "initial-message-1", MessageType: "SYSTEM", Visible: true, Timestamp: time.Now()}
			c.write(initialMessage)
		}
	}()
}
