package cmd

import (
	"log"
	"os/exec"
	"runtime"
	"strings"

	"github.com/codegangsta/cli"
	"github.com/exercism/cli/api"
	"github.com/exercism/cli/config"
)

// Open uses the given language and exercise and opens it in the browser
func Open(ctx *cli.Context) {
	c, err := config.New(ctx.GlobalString("config"))
	if err != nil {
		log.Fatal(err)
	}
	client := api.NewClient(c)

	args := ctx.Args()
	if len(args) != 2 {
		msg := "Usage: exercism open LANGUAGE EXERCISE"
		log.Fatal(msg)
	}

	language := args[0]
	exercise := args[1]
	submission, err := client.Submission(language, exercise)
	if err != nil {
		log.Fatal(err)
	}

	url := submission.URL
	// Escape characters not allowed by cmd/bash
	switch runtime.GOOS {
	case "windows":
		url = strings.Replace(url, "&", `^&`, -1)
	default:
		url = strings.Replace(url, "&", `\&`, -1)
	}

	// setup command to open browser
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", url)
	case "freebsd", "linux", "netbsd", "openbsd":
		cmd = exec.Command("xdg-open", url)
	case "windows":
		cmd = exec.Command("cmd", "/c", "start", url)
	}

	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}
}
