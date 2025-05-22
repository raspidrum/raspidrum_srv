# RaspiDrum Server - Basic Structure

## Overview
RaspiDrum Server is a Go-based service that manages drum kit presets and provides MIDI control functionality through gRPC interfaces.

## Directory Structure

### `/api`
Contains API definitions:
- `/grpc` - Protocol Buffer definitions for gRPC services (preset and channel control)

### `/cmd`
Application entry points:
- `/server` - Main server application that initializes and runs the gRPC server

### `/configs`
Configuration files:
- `dev.yaml` - Development configuration (server host, port, data directory)

### `/internal`
Core application code:

#### `/app`
Business logic implementations:
- `/channel_control` - channel control functionality
- `/mididevice` - MIDI device management
- `/preset` - Preset loading and management
- `/load_kit` - Kit loading functionality

#### `/model`
Domain models:
- Preset, Kit, Instrument, Layer, and Control data structures

#### `/repo`
Data access layer:
- `/db` - SQLite database operations for preset storage
- `/linuxsampler` - Interface with LinuxSampler for audio processing

#### `/pkg`
Shared internal packages:
- `/grpc` - Generated gRPC code

### `/libs`
External library integrations:
- `/liblscp-go` - Go bindings for LinuxSampler Control Protocol

### `/db`
Database-related files:
- SQLite database files
- Schema definitions

## Key Components

### Database (SQLite)
- Stores preset configurations
- Manages kit and instrument relationships

### gRPC Services
1. **Preset Service**
   - Loading presets
   - Managing preset configurations

2. **Channel Control Service**
   - Real-time MIDI control
   - Channel parameter adjustments

### LinuxSampler Integration
- Audio engine management
- Sample playback

## Configuration
The server uses Viper for configuration management with the following key settings:
- Host address and port
- Data directory location
- Database connection parameters 