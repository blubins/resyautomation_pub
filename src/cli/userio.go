package cli

import (
	"bufio"
	"fmt"
	"os"
	"resysniper/src/accgen"
	"resysniper/src/brute"
	"resysniper/src/colors"
	"resysniper/src/utils"
	"strings"
)

func GetUserInputAndPromptHeader() string {
	fmt.Printf("\n%sEnter your selection:%s",
		colors.MAGENTA, colors.RESET)
	reader := bufio.NewReader(os.Stdin)
	userInput, _ := reader.ReadString('\n')
	return strings.TrimSpace(userInput)
}

func GetUserInput() string {
	reader := bufio.NewReader(os.Stdin)
	userInput, _ := reader.ReadString('\n')
	return strings.TrimSpace(userInput)
}

func PrintTitle() {
	fmt.Printf("%s%s\n%s",
		colors.MAGENTA, utils.ResyTitle, colors.RESET)
	f, err := os.Stat("./config/logs.txt")
	if err != nil {
		return
	}
	fmt.Printf("%sLogs File Length: %.3f MB%s\n",
		colors.YELLOW, float32(f.Size())/1000000, colors.RESET)
}

func PrintCliOptions() {
	fmt.Print(colors.CYAN)
	fmt.Printf("---------Tasks-------\n")
	fmt.Printf("1. Safe Mode (Not Implemented)\n")
	fmt.Printf("2. Brute Menu\n")
	fmt.Printf("3. Start Account Generator\n")
	fmt.Printf("----Configuration----\n")
	fmt.Printf("4. Webhook Settings\n")
	fmt.Printf("5. Open reservation settings\n")
	fmt.Printf("6. Open generator settings\n")
	fmt.Printf("7. Open accounts list\n")
	fmt.Printf("8. Clear logs\n")
	fmt.Printf("---------------------\n")
	fmt.Printf("0. Exit\n")
	fmt.Print(colors.RESET)
}

func handleUserInput(userInput string) error {
	switch userInput {
	case "1":
		println("handle case handle safe mode")
	case "2":
		var TaskManager brute.TaskManager
		TaskManager.Begin()
	case "3":
		// acc generator start
		accgen.Start()
		fmt.Printf("%sPress enter to navigate back to the main menu\n%s", colors.MAGENTA, colors.RESET)
		GetUserInput()
	case "4":
		// open webhook settings
		WebhookMenu()
	case "5":
		// open reservation settings
		utils.OpenExplorer("./config/tasks.csv")
	case "6":
		// open generator settings
		utils.OpenExplorer("./config/accountgenconfig.csv")
	case "7":
		// open accounts list in config
		utils.OpenExplorer("./config/accounts.csv")
	case "8":
		utils.ClearFile("./config/logs.txt")
		// clear logs file without deleting it
	}

	return nil
}

func PrintWebhookHeader() {
	configF, err := utils.OpenConfigJson("./config/config.json")
	if err != nil {
		utils.Log("CLIENT", "Error locating config.json", "")
		return
	}
	var outString string
	isWebhookSet := len(configF.WebhookUrl) != 0
	if isWebhookSet {
		outString = fmt.Sprintf("%s%s%s",
			colors.GREEN, configF.WebhookUrl, colors.RESET)
	} else {
		outString = fmt.Sprintf("%sNot set%s",
			colors.RED, colors.RESET)
	}
	fmt.Printf("%sWebhook: %s\n%s---------------------\n%s",
		colors.YELLOW, outString, colors.CYAN, colors.RESET)
}

func PrintWebhookOptions() {
	fmt.Print(colors.CYAN)
	fmt.Printf("1. Test webhook\n")
	fmt.Printf("2. Edit webhook\n")
	fmt.Printf("---------------------\n")
	fmt.Printf("0. Main menu\n")
	fmt.Print(colors.RESET)
}

func handleUserWebhookInput(userInput string) {
	// utils.Log("CLIENT", fmt.Sprintf("Handling user input:%s", userInput), "")
	switch userInput {
	case "1":
		err := SendWebhookTest()
		if err != nil {
			fmt.Printf("Error sending webhook: %s", err.Error())
		}
	case "2":
		err := EditWebhook()
		if err != nil {
			fmt.Printf("Error editting webhook: %s", err.Error())
		}
	}

}
