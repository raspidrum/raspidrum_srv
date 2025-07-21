package usb

import (
	"context"
	"fmt"

	"github.com/google/gousb"
)

// Details contains detailed information about a USB device.
type Details struct {
	Class        uint8
	SubClass     uint8
	Protocol     uint8
	VendorID     gousb.ID
	ProductID    gousb.ID
	Manufacturer string
	Product      string
	SerialNumber string
	Bus          int
	Address      int
}

// DeviceManager provides access to USB devices.
type DeviceManager struct {
	ctx *gousb.Context
}

// NewDeviceManager creates a new DeviceManager.
func NewDeviceManager() *DeviceManager {
	return &DeviceManager{
		ctx: gousb.NewContext(),
	}
}

// Close closes the gousb context.
func (dm *DeviceManager) Close() {
	_ = dm.ctx.Close()
}

// GetDeviceDetails finds a device by vendor and product ID and returns its details.
func (dm *DeviceManager) GetDeviceDetails(ctx context.Context, vendorID, productID uint16) (*Details, error) {
	devs, err := dm.ctx.OpenDevices(func(desc *gousb.DeviceDesc) bool {
		return desc.Vendor == gousb.ID(vendorID) && desc.Product == gousb.ID(productID)
	})

	if err != nil {
		return nil, fmt.Errorf("failed to open devices: %w", err)
	}

	if len(devs) == 0 {
		return nil, nil // Device not found, but not an error
	}

	// Assume we are interested in the first device found.
	dev := devs[0]
	//defer dev.Close()

	manufacturer, _ := dev.Manufacturer()
	product, _ := dev.Product()
	serial, _ := dev.SerialNumber()

	details := &Details{
		Class:        uint8(dev.Desc.Class),
		SubClass:     uint8(dev.Desc.SubClass),
		Protocol:     uint8(dev.Desc.Protocol),
		VendorID:     dev.Desc.Vendor,
		ProductID:    dev.Desc.Product,
		Manufacturer: manufacturer,
		Product:      product,
		SerialNumber: serial,
		Bus:          dev.Desc.Bus,
		Address:      dev.Desc.Address,
	}

	return details, nil
}
