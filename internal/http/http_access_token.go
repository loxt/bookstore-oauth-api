package http

import (
	"github.com/gin-gonic/gin"
	"github.com/loxt/bookstore-oauth-api/internal/domain/access_token"
	"net/http"
	"strings"
)

type AccessTokenHandler interface {
	GetById(c *gin.Context)
}

type accessTokenHandler struct {
	service access_token.Service
}

func NewHandler(service access_token.Service) AccessTokenHandler {
	return &accessTokenHandler{
		service,
	}
}

func (h *accessTokenHandler) GetById(c *gin.Context) {
	accessTokenId := strings.TrimSpace(c.Param("access_token_id"))

	accessToken, err := h.service.GetById(accessTokenId)

	if err != nil {
		c.JSON(err.Status, err)
	}

	c.JSON(http.StatusOK, accessToken)
}
