package http

import (
	"github.com/khivuksergey/portmonetka.wallet/config"
	"github.com/khivuksergey/portmonetka.wallet/docs"
	"github.com/khivuksergey/portmonetka.wallet/internal/core/port/service"
	"github.com/khivuksergey/webserver/logger"
	"github.com/khivuksergey/webserver/router"
	"github.com/labstack/echo/v4"
	"net/http"
)

type Router struct {
	*echo.Echo
}

func NewRouter(cfg *config.Configuration, services *service.Manager, logger logger.Logger) http.Handler {
	handlers := newHandlers(services, logger)

	e := router.NewEchoRouter().
		WithConfig(cfg.Router).
		UseMiddleware(handlers.error.HandleError).
		UseHealthCheck().
		UseSwagger(docs.SwaggerInfo, cfg.Swagger)

	wallets := e.Group("users/:userId/wallets", handlers.authentication.AuthenticateJWT)
	wallets.GET("", handlers.wallet.GetWallets)
	wallets.POST("", handlers.wallet.CreateWallet)
	wallets.DELETE("/:walletId", handlers.wallet.DeleteWallet)
	wallets.PATCH("/:walletId", handlers.wallet.UpdateWallet)

	return e
}
