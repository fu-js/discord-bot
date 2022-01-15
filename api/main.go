package main

import (
	"github.com/bwmarrin/discordgo"
	"github.com/fu-js/discord-bot/api/config"
	"github.com/fu-js/discord-bot/api/controller"
	"github.com/fu-js/discord-bot/cmd/viblo/bot"
	"github.com/fu-js/discord-bot/cmd/viblo/services"
	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()

	api := e.Group("/api")

	dc, err := discordgo.New("Bot " + config.Global.Viblo.Token)
	if err != nil {
		panic(err)
	}

	vService := services.NewVibloService()
	vBot := bot.NewNewBot(dc, vService)

	vController := controller.NewVibloController(vBot)
	vController.Register(api)

	e.Logger.Fatal(e.Start(":" + config.Global.Port))
}
