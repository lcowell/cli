package cmd

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime"
	"time"

	"github.com/codegangsta/cli"
	"github.com/exercism/cli/config"
	"github.com/exercism/cli/paths"
)

// Debug provides information about the user's environment and configuration.
func Debug(ctx *cli.Context) {
	defer fmt.Printf("\nIf you are having trouble and need to file a GitHub issue (https://github.com/exercism/exercism.io/issues) please include this information (except your API key. Keep that private).\n")

	client := http.Client{Timeout: 5 * time.Second}

	fmt.Printf("\n**** Debug Information ****\n")
	fmt.Printf("Exercism CLI Version: %s\n", ctx.App.Version)

	rel, err := fetchLatestRelease(client)
	if err != nil {
		log.Println("unable to fetch latest release: " + err.Error())
	} else {
		if rel.Version() != ctx.App.Version {
			defer fmt.Printf("\nA newer version of the CLI (%s) can be downloaded here: %s\n", rel.TagName, rel.Location)
		}
		fmt.Printf("Exercism CLI Latest Release: %s\n", rel.Version())
	}

	fmt.Printf("OS/Architecture: %s/%s\n", runtime.GOOS, runtime.GOARCH)
	fmt.Printf("Build OS/Architecture %s/%s\n", BuildOS, BuildARCH)
	if BuildARM != "" {
		fmt.Printf("Build ARMv%s\n", BuildARM)
	}

	fmt.Printf("Home Dir: %s\n", paths.Home)

	c, err := config.New(ctx.GlobalString("config"))
	if err != nil {
		log.Fatal(err)
	}

	configured := true
	if _, err = os.Stat(c.File); err != nil {
		if os.IsNotExist(err) {
			configured = false
		} else {
			log.Fatal(err)
		}
	}

	if configured {
		fmt.Printf("Config file: %s\n", c.File)
		fmt.Printf("API Key: %s\n", c.APIKey)
	} else {
		fmt.Println("Config file: <not configured>")
		fmt.Println("API Key: <not configured>")
	}

	fmt.Printf("API: %s [%s]\n", c.API, pingURL(client, c.API))
	fmt.Printf("XAPI: %s [%s]\n", c.XAPI, pingURL(client, c.XAPI))
	fmt.Printf("Exercises Directory: %s\n", c.Dir)
}

func pingURL(client http.Client, url string) string {
	res, err := client.Get(url)
	if err != nil {
		return err.Error()
	}
	defer res.Body.Close()

	return "connected"
}
