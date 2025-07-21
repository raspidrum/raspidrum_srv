# Building on macOS (Darwin)

To build the `raspidrum_srv` application on macOS, you need to have Go and some system dependencies installed.

## Prerequisites

- [Go](https://golang.org/doc/install) (version 1.21 or higher)
- [Homebrew](https://brew.sh/)

## Installation

1.  **Install system dependencies:**

    The project uses `gousb` which requires `libusb` and `pkg-config`. You can install them using Homebrew:

    ```bash
    brew install libusb pkg-config
    ```

2.  **Patch libusb:**

    Due to an incompatibility between `libusb-1.0` and LLVM on macOS, you need to apply a patch. A script is provided for this purpose. See doc in gousb.

    ```bash
    # Find the path to libusb.h
    LIBUSB_PATH=$(brew --prefix libusb)/include/libusb-1.0/libusb.h
    ```

3.  **Build the application:**

    Navigate to the `raspidrum_srv` directory and run the build command:

    ```bash
    go build ./...
    ```

    This will compile the application. Note that USB monitoring features relying on `udev` are only available on Linux and will be disabled on macOS. The application will build and run, but will log a message indicating that USB monitoring is not available. 