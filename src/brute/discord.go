package brute

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"

	http "github.com/bogdanfinn/fhttp"
)

func (rt *ResyTask) SendSuccessWebhook() error {

	fields := []map[string]interface{}{
		{"name": "Restaurant ID", "value": rt.RestaurantId, "inline": true},
		{"name": "Restaurant", "value": rt.Venue.Name, "inline": true},
		{"name": "", "value": "", "inline": true},
		{"name": "Date", "value": rt.Venue.Date, "inline": true},
		{"name": "Room", "value": rt.Venue.RoomType, "inline": true},
		{"name": "Party Size", "value": rt.PartySize, "inline": true},
		{"name": "Email", "value": rt.EmailAddress, "inline": true},
		{"name": "Proxy", "value": fmt.Sprintf("||%s:%s:%s:%s", rt.Proxy.Host, rt.Proxy.Port, rt.Proxy.Username, rt.Proxy.Password), "inline": true},
	}

	title := fmt.Sprintf("Successfully booked %s!", rt.Venue.Name)

	webhookPayload := map[string]interface{}{
		"content": "",
		"embeds": []map[string]interface{}{
			{
				"title":       title,
				"url":         nil,
				"description": nil,
				"color":       0,
				"thumbnail": map[string]interface{}{
					"url": "https://images-ext-1.discordapp.net/external/Gl0H_fucnFv4rQa8bBvbdvynfyFdY6plQyXVCPCN4oM/https/i.gyazo.com/02d5fbd5dec03e7438f292d9e751f271.png?format=webp&quality=lossless",
				},
				"fields":    fields,
				"timestamp": time.Now().Format(time.RFC3339),
				"footer": map[string]interface{}{
					"text":     "Resy Sniper",
					"icon_url": "https://cdn.discordapp.com/avatars/851246893493518356/05a9fffd14e2bf44ecd730e1cc666f2d.webp?size=1280",
				},
			},
		},
		"username":    "Resy Sniper",
		"avatar_url":  "https://cdn.discordapp.com/avatars/851246893493518356/05a9fffd14e2bf44ecd730e1cc666f2d.webp?size=1280",
		"attachments": []interface{}{},
	}

	jsonPayload, err := json.Marshal(webhookPayload)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", rt.WebhookUrl, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return err
	}

	// discord api needs header to be set
	req.Header.Set("content-type", "application/json")
	resp, err := rt.Client.Do(req)
	if err != nil {
		return err
	}

	// status 204 is successful webhook sent
	if resp.StatusCode != 204 {
		return fmt.Errorf("non 200 status received: %d", resp.StatusCode)
	}

	return nil
}
