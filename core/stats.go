package core

import (
	"math"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/owncast/owncast/core/chat"
	"github.com/owncast/owncast/core/data"
	"github.com/owncast/owncast/geoip"
	"github.com/owncast/owncast/models"
	"github.com/owncast/owncast/utils"
)

var l = &sync.RWMutex{}

func setupStats() error {
	s := getSavedStats()
	_stats = &s

	statsSaveTimer := time.NewTicker(1 * time.Minute)
	go func() {
		for range statsSaveTimer.C {
			if err := saveStats(); err != nil {
				panic(err)
			}
		}
	}()

	return nil
}

// IsStreamConnected checks if the stream is connected or not.
func IsStreamConnected() bool {
	if !_stats.StreamConnected {
		return false
	}

	// Kind of a hack.  It takes a handful of seconds between a RTMP connection and when HLS data is available.
	// So account for that with an artificial buffer of four segments.
	timeSinceLastConnected := time.Since(_stats.LastConnectTime.Time).Seconds()
	waitTime := math.Max(float64(data.GetStreamLatencyLevel().SecondsPerSegment)*3.0, 7)
	if timeSinceLastConnected < waitTime {
		return false
	}

	return _stats.StreamConnected
}

// SetClientActive sets a client as active and connected.
func SetClientActive(client models.Client) {
	l.Lock()
	defer l.Unlock()

	// If this clientID already exists then update it.
	// Otherwise set a new one.
	if existingClient, ok := _stats.Clients[client.ClientID]; ok {
		existingClient.LastSeen = time.Now()
		existingClient.Username = client.Username
		existingClient.MessageCount = client.MessageCount
		existingClient.Geo = geoip.GetGeoFromIP(existingClient.IPAddress)
		_stats.Clients[client.ClientID] = existingClient
	} else {
		if client.Geo == nil {
			geoip.FetchGeoForIP(client.IPAddress)
		}
		_stats.Clients[client.ClientID] = client
	}

	// Don't update viewer counts if a live stream session is not active.
	if _stats.StreamConnected {
		_stats.SessionMaxViewerCount = int(math.Max(float64(len(_stats.Clients)), float64(_stats.SessionMaxViewerCount)))
		_stats.OverallMaxViewerCount = int(math.Max(float64(_stats.SessionMaxViewerCount), float64(_stats.OverallMaxViewerCount)))
	}
}

// RemoveClient removes a client from the active clients record.
func RemoveClient(clientID string) {
	log.Trace("Removing the client:", clientID)

	l.Lock()
	delete(_stats.Clients, clientID)
	l.Unlock()
}

func GetClients() []models.Client {
	l.RLock()
	clients := make([]models.Client, 0)
	for _, client := range _stats.Clients {
		chatClient := chat.GetClient(client.ClientID)
		if chatClient != nil {
			clients = append(clients, chatClient.GetViewerClientFromChatClient())
		} else {
			clients = append(clients, client)
		}
	}
	l.RUnlock()

	return clients
}

func saveStats() error {
	if err := data.SetPeakOverallViewerCount(_stats.OverallMaxViewerCount); err != nil {
		log.Errorln("error saving viewer count", err)
	}
	if err := data.SetPeakSessionViewerCount(_stats.SessionMaxViewerCount); err != nil {
		log.Errorln("error saving viewer count", err)
	}
	if err := data.SetLastDisconnectTime(_stats.LastConnectTime.Time); err != nil {
		log.Errorln("error saving disconnect time", err)
	}

	return nil
}

func getSavedStats() models.Stats {
	savedLastDisconnectTime, savedLastDisconnectTimeErr := data.GetLastDisconnectTime()
	result := models.Stats{
		Clients:               make(map[string]models.Client),
		SessionMaxViewerCount: data.GetPeakSessionViewerCount(),
		OverallMaxViewerCount: data.GetPeakOverallViewerCount(),
		LastDisconnectTime:    utils.NullTime{Time: savedLastDisconnectTime, Valid: savedLastDisconnectTimeErr == nil},
	}

	// If the stats were saved > 5min ago then ignore the
	// peak session count value, since the session is over.
	if !result.LastDisconnectTime.Valid || time.Since(result.LastDisconnectTime.Time).Minutes() > 5 {
		result.SessionMaxViewerCount = 0
	}

	return result
}
