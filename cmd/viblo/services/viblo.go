package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/fu-js/discord-bot/cmd/viblo/dtos"
	"github.com/fu-js/discord-bot/pkg/utils/log"
	"github.com/headzoo/surf/browser"
	"gopkg.in/headzoo/surf.v1"
	"html"
	"net/http"
	"strings"
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
	GetEditorChoices(limit int) ([]dtos.VibloPost, error)
	GetTrending(limit int) ([]dtos.VibloPost, error)
	SendMessage(session *discordgo.Session, channelID string, message string) error
	SendPost(session *discordgo.Session, channelID string, posts []dtos.VibloPost) []error
}

type vibloService struct {
	surf *browser.Browser
}

func NewVibloService() VibloService {
	b := surf.NewBrowser()

	return &vibloService{
		surf: b,
	}
}

func (s *vibloService) GetEditorChoices(limit int) ([]dtos.VibloPost, error) {
	data := dtos.VibloPostResponse{}
	tab := s.surf.NewTab()
	tab.AddRequestHeader("Accept", "text/json")
	tab.AddRequestHeader("Accept-Charset", "utf8")

	if err := tab.Open(fmt.Sprintf("https://viblo.asia/api/posts/editors-choice?limit=%v", limit)); err != nil {
		return nil, err
	}
	if tab.StatusCode() != http.StatusOK {
		err := errors.New(fmt.Sprintf("http request with non ok status: %v", tab.StatusCode()))
		log.Zap.Errorw("error when call viblo editor choices api", "error", err)
		return nil, err
	}
	if err := json.Unmarshal([]byte(tab.Body()), &data); err != nil {
		log.Zap.Errorw("error when decode viblo editor choices response", "error", err)
		return nil, err
	}
	return data.Data, nil
}

func (s *vibloService) GetTrending(limit int) ([]dtos.VibloPost, error) {
	data := dtos.VibloPostResponse{}

	tab := s.surf.NewTab()

	if err := tab.Open(fmt.Sprintf("https://viblo.asia/api/posts/trending?limit=%v", limit)); err != nil {
		return nil, err
	}

	if tab.StatusCode() != http.StatusOK {
		err := errors.New(fmt.Sprintf("http request with non ok status: %v", tab.StatusCode()))
		log.Zap.Errorw("error when call viblo trending api", "error", err)
		return nil, err
	}
	body := tab.Body()
	body = html.UnescapeString(body)
	if err := json.Unmarshal([]byte(body), &data); err != nil {
		log.Zap.Errorw("error when decode viblo trending response", "error", err)
		return nil, err
	}
	return data.Data, nil
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

func (s *vibloService) SendPost(session *discordgo.Session, channelID string, posts []dtos.VibloPost) []error {
	errs := make([]error, 0, len(posts))
	for i, post := range posts {
		tags := make([]string, 0, len(post.Tags.Data))
		for _, tag := range post.Tags.Data {
			tags = append(tags, tag.Name)
		}
		_, err := session.ChannelMessageSendEmbed(channelID, &discordgo.MessageEmbed{
			URL:       post.Url,
			Type:      discordgo.EmbedTypeArticle,
			Title:     post.Title,
			Timestamp: time.Time(post.PublishedAt).Format(time.RFC3339),
			Color:     colors[i%10],
			Footer: &discordgo.MessageEmbedFooter{
				Text: fmt.Sprintf("👀 %vm read time", post.ReadingTime),
			},
			Thumbnail: &discordgo.MessageEmbedThumbnail{
				URL: post.ThumbnailUrl,
			},
			Video:    nil,
			Provider: nil,
			Author: &discordgo.MessageEmbedAuthor{
				URL:     post.User.Data.Url,
				Name:    fmt.Sprintf("%v(%v)[%v🏆]", post.User.Data.Name, post.User.Data.Username, post.User.Data.Reputation),
				IconURL: fmt.Sprintf("https://images.viblo.asia/avatar/%v", post.User.Data.Avatar),
			},
			Fields: []*discordgo.MessageEmbedField{
				{
					Name:   "Point",
					Value:  fmt.Sprint(post.Points),
					Inline: true,
				},
				{
					Name:   "View",
					Value:  fmt.Sprint(post.ViewsCount),
					Inline: true,
				},
				{
					Name:   "Tags",
					Value:  strings.Join(tags, ", "),
					Inline: true,
				},
			},
		})
		log.Zap.Infow("send message", "#", i, "channel_id", channelID, "error", err)
		errs = append(errs, err)
	}
	return errs
}
