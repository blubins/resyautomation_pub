package cli

import (
	"resysniper/src/utils"
)

func MainMenu() {
	var input string
	for input != "0" {
		utils.ClearConsole()
		PrintTitle()
		PrintCliOptions()
		input = GetUserInputAndPromptHeader()
		handleUserInput(input)
	}
}

func WebhookMenu() {
	var input string
	for {
		utils.ClearConsole()
		PrintTitle()
		PrintWebhookHeader()
		PrintWebhookOptions()
		input = GetUserInputAndPromptHeader()
		handleUserWebhookInput(input)
		if input == "0" {
			return
		}
	}
}
