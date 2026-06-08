package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"u-hermes/internal/chat"
	"u-hermes/internal/config"
	"u-hermes/internal/server"
	"u-hermes/internal/store"
	"u-hermes/internal/tray"

	"github.com/skratchdot/open-golang/open"
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

	// Single instance check
	if !acquireLock(*portFlag) {
		// Another instance is running, just open browser
		open.Run(fmt.Sprintf("http://127.0.0.1:%d/chat", *portFlag))
		os.Exit(0)
	}

	// Find port
	port := findPort(*portFlag)

	// Determine paths
	exeDir := filepath.Dir(mustExePath())
	configPath := filepath.Join(exeDir, "config.json")
	dataDir := filepath.Join(exeDir, "data")

	// Load or create config
	cfg, err := config.Load(configPath)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	isFirstRun := len(cfg.Models) == 0

	// Open store
	st, err := store.Open(dataDir)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer st.Close()

	// Chat service
	chatSvc := chat.NewService()

	// Create server
	srv := server.New(st, cfg, chatSvc, configPath, webAssets, port)

	// Start HTTP server
	go func() {
		log.Printf("U-Hermes v%s listening on http://127.0.0.1:%d", version, port)
		if err := srv.Engine().Run(fmt.Sprintf("127.0.0.1:%d", port)); err != nil {
			log.Fatalf("Server error: %v", err)
		}
	}()

	// Open browser
	if !*noBrowser {
		targetURL := fmt.Sprintf("http://127.0.0.1:%d/chat", port)
		if isFirstRun {
			targetURL = fmt.Sprintf("http://127.0.0.1:%d/onboarding", port)
		}
		open.Run(targetURL)
	}

	// System tray
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	t := tray.New(port,
		func() { open.Run(fmt.Sprintf("http://127.0.0.1:%d/chat", port)) },
		func() { open.Run(fmt.Sprintf("http://127.0.0.1:%d/settings", port)) },
		func() {
			st.Close()
			os.Exit(0)
		},
	)

	go func() {
		<-sigCh
		t.Quit()
	}()

	t.Run()
}

func acquireLock(port int) bool {
	conn, err := net.Dial("tcp", fmt.Sprintf("127.0.0.1:%d", port))
	if err == nil {
		conn.Close()
		return false
	}
	return true
}

func findPort(preferred int) int {
	ln, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", preferred))
	if err == nil {
		ln.Close()
		return preferred
	}
	ln, err = net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		log.Fatalf("No available port: %v", err)
	}
	defer ln.Close()
	return ln.Addr().(*net.TCPAddr).Port
}

func mustExePath() string {
	p, err := os.Executable()
	if err != nil {
		log.Fatalf("Cannot determine executable path: %v", err)
	}
	return p
}
