#!/bin/bash
# This script installs dependencies required for udev and usb support.
set -e

echo "Installing development libraries for udev and usb..."
sudo apt update
sudo apt install -y libudev-dev libusb-1.0-0-dev libportmidi-dev