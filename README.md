**Raspidrum backend service**

# Build, deploy and remote debug

## Prerequisites
- Docker (with Buildx)
- Colima (for Docker on macOS)
- Go 1.21+
- (For releases) goreleaser
- Raspberry Pi with SSH access, user: drum

## Build for Raspberry Pi (ARM64)

```bash
make build         # Release build (ARM64, no debug info)
make build-debug   # Debug build (ARM64, with debug info)
make clean         # Clean build/
```

## Deploy to Raspberry Pi

```bash
make deploy        # Deploy binary and configs
make deploy-full   # Deploy binary, configs, and DB
```

## Service management (systemd on RPi)

```bash
make install-service   # Install/update systemd unit
make update-service    # Update and restart the service
make start-service     # Start the service
make stop-service      # Stop the service
make restart-service   # Restart the service
make status-service    # Check service status
make enable-service    # Enable autostart
make disable-service   # Disable autostart
```

## Debugging on Raspberry Pi

1. Build and deploy the debug version, then start the remote debugger:
   ```bash
   make debug-remote
   # or manually:
   make build-debug && make deploy && make start-debug
   ```
2. In VSCode, use the launch config `Srv RPi debug` or `Srv RPi attach`.
   - Port: 2345
   - Path on RPi: `/opt/raspidrum/raspidrum`
   - Host: raspidrum-aabf.local (or IP)
3. Stop the debugger: `make stop-debug`

## Logging

- Logging level: environment variable `SRV_LOG_LEVEL` (`DEBUG`, `INFO`, `WARNING`, `ERROR`)
- Example:
  ```bash
  SRV_LOG_LEVEL=DEBUG ./build/raspidrum
  ```

## Release build (multiplatform)

```bash
make release-it
# Uses goreleaser (see .goreleaser.yaml)
```

## CI/CD and environment setup on macOS

See the detailed guide in `doc/ci_cd_for_rpi_macos.md`.


# Logging

Default logging level is INFO.

Set level by ENV variable SRV_LOG_LEVEL.
Can be:
- DEBUG
- INFO
- WARNING
- ERROR