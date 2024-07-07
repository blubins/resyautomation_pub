package auth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"resysniper/src/cli"
	"resysniper/src/colors"
	"resysniper/src/utils"
	"strings"

	"github.com/denisbrodbeck/machineid"
)

func whopAuth(userLicense string) (bool, string, error) {
	machineId, err := machineid.ID()
	if err != nil {
		return false, "", err
	}
	requestDataStruct := map[string]interface{}{
		"metadata": map[string]string{
			"newKey": machineId,
		},
	}

	requestBody, err := json.Marshal(requestDataStruct)
	if err != nil {
		return false, "", err
	}

	url := fmt.Sprintf("https://api.whop.com/api/v2/memberships/%s/validate_license", userLicense)
	req, err := http.NewRequest(
		"POST",
		url,
		bytes.NewBuffer(requestBody),
	)
	if err != nil {
		return false, "", err
	}
	req.Header.Set("Authorization", "Bearer UnDKL3oroz-lAD8Hcav603kB8OfIu6jS6g8IPj1-rQ0")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		return false, "", err
	}
	defer response.Body.Close()
	switch response.StatusCode {
	case 201:
		return true, "Successfully verified your license", nil
	case 400:
		return false, "License already in use, please reset your key", nil
	case 404:
		return false, "License not found", nil
	default:
		return false, "Unexpected server reseponse please try again", nil
	}

}

// Takes users input for license and verifies with whop api
func takeInput() (bool, string, error) {
	fmt.Printf("%sPlease enter your license: %s", colors.MAGENTA, colors.RESET)
	userInput := cli.GetUserInput()
	trimmedInput := strings.TrimRight(userInput, "\r\n")
	verified, statusString, err := whopAuth(trimmedInput)
	if err != nil {
		return false, "", err
	}

	if verified {
		err = utils.UpdateConfigJson("./config/config.json", "License", trimmedInput)
		if err != nil {
			fmt.Print(err.Error())
		}
		return true, statusString, nil
	}
	fmt.Printf("%s\n", statusString)
	return false, "", err
}

// Will attempt to open user's existing config and read in and verify existing license
// if not found then will prompt user for key until validated indefinitely
func AuthorizeClient() (bool, string, error) {
	utils.Log("CLIENT", "authorizing license", colors.YELLOW)
	config, err := utils.OpenConfigJson("./config/config.json")
	if err != nil {
		return false, "error opening config file", err
	}

	if len(config.License) > 0 {
		verified, statusString, err := whopAuth(config.License)
		if err != nil {
			return false, "", err
		}

		if verified {
			utils.Log("CLIENT", "client authorized", colors.GREEN)
			return true, statusString, nil
		} else {
			takeInput()
			return false, statusString, nil
		}
	} else {
		takeInput()
		return false, "", nil
	}
}
