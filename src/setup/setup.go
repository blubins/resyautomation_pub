package setup

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"resysniper/src/colors"
	"resysniper/src/utils"
)

// Checks if a file at a given path exists
func CheckFile(path string) (bool, error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return false, err
	}
	_, err = os.Stat(absPath)
	if errors.Is(err, fs.ErrNotExist) {
		return false, fmt.Errorf("file does not exist")
	} else {
		return true, nil
	}
}

// Generates a file at a given path returning a pointer to the file
func GenerateFile(filePathName string) (*os.File, error) {
	file, err := os.Create(filePathName)
	if err != nil {
		return nil, err
	}
	return file, nil
}

// Generates a directory at a given path
func GenerateDirectory(path string) error {
	err := os.Mkdir(path, 0755)
	if err != nil {
		return err
	}
	return nil
}

// Initalizes the entire app structure of files
// writing in json fields and csv field data where necessary
func InitializeFiles() error {
	fmt.Printf(colors.RESET)
	utils.Log("CLIENT", "verifying app file structure", colors.YELLOW)

	var directoriesList []string = []string{
		"./config",
	}

	var configJsonWritables []string = []string{
		"./config/config.json",
	}

	var fileList []fileList = []fileList{
		{"./config/accountgenconfig.csv", "catch all,cc number,exp year,exp month,cvv,zip code,phone number,2captchakey,quantity"},
		{"./config/accounts.csv", "first name,last name,email,password,phone number,token,payment id\n"},
		{"./config/proxies.txt", ""},
		{"./config/tasks.csv", "email address,x resy auth token,day,time,room name,party size,payment id,restaurant id,delay,task schedule,run time (s)"},
	}

	utils.Log("CLIENT", "verifying app directories", colors.YELLOW)
	for directory := range directoriesList {
		var directoryExists bool
		utils.Log("CLIENT", "checking directory "+directoriesList[directory], colors.YELLOW)
		directoryExists, err := CheckFile(directoriesList[directory])
		if err != nil {
			utils.Log("ERROR", "error resolving files existence", colors.RED)
		}

		if !directoryExists {
			utils.Log("CLIENT", "directory: "+directoriesList[directory]+" does not exist, generating", colors.YELLOW)
			GenerateDirectory(directoriesList[directory])
		}
	}
	utils.Log("CLIENT", "verifying app files", colors.YELLOW)
	for file := range fileList {
		var fileExists bool
		utils.Log("CLIENT", "verifying file "+fileList[file].Path, colors.YELLOW)
		fileExists, err := CheckFile(fileList[file].Path)
		if err != nil && err.Error() != "file does not exist" {
			return err
		}

		if !fileExists {
			utils.Log("CLIENT", "file: "+fileList[file].Path+" does not exist, generating file", colors.YELLOW)
			f, err := GenerateFile(fileList[file].Path)
			if err != nil {
				return err
			}
			defer f.Close()

			_, err = f.Write([]byte(fileList[file].Content))
			if err != nil {
				return err
			}

		}
	}

	for i := range configJsonWritables {
		var fileExists bool
		utils.Log("CLIENT", "verifying file "+configJsonWritables[i], colors.YELLOW)
		fileExists, err := CheckFile(configJsonWritables[i])
		if err != nil && err.Error() != "file does not exist" {
			return err
		}

		if !fileExists {
			utils.Log("CLIENT", "file: "+configJsonWritables[i]+" does not exist, generating file", colors.YELLOW)
			f, err := GenerateFile(configJsonWritables[i])
			if err != nil {
				return err
			}
			defer f.Close()

			var writeConfig = userConfigJson{}
			dataToWrite, err := json.MarshalIndent(writeConfig, "", " ")
			if err != nil {
				return err
			}
			os.WriteFile(f.Name(), dataToWrite, 0644)
			f.Close()
		}
	}

	utils.Log("CLIENT", "app file structure verified", colors.GREEN)
	return nil
}
