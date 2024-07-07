package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"resysniper/src/colors"
	"runtime"
	"strings"
	"time"

	api2captcha "github.com/2captcha/2captcha-go"
)

// Solves captcha given a 2captcha api token
// will return solution token as a string
func SolveCaptcha(captchakey string) (string, error) {
	client := api2captcha.NewClient(captchakey)
	cap := api2captcha.ReCaptcha{
		SiteKey:   "6Lfw-dIZAAAAAESRBH4JwdgfTXj5LlS1ewlvvCYe",
		Url:       "https://resy.com/",
		Invisible: false,
		Action:    "verify",
	}
	req := cap.ToRequest()
	code, err := client.Solve(req)
	if err != nil {
		return "", err
	}
	return code, nil
}

// Opens the file explorer of a given path
func OpenExplorer(path string) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		fmt.Printf("Could not find reservation settings path")
	}
	cmd := exec.Command("explorer", "/open,", absPath)
	cmd.Run()
}

// Will truncate a file to of length 0
func ClearFile(path string) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		fmt.Printf("Could not find logs file path\n")
	}
	err = os.Truncate(absPath, 0)
	if err != nil {
		fmt.Printf("Error truncating file %s\n", path)
	}
}

// Clears the terminal window
func ClearConsole() {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("cmd", "/c", "cls")
	case "darwin", "linux":
		cmd = exec.Command("clear")
	default:
		fmt.Printf("Unsupported operating system\n")
		return
	}
	cmd.Stdout = os.Stdout
	cmd.Run()
}

// Opens and returns a config file as a string
func OpenConfigFile(path string) (string, error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return "", fmt.Errorf("failed to get absolute path: %v", err)
	}
	f, err := os.Open(absPath)
	if err != nil {
		return "", fmt.Errorf("failed to open file: %s", absPath)
	}
	defer f.Close()
	var data string
	buf := make([]byte, 1024)
	for {
		n, err := f.Read(buf)
		if err == io.EOF {
			break
		}
		if err != nil {
			return "", fmt.Errorf("error reading file %s: %v", absPath, err)
		}
		if n > 0 {
			data += string(buf[:n])
		}
	}
	if err := f.Close(); err != nil {
		return "", fmt.Errorf("failed to close file %s: %v", absPath, err)
	}
	return data, nil
}

// Opens the user's config.json and returns UserConfigJson struct
func OpenConfigJson(path string) (UserConfigJson, error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return UserConfigJson{}, err
	}

	f, err := os.Open(absPath)
	if err != nil {
		return UserConfigJson{}, err
	}

	defer f.Close()

	byteValue, _ := io.ReadAll(f)
	var internalConfig UserConfigJson
	json.Unmarshal(byteValue, &internalConfig)
	return internalConfig, nil
}

// Will update a JSON value of a given file
// given it's path, fieldName and newValue
func UpdateConfigJson(path, fieldName, newValue string) error {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return err
	}

	fileContent, err := os.ReadFile(absPath)
	if err != nil {
		return err
	}

	var existingConfig UserConfigJson
	err = json.Unmarshal(fileContent, &existingConfig)
	if err != nil {
		return err
	}

	reflectedValue := reflect.ValueOf(&existingConfig).Elem()
	field := reflectedValue.FieldByName(fieldName)
	if !field.IsValid() {
		return fmt.Errorf("field name %s does not exist in config", fieldName)
	}
	field.SetString(newValue)

	updatedData, err := json.MarshalIndent(existingConfig, "", " ")
	if err != nil {
		return err
	}

	err = os.WriteFile(absPath, updatedData, 0644)
	if err != nil {
		return err
	}

	return nil
}

// appends text to a given file as name implies
func AppendFileSync(filename, content string) error {
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(content)
	if err != nil {
		return err
	}

	return nil
}

// Log saves the given data to the local ./config/logs.txt file and prints the message to the console.
// The prefix parameter will be formatted within brackets [%s] after the current date and time.
// The message parameter is the text to be saved and printed to the console.
// The color parameter specifies the color of the text printed to the console.
// The function returns an error if it fails to write to the file.
func Log(prefix, message, color string) {
	currentTime := time.Now()
	currentDate := currentTime.Format("2006-01-02 15:04:05")

	logMessage := fmt.Sprintf(`%s[%s] [%s] - %s%s`,
		color, currentDate, prefix, message, colors.RESET)

	logMessageForFile := fmt.Sprintf("[%s] [%s] - %s",
		currentDate, prefix, message)

	err := AppendFileSync("./config/logs.txt", logMessageForFile+"\n")
	if err != nil {
		fmt.Printf("error appending logs: %v", err.Error())
		return
	}

	fmt.Printf("%s\n", logMessage)
}

// TODO find less jank solution
// url.QueryEscape is terrible does not work
// using this instead as placeholder
func UrlEncode(s string) string {
	outString := strings.ReplaceAll(s, "{", "%7B")
	outString = strings.ReplaceAll(outString, "}", "%7D")
	outString = strings.ReplaceAll(outString, ":", "%3a")
	outString = strings.ReplaceAll(outString, `"`, "%22")
	outString = strings.ReplaceAll(outString, "=", "%3D")
	return outString
}

// Prints requirement information about brute mode
func PrintBruteHelp() {
	fmt.Printf("\n%sBrute mode will constantly monitor the sites reservation slots.\n", colors.GREEN)
	fmt.Printf("Once a slot is found it will randomly chose one and immediately\n")
	fmt.Printf("book the reservation.\n\n")
	fmt.Printf("The required fields for brute within tasks.csv are:\n")
	fmt.Printf("------Options--------------Examples----\n")
	fmt.Printf("1. Email Address     | blubs@gmail.com\n")
	fmt.Printf("2. X Resy Auth Token | eyJ0eXAiOiJK...\n")
	fmt.Printf("3. Day               | 2024-06-12\n")
	fmt.Printf("4. Party Size        | 2, 4, etc.\n")
	fmt.Printf("5. Restaurant ID     | 54139\n")
	fmt.Printf("6. Delay (ms)        | 5000\n")
	fmt.Printf("---------------------------------------\n")
	fmt.Printf("You do not need:\n")
	fmt.Printf("1. Time\n")
	fmt.Printf("2. Room Name\n")
	fmt.Printf("3. Task Schedule\n")
	fmt.Printf("4. Run Time (s)\n%s", colors.RESET)
}

const ResyTitle string = `
 /$$$$$$$                                    /$$$$$$          /$$                           
| $$__  $$                                  /$$__  $$        |__/                           
| $$  \ $$ /$$$$$$  /$$$$$$$/$$   /$$      | $$  \__//$$$$$$$ /$$ /$$$$$$  /$$$$$$  /$$$$$$ 
| $$$$$$$//$$__  $$/$$_____| $$  | $$      |  $$$$$$| $$__  $| $$/$$__  $$/$$__  $$/$$__  $$
| $$__  $| $$$$$$$|  $$$$$$| $$  | $$       \____  $| $$  \ $| $| $$  \ $| $$$$$$$| $$  \__/
| $$  \ $| $$_____/\____  $| $$  | $$       /$$  \ $| $$  | $| $| $$  | $| $$_____| $$      
| $$  | $|  $$$$$$$/$$$$$$$|  $$$$$$$      |  $$$$$$| $$  | $| $| $$$$$$$|  $$$$$$| $$      
|__/  |__/\_______|_______/ \____  $$       \______/|__/  |__|__| $$____/ \_______|__/      
                            /$$  | $$                           | $$                        
                           |  $$$$$$/                           | $$                        
                            \______/                            |__/                        
`
