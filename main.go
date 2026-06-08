package main

import (
	"flag"
	"fmt"
	"os"
)

var version = "0.1.0"

func main() {
	showVersion := flag.Bool("version", false, "Show version")
	resetFlag := flag.Bool("reset", false, "Factory reset")
	portFlag := flag.Int("port", 21475, "Port to listen on")
	noBrowser := flag.Bool("no-browser", false, "Don't open browser")
	flag.Parse()

	if *showVersion {
		fmt.Printf("u-hermes v%s\n", version)
		os.Exit(0)
	}

	if *resetFlag {
		fmt.Print("This will delete all config and data. Continue? (y/N): ")
		var answer string
		fmt.Scanln(&answer)
		if answer == "y" || answer == "Y" {
			os.Remove("config.json")
			os.RemoveAll("data")
			fmt.Println("Factory reset complete. Restart u-hermes to reconfigure.")
		}
		os.Exit(0)
	}

	fmt.Printf("u-hermes v%s starting on port %d...\n", version, *portFlag)
	_ = noBrowser
	// TODO: server start, tray, browser open (later tasks)
	select {}
}
