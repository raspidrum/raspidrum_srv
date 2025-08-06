//go:build linux
// +build linux

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

	"github.com/raspidrum-srv/internal/repo/udev"
	libalsa "github.com/raspidrum-srv/libs/libalsa"
)

// formatEvent formats the event for output
func formatEvent(e *udev.Event, verbose bool) string {
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
	result += "\n"
	if verbose {
		for key, val := range e.Env {
			result += fmt.Sprintf("  %s=%s\n", key, val)
		}
	}
	return result
}

// printHelp prints help message
func printHelp() {
	fmt.Println("Udev device monitoring CLI")
	fmt.Println("Usage:")
	fmt.Println("  --cards      Show ALSA sound and MIDI cards, then exit")
	fmt.Println("  --monitor    Monitor device connect/disconnect events")
	fmt.Println("  -v           Enable verbose output. Only for --monitor")
	fmt.Println("  -h, --help   Show this help message")
	fmt.Println("\nNo arguments: show this help")
}

// getCardPortFromDevPath extracts card and port from DEVPATH string
func getCardPortFromDevPath(devPath string) (card, port int, ok bool) {
	// Example: /devices/platform/.../sound/card3/midiC3D0
	re := regexp.MustCompile(`sound/card(\d+)/midiC(\d+)D(\d+)`)
	matches := re.FindStringSubmatch(devPath)
	if len(matches) == 4 {
		card = atoi(matches[1])
		port = atoi(matches[3])
		return card, port, true
	}
	return 0, 0, false
}

func endsWithCardNumber(path string) bool {
	// Регулярное выражение для проверки /card и числа в конце
	re, err := regexp.Compile(`/card\d+$`)
	if err != nil {
		return false
	}
	return re.MatchString(path)
}

func atoi(s string) int {
	n, _ := strconv.Atoi(s)
	return n
}

// getAlsaCards возвращает список устройств ALSA и ошибки
func getAlsaCards() ([]string, error) {
	var result []string
	cards, err := libalsa.GetAllCards()
	if err != nil {
		return nil, err
	}
	for _, cardNum := range cards {
		cardInfo, err := libalsa.GetCardInfo(cardNum)
		if err != nil {
			result = append(result, fmt.Sprintf("Error getting info for card %d: %v", cardNum, err))
			continue
		}
		hwInfo, _ := libalsa.GetHardwareInfo(cardNum)
		info := fmt.Sprintf("Card %d: %s (%s), Driver: %s, Mixer: %s, Playback: %s, Capture: %s",
			cardInfo.ID, cardInfo.Name, cardInfo.LongName, cardInfo.Driver, cardInfo.Mixer,
			hwInfo["playback_devices"], hwInfo["capture_devices"])
		result = append(result, info)
	}

	// Print MIDI sequencer clients/ports
	midiPorts, err := libalsa.ListMidiPorts()
	if err == nil {
		for _, port := range midiPorts {
			var portType string
			switch port.PortType {
			case libalsa.MidiPortTypeHardware:
				portType = "hardware"
			case libalsa.MidiPortTypeSoftware:
				portType = "software"
			}
			result = append(result, fmt.Sprintf("MIDI: client=%d port=%d card=%d name=%s portname=%s\n type=%s input=%t output=%t",
				port.ClientId, port.PortId, port.CardId, port.ClientName, port.PortName, portType, port.IsInput, port.IsOutput))
		}
	}
	return result, nil
}

func listCards() {
	devices, err := getAlsaCards()
	if err != nil {
		fmt.Printf("Error listing ALSA devices: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("ALSA cards (sound and MIDI):")
	for _, dev := range devices {
		fmt.Println(dev)
	}
}

func runMonitor(verbose bool) {
	for {
		monitor, err := udev.NewMonitor()
		if err != nil {
			log.Fatalf("Failed to create monitor: %v", err)
		}
		ctx, cancel := context.WithCancel(context.Background())
		signalCh := make(chan os.Signal, 1)
		signal.Notify(signalCh, os.Interrupt, syscall.SIGTERM)
		terminatedBySignal := false
		go func() {
			<-signalCh
			terminatedBySignal = true
			fmt.Println("\nTermination signal received...")
			cancel()
		}()

		eventsCh, errCh, err := monitor.Start(ctx)
		if err != nil {
			log.Fatalf("Failed to start monitoring: %v", err)
		}
		fmt.Println("Starting udev device monitoring...")
		fmt.Println("Press Ctrl+C to stop")
		monitoringStopped := false
		for !monitoringStopped {
			select {
			case event, ok := <-eventsCh:
				if !ok {
					fmt.Println("Monitoring stopped. Restarting in 1s...")
					monitoringStopped = true
					break
				}
				fmt.Println(formatEvent(event, verbose))
				if (event.Subsystem == "sound" || event.Subsystem == "snd_seq") && (event.Action == "add" || event.Action == "change" || event.Action == "remove") {
					if event.Action == "change" && strings.Contains(event.DevPath, "card") {
						card, port, ok := getCardPortFromDevPath(event.DevPath)
						devices, err := getAlsaCards()
						if err != nil {
							fmt.Printf("Error listing ALSA devices: %v\n", err)
						} else {
							for _, dev := range devices {
								if ok && strings.Contains(dev, fmt.Sprintf("port=%d card=%d", port, card)) {
									fmt.Printf("* %s <-- just connected\n", dev)
								} else {
									fmt.Printf("  %s\n", dev)
								}
							}
						}
					}
					if event.Action == "delete" && endsWithCardNumber(event.DevPath) {
						devices, err := getAlsaCards()
						if err != nil {
							fmt.Printf("Error listing ALSA devices: %v\n", err)
						} else {
							for _, dev := range devices {
								fmt.Printf("  %s\n", dev)
							}
						}
					}
				}
			case err, ok := <-errCh:
				if ok && err != nil {
					fmt.Printf("Monitor error: %v\n", err)
				}
				fmt.Println("Monitoring stopped due to error. Restarting in 1s...")
				monitoringStopped = true
			}
		}
		cancel()
		if terminatedBySignal {
			fmt.Println("Exiting monitor loop due to termination signal.")
			break
		}
		time.Sleep(1 * time.Second)
		fmt.Println("Restarting udev monitor...")
	}
}

func main() {
	if len(os.Args) == 1 {
		printHelp()
		return
	}

	// Parse flags manually
	args := os.Args[1:]
	var monitor, verbose, cards, help bool
	for _, arg := range args {
		switch arg {
		case "--monitor":
			monitor = true
		case "-v":
			verbose = true
		case "--cards":
			cards = true
		case "-h", "--help":
			help = true
		}
	}

	if help {
		printHelp()
		return
	}
	if cards {
		listCards()
		return
	}
	if monitor {
		runMonitor(verbose)
		return
	}
	printHelp()
}
