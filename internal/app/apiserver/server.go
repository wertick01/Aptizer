package apiserver

import (
	"context"
	"database/sql"
	"net/http"

	"aptizer.com/internal/app/handlers"
	"aptizer.com/internal/app/handlers/middleware"
	"aptizer.com/internal/pkg/logger"
	"aptizer.com/internal/pkg/swagger"
	"github.com/gin-gonic/gin"
)

type Server struct {
	Handlers *handlers.Handler
	Router   *gin.Engine
	db       *sql.DB
	ctx      context.Context
	Logger   logger.ILogger
}

// CreateRoutes - Addresses of function calls.
func (server *Server) CreateRoutes() {
	server.Router.GET("/swagger/*any", swagger.GetGinHandler())
	server.Router.Use(logger.GinMiddlewareLog(server.Logger, true))

	public := server.Router.Group("/")

	public.POST("/api/login", server.Handlers.Authorizer.Login)
	public.POST("/api/refresh", server.Handlers.Authorizer.Refresh)
	public.POST("/api/registration", server.Handlers.Create)

	public.GET("/api/users", server.Handlers.List)
	public.GET("/api/users/:id", server.Handlers.Find)

	private := server.Router.Group("/")
	private.Use(middleware.CheckJWTToken())

	private.PUT("/api/users", server.Handlers.Change)
	private.GET("/api/users/phone", server.Handlers.FindByPhone)
	private.DELETE("/api/users/:id", server.Handlers.Delete)

	server.Router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"code": http.StatusNotFound, "message": "Page not found"})
	})
}
