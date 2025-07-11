#!/bin/bash
# Install Go
GO_VERSION=1.24.4
GO_ARCH=arm64
wget https://go.dev/dl/go${GO_VERSION}.linux-${GO_ARCH}.tar.gz
sudo tar -C /usr/local -xzf go${GO_VERSION}.linux-${GO_ARCH}.tar.gz
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
source ~/.bashrc

# Install delve
go install github.com/go-delve/delve/cmd/dlv@latest

# Add delve to PATH
echo 'export PATH=$PATH:~/go/bin' >> ~/.bashrc
source ~/.bashrc

sudo apt install tmux -y
