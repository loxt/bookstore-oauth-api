package app

import (
	"github.com/gin-gonic/gin"
	"github.com/loxt/bookstore-oauth-api/internal/domain/access_token"
	"github.com/loxt/bookstore-oauth-api/internal/http"
	"github.com/loxt/bookstore-oauth-api/internal/repository/db"
)

var (
	router = gin.Default()
)

func StartApplication() {
	dbRepository := db.New()
	service := access_token.NewService(dbRepository)
	atHandler := http.NewHandler(service)

	router.GET("/oauth/access_token/:access_token_id", atHandler.GetById)
	router.POST("/oauth/access_token", atHandler.Create)

	router.Run(":8080")
}
