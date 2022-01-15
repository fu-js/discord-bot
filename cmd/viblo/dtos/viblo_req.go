package dtos

type VibloBotPublishRequest struct {
	ChannelID string `json:"channel_id,omitempty" query:"channel_id"`
	Limit     int    `json:"limit,omitempty" query:"limit"`
}
