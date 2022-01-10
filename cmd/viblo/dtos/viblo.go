package dtos

import (
	"strings"
	"time"
)

type VibloDate time.Time

func (d *VibloDate) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), `"`)
	t, err := time.Parse("2006-01-02 15:04:05", s)
	if err != nil {
		return err
	}
	*d = VibloDate(t)
	return nil
}

type VibloUser struct {
	ID             int64  `json:"id"`
	Url            string `json:"url"`
	Avatar         string `json:"avatar"`
	Name           string `json:"name"`
	Username       string `json:"username"`
	FollowersCount int64  `json:"followers_count"`
	Reputation     int64  `json:"reputation"`
	PostsCount     int64  `json:"posts_count"`
}

type VibloTag struct {
	Slug    string `json:"slug"`
	Name    string `json:"name"`
	Primary bool   `json:"primary"`
	Image   string `json:"image"`
}

type VibloPost struct {
	ID            int64     `json:"id"`
	Title         string    `json:"title"`
	Slug          string    `json:"slug"`
	Url           string    `json:"url"`
	UserID        int64     `json:"user_id"`
	ContentsShort string    `json:"contents_short"`
	Contents      string    `json:"contents"`
	PublishedAt   VibloDate `json:"published_at"`
	UpdatedAt     VibloDate `json:"updated_at"`
	ReadingTime   int64     `json:"reading_time"`
	Points        int64     `json:"points"`
	ViewsCount    int64     `json:"views_count"`
	ClipsCount    int64     `json:"clips_count"`
	CommentsCount int64     `json:"comments_count"`
	ThumbnailUrl  string    `json:"thumbnail_url"`
	User          struct {
		Data VibloUser `json:"data"`
	} `json:"user"`
	Tags struct {
		Data []VibloTag `json:"data"`
	} `json:"tags"`
}

type VibloPostResponse struct {
	Data []VibloPost `json:"data"`
}
