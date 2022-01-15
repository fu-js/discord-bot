package services

import (
	"encoding/xml"
	"errors"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/fu-js/discord-bot/cmd/viblo/dtos"
	"github.com/fu-js/discord-bot/pkg/utils/log"
	"github.com/fu-js/discord-bot/pkg/utils/math"
	"net/http"
	"time"
)

var colors = []int{
	0xfff100,
	0xff8c00,
	0xe81123,
	0xec008c,
	0x68217a,
	0x00188f,
	0x00bcf2,
	0x00b294,
	0x009e49,
	0xbad80a,
}

type VibloService interface {
	GetEditorChoices(limit int) ([]dtos.VibloRSSItem, error)
	GetTrending(limit int) ([]dtos.VibloRSSItem, error)
	SendMessage(session *discordgo.Session, channelID string, message string) error
	SendPost(session *discordgo.Session, channelID string, posts []dtos.VibloRSSItem) []error
}

type vibloService struct {
}

func NewVibloService() VibloService {
	return &vibloService{}
}

func (s *vibloService) GetEditorChoices(limit int) ([]dtos.VibloRSSItem, error) {
	data := dtos.VibloRSS{}

	resp, err := http.Get("https://viblo.asia/rss/posts/editors-choice.rss")
	if err != nil {
		log.Zap.Errorw("error when call viblo editor choices rss", "error", err)
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		err := errors.New(fmt.Sprintf("http request with non ok status: %v", resp.StatusCode))
		log.Zap.Errorw("error when call viblo editor choices api", "error", err)
		return nil, err
	}
	if err := xml.NewDecoder(resp.Body).Decode(&data); err != nil {
		log.Zap.Errorw("error when decode viblo editor choices response", "error", err)
		return nil, err
	}
	limit = math.MinInt(len(data.Channel.Item), limit)
	return data.Channel.Item[:limit], nil
}

func (s *vibloService) GetTrending(limit int) ([]dtos.VibloRSSItem, error) {
	data := dtos.VibloRSS{}

	resp, err := http.Get("https://viblo.asia/rss/posts/trending.rss")
	if err != nil {
		log.Zap.Errorw("error when call viblo trending rss", "error", err)
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		err := errors.New(fmt.Sprintf("http request with non ok status: %v", resp.StatusCode))
		log.Zap.Errorw("error when call viblo trending api", "error", err)
		return nil, err
	}
	if err := xml.NewDecoder(resp.Body).Decode(&data); err != nil {
		log.Zap.Errorw("error when decode viblo trending response", "error", err)
		return nil, err
	}
	limit = math.MinInt(len(data.Channel.Item), limit)
	return data.Channel.Item[:limit], nil
}

func (s *vibloService) SendMessage(session *discordgo.Session, channelID string, message string) error {
	if _, err := session.ChannelMessageSend(channelID, message); err != nil {
		log.Zap.Errorw("error when send text msg to channel",
			"error", err,
			"channel_id", channelID,
			"message", message,
		)
		return err
	}
	return nil
}

func (s *vibloService) SendPost(session *discordgo.Session, channelID string, posts []dtos.VibloRSSItem) []error {
	errs := make([]error, 0, len(posts))
	for i, post := range posts {
		pubDate, _ := time.Parse("2006-01-02 15:04:05", post.PubDate)
		msg := &discordgo.MessageEmbed{
			URL:         post.Link,
			Type:        discordgo.EmbedTypeArticle,
			Title:       post.Title,
			Timestamp:   pubDate.Format(time.RFC3339),
			Color:       colors[i%10],
			Description: post.Description,
			Author: &discordgo.MessageEmbedAuthor{
				Name: post.Creator.Text,
			},
		}
		if post.Category != "" {
			msg.Fields = append(msg.Fields, &discordgo.MessageEmbedField{
				Name:   "Category",
				Value:  post.Category,
				Inline: true,
			})
		}
		_, err := session.ChannelMessageSendEmbed(channelID, msg)
		log.Zap.Infow("send message", "#", i, "channel_id", channelID, "error", err)
		errs = append(errs, err)
	}
	return errs
}
