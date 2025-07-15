# CI/CD and Build Environment for Raspberry Pi on macOS

## 1. Install Docker and Colima

- Install Colima (recommended for ARM/M1):
  ```bash
  brew install colima
  colima start -r docker
  ```

- Install Docker CLI via Homebrew:
  ```bash
  brew install docker
  ```

## 2. Install Docker Buildx

- Check if buildx is installed:
  ```bash
  docker buildx version
  ```
- If buildx is not installed or you need a newer version:
  ```bash
  ARCH=arm64
  VERSION=v0.25.0
  curl -LO https://github.com/docker/buildx/releases/download/${VERSION}/buildx-${VERSION}.darwin-${ARCH}
  mkdir -p ~/.docker/cli-plugins
  mv buildx-${VERSION}.darwin-${ARCH} ~/.docker/cli-plugins/docker-buildx
  chmod +x ~/.docker/cli-plugins/docker-buildx
  docker buildx version # verify installation
  ```

## 3. Setup builder

- Create a builder for cross-compilation:
  ```bash
  docker buildx create --name mybuilder --driver docker-container --use mybuilder
  docker buildx inspect --bootstrap
  ```

## 4. Install goreleaser

- Via Homebrew:
  ```bash
  brew install --cask goreleaser/tap/goreleaser
  ```
- Or manually: https://goreleaser.com/install/
- Or by go:
  ```bash
  go install github.com/goreleaser/goreleaser/v2@latest
  ```

## 5. Verify environment

- Make sure everything is working:
  ```bash
  docker info
  docker buildx ls
  goreleaser --version
  colima status
  ```

## 6. Build and deploy

- Build for Raspberry Pi:
  ```bash
  make build         # Release build (ARM64)
  make build-debug   # Debug build (ARM64)
  ```
- Deploy to Raspberry Pi:
  ```bash
  make deploy        # Only binary and configs
  make deploy-full   # + DB
  ```

## 7. Debugging on Raspberry Pi

- Start the remote debugger:
  ```bash
  make debug-remote
  # or manually:
  make build-debug && make deploy && make start-debug
  ```
- Connect from VSCode: use the launch config `Srv RPi debug` (port 2345, path /opt/raspidrum/raspidrum).
- Stop the debugger: `make stop-debug`

## 8. Buildkit quirks and troubleshooting on macOS

- Sometimes buildkit/buildx requires manual installation (see step 2).
- If you encounter builder errors, recreate the builder:
  ```bash
  docker buildx rm mybuilder
  docker buildx create --name mybuilder --driver docker-container --use mybuilder
  docker buildx inspect --bootstrap
  ```
- For best stability, use Colima instead of standard Docker Desktop on M1/M2.

## 9. Useful links
- https://github.com/docker/buildx
- https://github.com/abiosoft/colima
- https://goreleaser.com/ 