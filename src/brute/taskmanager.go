package brute

import (
	"fmt"
	"resysniper/src/colors"
	"resysniper/src/utils"
	"strconv"
	"strings"

	"github.com/google/uuid"
)

type TaskManager struct {
	Tasks []*ResyTask
}

// brute task menu and manager main function
// prompts user with options to start all or stop all given tasks and exit
func (tm *TaskManager) Begin() error {
	utils.ClearConsole()
	fmt.Printf("%s%s%s", colors.RED, utils.ResyTitle, colors.RESET)
	err := tm.InitTasks()
	if err != nil {
		utils.Log("CLIENT", fmt.Sprintf("Error initializing tasks: %s", err.Error()), colors.RED)
	}
	var userInput string
	utils.Log("CLIENT", fmt.Sprintf("%d tasks initialized", len(tm.Tasks)), colors.YELLOW)
	for userInput != "0" {
		fmt.Printf("%s-------------\n1. Start all tasks\n2. Stop all tasks\n3. Help\n-------------\n0. Exit, stop all tasks and return to main menu\n%s", colors.CYAN, colors.RESET)
		fmt.Scanln(&userInput)
		if userInput == "2" {
			tm.StopAllTasks()
		}
		if userInput == "1" {
			tm.StartAllTasks()
		}
		if userInput == "3" {
			fmt.Printf(colors.GREEN)
			utils.PrintBruteHelp()
			fmt.Printf(colors.RESET)
		}
	}
	tm.StopAllTasks()
	return nil
}

// initializes all tasks from config into task managers internal slice
func (tm *TaskManager) InitTasks() error {
	utils.Log("CLIENT", "Initializing tasks", colors.YELLOW)
	// make sure to reinitalize slice back to 0 incase of navigating
	// between windows and appending a ton of the same tasks
	tm.Tasks = tm.Tasks[:0]
	usersTaskConfigFile, err := utils.OpenConfigFile("./config/tasks.csv")
	if err != nil {
		return err
	}

	splitConfig := strings.Split(usersTaskConfigFile, "\n")
	if len(splitConfig) < 1 {
		return fmt.Errorf("tasks.csv not properly configured")
	}

	userConfigJson, err := utils.OpenConfigJson("./config/config.json")
	if err != nil {
		return err
	}
	// csv format
	// [0] = email address
	// [1] = x resy auth token
	// [2] = day
	// [3] = time
	// [4] = room name
	// [5] = party size
	// [6] = payment id
	// [7] = restaurant id
	// [8] = delay
	// [9] = task schedule
	// [10] = run time (s)
	for i := 1; i < len(splitConfig); i++ {
		splitTask := strings.Split(splitConfig[i], ",")
		configLength := len(splitTask)
		// in case of empty rows in csv return nil to skip over empty rows
		if configLength == 1 {
			return nil
		}
		if configLength != 11 {
			return fmt.Errorf("invalid brute task configuration, missing or invalid field length")
		}
		for i := 0; i < len(splitTask); i++ {
			splitTask[i] = strings.ReplaceAll(splitTask[i], "\r", "")
		}

		var rt ResyTask
		rt.TaskIdentifier = uuid.NewString()
		// required for neither
		rt.WebhookUrl = userConfigJson.WebhookUrl
		// required for brute
		rt.EmailAddress = splitTask[0]
		// required for both
		rt.XResyAuthToken = splitTask[1]
		// required for both
		rt.Day = splitTask[2]
		// required for safe
		rt.Time = splitTask[3]
		// required for safe
		rt.RoomName = splitTask[4]
		// required for both
		partySize, _ := strconv.Atoi(splitTask[5])
		rt.PartySize = partySize

		// required for both
		rt.PaymentId = splitTask[6]

		// required for both
		rt.RestaurantId = splitTask[7]

		// required for both
		delayInt, _ := strconv.Atoi(splitTask[8])
		rt.Delay = delayInt

		// required for safe
		rt.TaskSchedule = splitTask[9]

		// required for safe
		runtimeInt, _ := strconv.Atoi(splitTask[10])

		rt.RunTime = runtimeInt
		rt.StopChannel = make(chan struct{})
		rt.Running = false
		err = rt.ValidateTaskData()
		if err != nil {
			return err
		}
		tm.Tasks = append(tm.Tasks, &rt)
	}

	return nil
}

// Starts all tasks from TaskManager's internal slice
func (tm *TaskManager) StartAllTasks() error {
	utils.Log("CLIENT", fmt.Sprintf("Starting %d tasks", len(tm.Tasks)), colors.CYAN)
	for _, task := range tm.Tasks {
		if !task.Running {
			task.Running = true
			task.StopChannel = make(chan struct{})
			go task.Start()
		}
	}
	return nil
}

// Stops all tasks from TaskManager's internal slice
func (tm *TaskManager) StopAllTasks() error {
	utils.Log("CLIENT", fmt.Sprintf("Stopping %d tasks", len(tm.Tasks)), colors.YELLOW)
	for _, task := range tm.Tasks {
		if task.Running {
			task.Kill()
		}
	}
	return nil
}
