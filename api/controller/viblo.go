package controller

import (
	"github.com/fu-js/discord-bot/cmd/viblo/bot"
	"github.com/fu-js/discord-bot/cmd/viblo/dtos"
	"github.com/labstack/echo/v4"
)

type VibloController interface {
	BaseController
}

type vibloController struct {
	bot bot.VibloBot
}

func NewVibloController(bot bot.VibloBot) VibloController {
	return &vibloController{
		bot: bot,
	}
}

func (v *vibloController) Register(g *echo.Group) {
	group := g.Group("/viblo")
	group.GET("/news/publish", v.PublishNews)
}

func (v *vibloController) PublishNews(c echo.Context) error {
	req := dtos.VibloBotPublishRequest{}
	if err := c.Bind(&req); err != nil {
		return err
	}
	if err := v.bot.PostTrending(req.ChannelID, req.Limit); err != nil {
		return c.JSON(500, dtos.Response{
			Message: err.Error(),
		})
	}
	return c.JSON(200, dtos.Response{
		Message: "success",
	})
}
