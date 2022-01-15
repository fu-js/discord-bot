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

type VibloRSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	Creator     struct {
		Text string `xml:",chardata"`
		Dc   string `xml:"dc,attr"`
	} `xml:"creator"`
	PubDate  string `xml:"pubDate"`
	Category string `xml:"category"`
}

type VibloRSS struct {
	Channel struct {
		Item []VibloRSSItem `xml:"item"`
	} `xml:"channel"`
}
