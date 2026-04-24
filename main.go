package main

import (
	"fmt"
	"log"
	"log/slog"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"

	docs "vrcomandaapi/docs"
	"vrcomandaapi/internal/config"
	"vrcomandaapi/internal/database"
	"vrcomandaapi/internal/modules/atendente"
	"vrcomandaapi/internal/modules/comanda"
	"vrcomandaapi/internal/modules/global"
	"vrcomandaapi/internal/modules/lancamento"
	"vrcomandaapi/internal/modules/mesa"
	"vrcomandaapi/internal/shared/middleware"
	"vrcomandaapi/internal/shared/models"
	"vrcomandaapi/internal/shared/utils"
)

var fatalLog = log.Fatal

var runServer = func(router *gin.Engine, port string) error {
	return router.Run(port)
}

var autoMigrateSQLite = func(db *gorm.DB) error {
	return db.AutoMigrate(
		&models.LancamentoComanda{},
		&models.LancamentoComandaItem{},
	)
}

// @title VRComandaAPI
// @version 1.0
// @description API backend para operacao de comandas.
// @host localhost:8080
// @BasePath /
// @schemes http

// bootstrap wires infrastructure and routes in one place.
//
// The application is organized by domain modules so each feature owns its
// handler, service, and repository. This keeps boundaries explicit, reduces
// coupling, and makes future changes easier to localize.
func bootstrap() (*gin.Engine, error) {
	_ = godotenv.Load()
	cfg := config.Load()
	docs.SwaggerInfo.Host = swaggerHostFromPort(cfg.HTTPPort)
	docs.SwaggerInfo.Schemes = []string{"http", "https"}

	sqliteDB, err := database.ConnectSQLite(cfg)
	if err != nil {
		return nil, fmt.Errorf("erro ao conectar no SQLite: %w", err)
	}

	mongoManager := database.NewMongoManager(cfg)
	mongoManager.Start(30 * time.Second)

	if err := autoMigrateSQLite(sqliteDB); err != nil {
		return nil, fmt.Errorf("erro no auto-migrate do SQLite: %w", err)
	}

	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(middleware.Logger(), middleware.Recovery())
	router.GET("/health", utils.HealthHandler)
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Each module receives only the dependency it needs.
	comanda.RegisterRoutes(router, mongoManager.DB, mongoManager.InvalidateConnection)
	mesa.RegisterRoutes(router, mongoManager.DB, mongoManager.InvalidateConnection)
	atendente.RegisterRoutes(router, mongoManager.DB, mongoManager.InvalidateConnection)

	lancamento.RegisterRoutes(router, sqliteDB)

	global.RegisterRoutes(
		router,
		lancamento.NewService(lancamento.NewRepository(sqliteDB)),
		comanda.NewService(comanda.NewMongoRepository(mongoManager.DB, mongoManager.InvalidateConnection, "comandas")),
		mesa.NewService(mesa.NewMongoRepository(mongoManager.DB, mongoManager.InvalidateConnection, "mesas")),
	)

	return router, nil
}

// main is the process entrypoint.
func main() {
	if err := run(bootstrap, config.Load); err != nil {
		fatalLog(err)
	}
}

func run(bootstrapFn func() (*gin.Engine, error), loadConfigFn func() config.Config) error {
	router, err := bootstrapFn()
	if err != nil {
		return err
	}

	port := loadConfigFn().HTTPPort
	slog.Info("VRComandaAPI starting", "port", port)

	if err := runServer(router, port); err != nil {
		return fmt.Errorf("erro ao iniciar servidor HTTP na porta %s: %w", port, err)
	}

	return nil
}

func swaggerHostFromPort(port string) string {
	trimmed := strings.TrimSpace(port)
	if trimmed == "" {
		return "localhost:8080"
	}

	if strings.HasPrefix(trimmed, ":") {
		return "localhost" + trimmed
	}

	if strings.Contains(trimmed, ":") {
		return trimmed
	}

	return "localhost:" + trimmed
}
