package utils

// config/config.json file struct for unmarshal
type UserConfigJson struct {
	License    string `json:"license"`
	WebhookUrl string `json:"webhookurl"`
}
