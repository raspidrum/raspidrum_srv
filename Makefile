include .env_rpi
APP_NAME=raspidrum
RD_USER=drum
#RD_HOST=raspidrum-aabf.local
RD_HOST=$(RHOST)
APP_PATH=/opt/raspidrum
SRC_SERVICE_DIR=linux
SERVICE_FILE=raspidrum.service
REMOTE_SYSTEMD_DIR=/etc/systemd/system
SRC_DIR=cmd/server
SRC_DB_DIR=db
SRC_CFG_DIR=configs
RD_CONFIG=$(RDRUM_CONFIG)

BOLD=\033[1m
REGULAR=\033[0m

.PHONY: build build-debug clean prepare-builder

prepare-builder:
	docker buildx build \
		--cache-from type=local,src=/tmp/buildkit-cache \
  	--cache-to type=local,dest=/tmp/buildkit-cache,mode=max \
		--file build.Dockerfile \
		--output type=docker \
		-t raspidrum-builder .

# Build for Raspberry Pi
build: prepare-builder
	@echo "Building for Raspberry Pi ARM64... (without debug info)"
	@mkdir -p build
	docker run --rm \
		--platform linux/arm64 \
		--mount type=bind,src=.,dst=/src \
	  -v /tmp/buildkit-cache:/root/.cache/go-build \
	  -v /tmp/buildkit-cache:/go/pkg/mod \
	  --name raspidrum-builder \
	  raspidrum-builder \
	 	go build -ldflags="-s -w" -buildvcs=false -trimpath -o ./build/$(APP_NAME) ./$(SRC_DIR)

# Build debug version
build-debug: prepare-builder
	@echo "Building debug version for Raspberry Pi ARM64..."
	@mkdir -p build
	docker run --rm \
		--platform linux/arm64 \
		--mount type=bind,src=.,dst=/src \
	  -v /tmp/buildkit-cache:/root/.cache/go-build \
	  -v /tmp/buildkit-cache:/go/pkg/mod \
	  --name raspidrum-builder \
	  raspidrum-builder \
	  go build -gcflags="all=-N -l" -buildvcs=false -trimpath -o ./build/$(APP_NAME) ./$(SRC_DIR)

# Clean build files
clean:
	@echo "Cleaning build directory..."
	rm -rf build


.PHONY: deploy deploy-full
deploy:
	@echo "Deploying to Raspberry Pi..."
	ssh $(RD_USER)@$(RD_HOST) "mkdir -p $(APP_PATH)"
	scp -r $(SRC_CFG_DIR) build/$(APP_NAME) $(RD_USER)@$(RD_HOST):$(APP_PATH)/
	ssh $(RD_USER)@$(RD_HOST) "chmod +x $(APP_PATH)/$(APP_NAME)"

deploy-full:
	@echo "Deploying to Raspberry Pi..."
	ssh $(RD_USER)@$(RD_HOST) "mkdir -p $(APP_PATH)"
	scp -r $(SRC_DB_DIR) $(SRC_CFG_DIR) build/$(APP_NAME) $(RD_USER)@$(RD_HOST):$(APP_PATH)/
	ssh $(RD_USER)@$(RD_HOST) "chmod +x $(APP_PATH)/$(APP_NAME)"

.PHONY: start-debug debug-remote stop-debug
# Start remote debugging
debug-remote: build-debug deploy start-debug

start-debug: stop-debug
	@echo "Starting remote debugger..."
	ssh $(RD_USER)@$(RD_HOST) "cd $(APP_PATH) && tmux new-session -d 'env RDRUM_CONFIG=$(RD_CONFIG) ~/go/bin/dlv dap --listen=:2345 > /tmp/dlv.log 2>&1'"
	@echo "Debugger started on $(RD_HOST):2345"
	@echo "Connect via VS Code or run: dlv connect $(RD_HOST):2345"

# Stop debugger
stop-debug:
	@echo "Stopping remote debugger..."
	ssh $(RD_USER)@$(RD_HOST) "pkill -f dlv" || true




.PHONY: start-service stop-service restart-service status-service enable-service disable-service update-service
# Install service
install-service:
	@echo "Deploying service $(APP_NAME) on $(RD_USER)@$(RD_HOST)..."
	
	# Copy service file to remote machine
	scp $(SRC_SERVICE_DIR)/$(SERVICE_FILE) $(RD_USER)@$(RD_HOST):/tmp/
	
	# Install and activate service
	ssh $(RD_USER)@$(RD_HOST) '\
		sudo mv /tmp/$(SERVICE_FILE) $(REMOTE_SYSTEMD_DIR)/ && \
		sudo chown root:root $(REMOTE_SYSTEMD_DIR)/$(SERVICE_FILE) && \
		sudo chmod 644 $(REMOTE_SYSTEMD_DIR)/$(SERVICE_FILE) && \
		sudo systemctl daemon-reload && \
		sudo systemctl enable $(APP_NAME) && \
		sudo systemctl restart $(APP_NAME) && \
		sudo systemctl status $(APP_NAME)'
	
	@echo "Service $(APP_NAME) successfully installed and started!"


# Start service
start-service:
	@echo "Starting $(APP_NAME) service..."
	ssh $(RD_USER)@$(RD_HOST) "sudo systemctl start $(APP_NAME)"

# Stop service
stop-service:
	@echo "Stopping $(APP_NAME) service..."
	ssh $(RD_USER)@$(RD_HOST) "sudo systemctl stop $(APP_NAME)"

# Restart service
restart-service:
	@echo "Restarting $(APP_NAME) service..."
	ssh $(RD_USER)@$(RD_HOST) "sudo systemctl restart $(APP_NAME)"

# Check service status
status-service:
	@echo "Checking $(APP_NAME) service status..."
	ssh $(RD_USER)@$(RD_HOST) "sudo systemctl status $(APP_NAME)"

# Enable service (auto-start on boot)
enable-service:
	@echo "Enabling $(APP_NAME) service for auto-start..."
	ssh $(RD_USER)@$(RD_HOST) "sudo systemctl enable $(APP_NAME)"

# Disable service (no auto-start on boot)
disable-service:
	@echo "Disabling $(APP_NAME) service auto-start..."
	ssh $(RD_USER)@$(RD_HOST) "sudo systemctl disable $(APP_NAME)"

# Update and restart service (for quick updates)
update-service: stop-service deploy start-service
	@echo "Service updated and restarted"


.PHONY: test-connection service-logs service-logs-tail logs
# Test connection to Raspberry Pi
test-connection:
	@echo "Testing connection to Raspberry Pi..."
	ssh $(RD_USER)@$(RD_HOST) "echo 'Connection successful'"

service-logs:
	@echo "Showing $(APP_NAME) service logs..."
	ssh $(RD_USER)@$(RD_HOST) "sudo journalctl -f -u $(APP_NAME)"

service-logs-tail:
	@echo "Showing last 50 lines of $(APP_NAME) service logs..."
	ssh $(RD_USER)@$(RD_HOST) "sudo journalctl -n 50 -u $(APP_NAME)"

logs:
	@echo "Showing $(APP_NAME) debug logs..."
	ssh $(RD_USER)@$(RD_HOST) "tail -f /tmp/dlv.log"

# Build release
PACKAGE_NAME          := github.com/raspidrum-srv
GOLANG_CROSS_VERSION  ?= v1.21.5

.PHONY: release-it
release-it:
	@docker run \
		--rm \
		-e CGO_ENABLED=1 \
		-e GOPROXY="direct" \
		-v /var/run/docker.sock:/var/run/docker.sock \
		-v `pwd`:/go/src/$(PACKAGE_NAME) \
		-w /go/src/$(PACKAGE_NAME) \
		ghcr.io/goreleaser/goreleaser-cross:${GOLANG_CROSS_VERSION} \
		--clean --skip=validate --skip=publish --snapshot

.PHONY: help
# Show help
help:
	@echo "Usage: make <target>"
	@echo "Targets:"
	@echo "# Build"
	@echo "  ${BOLD}build${REGULAR} - Build for Raspberry Pi ARM64"
	@echo "  ${BOLD}build-debug${REGULAR} - Build debug version for Raspberry Pi ARM64"
	@echo "  ${BOLD}clean${REGULAR} - Clean build files"
	@echo "  ${BOLD}release-it${REGULAR} - Build multiplatform release and pack"
	@echo "\n# Deploy"
	@echo "  ${BOLD}update-service${REGULAR} - Update and restart $(APP_NAME) service"
	@echo "  ${BOLD}deploy${REGULAR} - Deploy to Raspberry Pi"
	@echo "  ${BOLD}deploy-full${REGULAR} - Deploy with DB to Raspberry Pi"
	@echo "\n# Service management"
	@echo "  ${BOLD}start-service${REGULAR} - Start $(APP_NAME) service"
	@echo "  ${BOLD}stop-service${REGULAR} - Stop $(APP_NAME) service"
	@echo "  ${BOLD}restart-service${REGULAR} - Restart $(APP_NAME) service"
	@echo "  ${BOLD}status-service${REGULAR} - Check $(APP_NAME) service status"
	@echo "  ${BOLD}enable-service${REGULAR} - Enable $(APP_NAME) service auto-start"
	@echo "  ${BOLD}disable-service${REGULAR} - Disable $(APP_NAME) service auto-start"
	@echo "\n# Debugging"
	@echo "  ${BOLD}debug-remote${REGULAR} - Start remote debugger with build and deploy"
	@echo "  ${BOLD}start-debug${REGULAR} - Only start remote debugger"
	@echo "  ${BOLD}stop-debug${REGULAR} - Stop remote debugger"
	@echo "\n# Diagnostics"
	@echo "  ${BOLD}test-connection${REGULAR} - Test connection to Raspberry Pi"
	@echo "  ${BOLD}logs${REGULAR} - Show $(APP_NAME) service logs"
	@echo "  ${BOLD}logs-tail${REGULAR} - Show last 50 lines of $(APP_NAME) service logs"

.PHONY: build-cli-debug
# Build debug version of CLI (cmd/cli/main.go) for Linux arm64
build-cli-debug: prepare-builder
	@echo "Building debug CLI (cmd/cli/main.go) for linux/arm64..."
	@mkdir -p build
	docker run --rm \
		--platform linux/arm64 \
		--mount type=bind,src=.,dst=/src \
	  -v /tmp/buildkit-cache:/root/.cache/go-build \
	  -v /tmp/buildkit-cache:/go/pkg/mod \
	  --name raspidrum-builder \
	  raspidrum-builder \
		go build -gcflags="all=-N -l" -buildvcs=false -trimpath -o ./build/raspidrum_cli ./cmd/cli/main.go

