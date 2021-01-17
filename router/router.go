package router

import (
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/owncast/owncast/config"
	"github.com/owncast/owncast/controllers"
	"github.com/owncast/owncast/controllers/admin"
	"github.com/owncast/owncast/core/chat"
	"github.com/owncast/owncast/router/middleware"
	"github.com/owncast/owncast/yp"
)

// Start starts the router for the http, ws, and rtmp.
func Start() error {
	// static files
	http.HandleFunc("/", controllers.IndexHandler)

	// admin static files
	http.HandleFunc("/admin/", middleware.RequireAdminAuth(admin.ServeAdmin))

	// status of the system
	http.HandleFunc("/api/status", controllers.GetStatus)

	// custom emoji supported in the chat
	http.HandleFunc("/api/emoji", controllers.GetCustomEmoji)

	// websocket chat server
	go func() {
		err := chat.Start()
		if err != nil {
			log.Fatalln(err)
		}
	}()

	// chat rest api
	http.HandleFunc("/api/chat", controllers.GetChatMessages)

	// web config api
	http.HandleFunc("/api/config", controllers.GetWebConfig)

	// chat embed
	http.HandleFunc("/embed/chat", controllers.GetChatEmbed)

	// video embed
	http.HandleFunc("/embed/video", controllers.GetVideoEmbed)

	http.HandleFunc("/api/yp", yp.GetYPResponse)

	// Authenticated admin requests

	// Current inbound broadcaster
	http.HandleFunc("/api/admin/status", middleware.RequireAdminAuth(admin.Status))

	// Disconnect inbound stream
	http.HandleFunc("/api/admin/disconnect", middleware.RequireAdminAuth(admin.DisconnectInboundConnection))

	// Change the current streaming key in memory
	http.HandleFunc("/api/admin/changekey", middleware.RequireAdminAuth(admin.ChangeStreamKey))

	// Change the current streaming name in memory
	http.HandleFunc("/api/admin/changename", middleware.RequireAdminAuth(admin.ChangeStreamName))

	// Change the current streaming name in memory
	http.HandleFunc("/api/admin/changetitle", middleware.RequireAdminAuth(admin.ChangeStreamTitle))

	// Change the current streaming name in memory
	http.HandleFunc("/api/admin/changetags", middleware.RequireAdminAuth(admin.ChangeStreamTags))

	// Change the extra page content in memory
	http.HandleFunc("/api/admin/changeextrapagecontent", middleware.RequireAdminAuth(admin.ChangeExtraPageContent))

	// Server config
	http.HandleFunc("/api/admin/serverconfig", middleware.RequireAdminAuth(admin.GetServerConfig))

	// Get viewer count over time
	http.HandleFunc("/api/admin/viewersOverTime", middleware.RequireAdminAuth(admin.GetViewersOverTime))

	// Get hardware stats
	http.HandleFunc("/api/admin/hardwarestats", middleware.RequireAdminAuth(admin.GetHardwareStats))

	// Get a a detailed list of currently connected clients
	http.HandleFunc("/api/admin/clients", middleware.RequireAdminAuth(controllers.GetConnectedClients))

	// Get all logs
	http.HandleFunc("/api/admin/logs", middleware.RequireAdminAuth(admin.GetLogs))

	// Get warning/error logs
	http.HandleFunc("/api/admin/logs/warnings", middleware.RequireAdminAuth(admin.GetWarnings))

	// Get all chat messages for the admin, unfiltered.
	http.HandleFunc("/api/admin/chat/messages", middleware.RequireAdminAuth(admin.GetChatMessages))

	// Update chat message visibilty
	http.HandleFunc("/api/admin/chat/updatemessagevisibility", middleware.RequireAdminAuth(admin.UpdateMessageVisibility))

	port := config.Config.GetPublicWebServerPort()

	log.Tracef("Web server running on port: %d", port)

	return http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}
