package bot

import (
	"github.com/fu-js/discord-bot/cmd/viblo/services"
	"github.com/bwmarrin/discordgo"
)

type VibloBot interface {
	PostTrending(channelID string, limit int) error
	PostEditorChoice(channelID string, limit int) error
}

type vibloBot struct {
	session      *discordgo.Session
	vibloService services.VibloService
}

func NewNewBot(
	session *discordgo.Session,
	vibloService services.VibloService,
) VibloBot {
	return &vibloBot{
		session:      session,
		vibloService: vibloService,
	}
}

func (b *vibloBot) PostTrending(channelID string, limit int) error {
	if err := b.session.Open(); err != nil {
		return err
	}
	defer b.session.Close()
	trending, err := b.vibloService.GetTrending(limit)
	if err != nil {
		return err
	}
	b.vibloService.SendPost(b.session, channelID, trending)
	return nil
}

func (b *vibloBot) PostEditorChoice(channelID string, limit int) error {
	if err := b.session.Open(); err != nil {
		return err
	}
	defer b.session.Close()
	editorChoices, err := b.vibloService.GetEditorChoices(limit)
	if err != nil {
		return err
	}
	b.vibloService.SendPost(b.session, channelID, editorChoices)
	return nil
}
