//go:build linux

package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"regexp"
	"strconv"
	"strings"
	"syscall"
	"time"

	"gitlab.com/gomidi/rtmididrv"

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
	if val, ok := e.Env["ID_MODEL_ENC"]; ok {
		result += fmt.Sprintf(", model=%s", val)
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

// getCardPortFromDevPath extracts card and port from DEVPATH string
func getCardPortFromDevPath(devPath string) (card, port int, ok bool) {
	// Example: /devices/platform/.../sound/card3/seq-midi-3-0
	re := regexp.MustCompile(`sound/card(\d+)/seq-midi-(\d+)-(\d+)`)
	matches := re.FindStringSubmatch(devPath)
	if len(matches) == 4 {
		card = atoi(matches[1])
		port = atoi(matches[3])
		return card, port, true
	}
	return 0, 0, false
}

func atoi(s string) int {
	n, _ := strconv.Atoi(s)
	return n
}

// listMidiDevices returns a slice of info about MIDI devices using gomidi/rtmididrv
func listMidiDevices() ([]string, error) {
	drv, err := rtmididrv.New()
	if err != nil {
		return nil, fmt.Errorf("could not init rtmididrv: %w", err)
	}
	defer drv.Close()
	ins, err := drv.Ins()
	if err != nil {
		return nil, fmt.Errorf("could not get MIDI inputs: %w", err)
	}
	var devices []string
	for _, in := range ins {
		devices = append(devices, fmt.Sprintf("id=%d name=%s", in.Number(), in.String(), drv.String()))
	}
	return devices, nil
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
		if (event.Subsystem == "sound" || event.Subsystem == "snd_seq") && (event.Action == "add" || event.Action == "remove") {
			fmt.Println(formatEvent(event))
			if event.Action == "add" && strings.Contains(event.DevPath, "seq-midi-") {
				card, port, ok := getCardPortFromDevPath(event.DevPath)
				fmt.Println("MIDI devices in system:")
				devices, err := listMidiDevices()
				if err != nil {
					fmt.Printf("Error listing MIDI devices: %v\n", err)
				} else {
					for _, dev := range devices {
						if ok && strings.Contains(dev, fmt.Sprintf("card=%d port=%d", card, port)) {
							fmt.Printf("* %s <-- just connected\n", dev)
						} else {
							fmt.Printf("  %s\n", dev)
						}
					}
				}
			}
		}
	}

	fmt.Println("Monitoring stopped.")
}
