package viblo

import (
	"github.com/bwmarrin/discordgo"
	"github.com/fu-js/discord-bot/cmd/viblo/bot"
	"github.com/fu-js/discord-bot/cmd/viblo/services"
	"github.com/fu-js/discord-bot/pkg/utils/log"
	"github.com/spf13/cobra"
	"os"
)

var (
	token     string
	channelID string
	option    string
	limit     int
)

var Cmd = &cobra.Command{
	Use:   "viblo",
	Short: "viblo editor choice / trending",
	Run: func(cmd *cobra.Command, args []string) {
		dc, err := discordgo.New("Bot " + token)
		if err != nil {
			log.Zap.Errorw("error when create bot", "error", err)
			os.Exit(1)
		}
		service := services.NewVibloService()

		b := bot.NewNewBot(dc, service)

		switch option {
		case "EC":
			if err := b.PostEditorChoice(channelID, limit); err != nil {
				log.Zap.Errorw("error when post editor choices", "error", err)
			}
		case "T":
			if err := b.PostTrending(channelID, limit); err != nil {
				log.Zap.Errorw("error when post trending", "error", err)
			}
		}
	},
}

func init() {
	Cmd.Flags().StringVarP(&token, "token", "t", os.Getenv("TOKEN"), "bot token")
	Cmd.Flags().StringVarP(&channelID, "channel", "c", os.Getenv("CHANNEL_ID"), "channel id")
	Cmd.Flags().StringVarP(&option, "option", "o", "EC", "EC for editor choices, T for trending")
	Cmd.Flags().IntVarP(&limit, "limit", "l", 10, "Number of post")
	Cmd.MarkFlagRequired("token")
	Cmd.MarkFlagDirname("channel")

}
