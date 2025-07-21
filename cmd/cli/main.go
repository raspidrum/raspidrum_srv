//go:build linux

package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/raspidrum-srv/internal/repo/udev"
)

// formatEvent formats the event for output
func formatEvent(e *udev.Event) string {
	if e.Action == "" || e.Subsystem == "" {
		return ""
	}

	timeStr := time.Now().Format("2006-01-02 15:04:05")
	var prefix string

	switch e.Action {
	case "add":
		prefix = "✓ CONNECTED:"
	case "remove":
		prefix = "✗ DISCONNECTED:"
	case "change":
		prefix = "△ CHANGED:"
	default:
		prefix = strings.ToUpper(e.Action) + ":"
	}

	result := fmt.Sprintf("[%s] %s subsystem=%s", timeStr, prefix, e.Subsystem)
	if e.DevName != "" {
		result += fmt.Sprintf(", device=%s", e.DevName)
	}
	if e.DevType != "" {
		result += fmt.Sprintf(", type=%s", e.DevType)
	}
	if e.DevPath != "" {
		result += fmt.Sprintf(", path=%s", e.DevPath)
	}
	return result
}

// printHelp prints help message
func printHelp() {
	fmt.Println("Udev device monitoring")
	fmt.Println("Usage: ./udev_monitor [options]")
	fmt.Println("\nOptions:")
	fmt.Println("  -h, --help     Show this message")
	fmt.Println("  -v, --verbose  Verbose output (not implemented)")
	fmt.Println("\nPress Ctrl+C to stop")
}

func main() {
	if len(os.Args) > 1 && (os.Args[1] == "-h" || os.Args[1] == "--help") {
		printHelp()
		return
	}

	monitor, err := udev.NewMonitor()
	if err != nil {
		log.Fatalf("Failed to create monitor: %v", err)
	}
	// The monitor's lifecycle is managed by the context, no need for defer.

	ctx, cancel := context.WithCancel(context.Background())

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-signalCh
		fmt.Println("\nTermination signal received...")
		cancel()
	}()

	eventsCh, err := monitor.Start(ctx)
	if err != nil {
		log.Fatalf("Failed to start monitoring: %v", err)
	}

	fmt.Println("Starting udev device monitoring...")
	fmt.Println("Press Ctrl+C to stop")

	for event := range eventsCh {
		if event.Subsystem == "sound" && (event.Action == "add" || event.Action == "remove") {
			fmt.Println(formatEvent(event))
		}
	}

	fmt.Println("Monitoring stopped.")
}
