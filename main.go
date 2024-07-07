package main

import (
	"os"
	"os/exec"
	"resysniper/src/auth"
	"resysniper/src/cli"
	"resysniper/src/setup"
	"time"

	"golang.org/x/sys/windows"
)

func initConsoleWindow() error {
	stdout := windows.Handle(os.Stdout.Fd())
	var originalMode uint32

	windows.GetConsoleMode(stdout, &originalMode)
	windows.SetConsoleMode(stdout, originalMode|windows.ENABLE_VIRTUAL_TERMINAL_PROCESSING)
	cmd := exec.Command("cmd", "/C", "title", "Resy Sniper 1.0.5")
	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}

func main() {
	initConsoleWindow()

	err := setup.InitializeFiles()
	if err != nil {
		println(err.Error())
	}

	for {
		isAuthed, _, _ := auth.AuthorizeClient()
		if isAuthed {
			break
		}
		// make sure to not spam whop api before attempting again
		time.Sleep(time.Second * 1)
	}

	cli.MainMenu()
}
