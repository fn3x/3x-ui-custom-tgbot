package controller

import (
	"x-ui/web/service"

	"github.com/gin-gonic/gin"
)

type WebhookServerController struct {
	webhookService service.WebhookService
}

func NewWebhookController(g *gin.RouterGroup) *WebhookServerController {
	a := &WebhookServerController{}
	a.initRouter(g)
	return a
}

func (a *WebhookServerController) initRouter(g *gin.RouterGroup) {
	g.POST("/webhooks", a.webhooks)
}

func (a *WebhookServerController) webhooks(c *gin.Context) {
	a.webhookService.WebhookHandler(c.Writer, c.Request)
}
