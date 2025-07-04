package main

import (
	"context"
	"log"
	"peekaping/docs"
	"peekaping/src/modules/auth"
	"peekaping/src/modules/cleanup"
	"peekaping/src/modules/events"
	"peekaping/src/modules/healthcheck"
	"peekaping/src/modules/heartbeat"
	"peekaping/src/modules/maintenance"
	"peekaping/src/modules/monitor"
	"peekaping/src/modules/monitor_notification"
	"peekaping/src/modules/monitor_status_page"
	"peekaping/src/modules/notification_channel"
	"peekaping/src/modules/proxy"
	"peekaping/src/modules/setting"
	"peekaping/src/modules/stats"
	"peekaping/src/modules/status_page"
	"peekaping/src/modules/websocket"
	"peekaping/src/utils"

	"go.uber.org/dig"
	"go.uber.org/zap"
)

// @title			Peekaping API
// @BasePath	/api/v1
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	docs.SwaggerInfo.Version = Version

	utils.RegisterCustomValidators()

	container := dig.New()

	// Provide dependencies
	container.Provide(ProvideConfig)
	container.Provide(ProvideLogger)
	container.Provide(ProvideServer)
	container.Provide(ProvideMongoDB)
	container.Provide(websocket.NewServer)

	// Register dependencies in the correct order to handle circular dependencies
	events.RegisterDependencies(container)
	heartbeat.RegisterDependencies(container)
	monitor.RegisterDependencies(container)
	healthcheck.RegisterDependencies(container)
	auth.RegisterDependencies(container)
	notification_channel.RegisterDependencies(container)
	monitor_notification.RegisterDependencies(container)
	proxy.RegisterDependencies(container)
	setting.RegisterDependencies(container)
	stats.RegisterDependencies(container)
	maintenance.RegisterDependencies(container)
	status_page.RegisterDependencies(container)
	monitor_status_page.RegisterDependencies(container)

	// Start the event healthcheck listener
	err := container.Invoke(func(listener *healthcheck.EventListener, eventBus *events.EventBus) {
		listener.Start(eventBus)
	})

	if err != nil {
		log.Fatal(err)
	}

	// Start cleanup cron job(s)
	err = container.Invoke(func(heartbeatService heartbeat.Service, settingService setting.Service, logger *zap.SugaredLogger) {
		cleanup.StartCleanupCron(heartbeatService, settingService, logger)
	})
	if err != nil {
		log.Fatal(err)
	}

	// Start the health check supervisor
	err = container.Invoke(func(supervisor *healthcheck.HealthCheckSupervisor) {
		if err := supervisor.StartAll(context.Background()); err != nil {
			log.Fatal(err)
		}
	})
	if err != nil {
		log.Fatal(err)
	}

	err = container.Invoke(func(listener *notification_channel.NotificationEventListener, eventBus *events.EventBus) {
		listener.Subscribe(eventBus)
	})
	if err != nil {
		log.Fatal(err)
	}

	// Start the server
	err = container.Invoke(func(server *Server) {
		docs.SwaggerInfo.Host = "localhost:" + server.cfg.Port

		port := server.cfg.Port
		if port == "" {
			port = "8080"
		}
		if port[0] != ':' {
			port = ":" + port
		}
		if err := server.router.Run(port); err != nil {
			log.Fatal(err)
		}
	})

	if err != nil {
		log.Fatal(err)
	}
}
