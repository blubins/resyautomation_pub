package cli

import (
	"bytes"
	"encoding/json"
	"fmt"
	"resysniper/src/colors"
	"resysniper/src/utils"
	"strings"
	"time"

	http "github.com/bogdanfinn/fhttp"
	tls_client "github.com/bogdanfinn/tls-client"
)

// Sends a test webhook with the users set webhook url
// does not care or check really if it is valid url only that it is filled with something
func SendWebhookTest() error {
	userConfig, err := utils.OpenConfigJson("./config/config.json")
	if err != nil {
		return err
	}
	if len(userConfig.WebhookUrl) == 0 {
		return fmt.Errorf("no webhook found")
	}

	title := "Webhook Test"
	webhookPayload := map[string]interface{}{
		"content": "",
		"embeds": []map[string]interface{}{
			{
				"title":       title,
				"url":         nil,
				"description": "```md\nResy - Sniper\n=============\n```",
				"color":       0,
				"thumbnail": map[string]interface{}{
					"url": "https://images-ext-1.discordapp.net/external/Gl0H_fucnFv4rQa8bBvbdvynfyFdY6plQyXVCPCN4oM/https/i.gyazo.com/02d5fbd5dec03e7438f292d9e751f271.png?format=webp&quality=lossless",
				},
				"fields":    nil,
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

	client, err := tls_client.NewHttpClient(tls_client.NewNoopLogger(), tls_client.DefaultOptions...)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", userConfig.WebhookUrl, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return err
	}

	// discord api needs header to be set
	req.Header.Set("content-type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	// status 204 is success
	if resp.StatusCode != 204 {
		return fmt.Errorf("non 200 status received: %d", resp.StatusCode)
	}

	return nil
}

// Prompt for user to edit their webhook saved in /config/config.json
func EditWebhook() error {
	utils.ClearConsole()
	PrintTitle()
	PrintWebhookHeader()
	fmt.Printf("%sEnter your new webhook or \"0\" to return without saving:%s",
		colors.MAGENTA, colors.RESET)
	input := GetUserInput()
	if input == "0" {
		return nil
	}
	err := utils.UpdateConfigJson("./config/config.json", "WebhookUrl", strings.TrimRight(input, "\r\n"))
	if err != nil {
		time.Sleep(time.Second * 3)
		return err
	}
	return nil
}
