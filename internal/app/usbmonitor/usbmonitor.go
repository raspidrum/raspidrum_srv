package usbmonitor

import (
	"context"
	"fmt"
	"log/slog"
	"strconv"
	"strings"

	"github.com/raspidrum-srv/internal/repo/udev"
	"github.com/raspidrum-srv/internal/repo/usb"
)

const (
	ClassAudio       = 0x01
	ClassMassStorage = 0x08
)

// MonitorService orchestrates USB device monitoring.
type MonitorService struct {
	udevMonitor *udev.Monitor
	usbManager  *usb.DeviceManager
}

// NewMonitorService creates a new MonitorService.
func NewMonitorService() (*MonitorService, error) {
	udevMon, err := udev.NewMonitor()
	if err != nil {
		return nil, err
	}
	usbMgr := usb.NewDeviceManager()
	return &MonitorService{
		udevMonitor: udevMon,
		usbManager:  usbMgr,
	}, nil
}

// Start begins the monitoring process.
func (s *MonitorService) Start(ctx context.Context) {
	events, err := s.udevMonitor.Start(ctx)
	if err != nil {
		slog.Info("USB monitoring is not available on this system", "error", err)
		return
	}

	go func() {
		for {
			select {
			case <-ctx.Done():
				slog.Info("Stopping USB monitor service.")
				s.udevMonitor.Close()
				s.usbManager.Close()
				return
			case event, ok := <-events:
				if !ok {
					slog.Info("Udev event channel closed.")
					return
				}
				if event.Subsystem == "usb" && event.Action == "add" {
					s.handleDeviceEvent(ctx, event)
				}
			}
		}
	}()
}

func (s *MonitorService) handleDeviceEvent(ctx context.Context, event *udev.Event) {
	vendorID, productID, err := getDeviceIDsFromEvent(event)
	if err != nil {
		slog.Warn("Failed to get device IDs from udev event", "error", err)
		return
	}

	details, err := s.usbManager.GetDeviceDetails(ctx, vendorID, productID)
	if err != nil {
		slog.Warn("Failed to get device details", "error", err, "vendor", vendorID, "product", productID)
		return
	}
	if details == nil {
		slog.Debug("Device details not found, skipping", "vendor", vendorID, "product", productID)
		return
	}

	if details.Class == ClassAudio || details.Class == ClassMassStorage {
		slog.Info("Detected USB device",
			slog.Uint64("class", uint64(details.Class)),
			slog.Uint64("subclass", uint64(details.SubClass)),
			slog.Uint64("vendor_id", uint64(details.VendorID)),
			slog.Uint64("product_id", uint64(details.ProductID)),
			slog.String("manufacturer", details.Manufacturer),
			slog.String("product", details.Product),
			slog.String("serial_number", details.SerialNumber),
			slog.Int("bus", details.Bus),
			slog.Int("address", details.Address),
		)
	}
}

func getDeviceIDsFromEvent(event *udev.Event) (uint16, uint16, error) {
	productStr, ok := event.Env["PRODUCT"]
	if !ok {
		return 0, 0, fmt.Errorf("PRODUCT not found in uevent env")
	}

	parts := strings.Split(productStr, "/")
	if len(parts) < 2 {
		return 0, 0, fmt.Errorf("invalid PRODUCT format: %s", productStr)
	}

	vendorID, err := strconv.ParseUint(parts[0], 16, 16)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to parse vendor id: %w", err)
	}

	productID, err := strconv.ParseUint(parts[1], 16, 16)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to parse product id: %w", err)
	}

	return uint16(vendorID), uint16(productID), nil
}
