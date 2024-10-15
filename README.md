# UBX-GPSD Bridge
---

Fawad Mazhar <fawadmazhar@hotmail.com> 2024

---


## Overview

UBX-GPSD Bridge is a Go application that acts as an intermediary between UBX-protocol GPS devices connected via I2C and GPSD (GPS daemon) clients connected via TCP. This bridge allows GPSD to work with GPS modules that use the UBX protocol and are connected through I2C, which is not natively supported by GPSD.

## Features

- Reads UBX protocol data from an I2C-connected GPS device
- Serves GPS data to GPSD clients over a TCP connection
- Handles bidirectional communication, allowing GPSD to send commands to the GPS device
- Supports both UBX and NMEA sentence parsing

## Prerequisites

- Go 1.16 or higher
- Access to an I2C bus (usually available on embedded systems like Raspberry Pi)
- A UBX-protocol compatible GPS module connected via I2C
- GPSD installed on your system (for client connections)

## Installation

1. Clone the repository:
   ```
   git clone https://github.com/fawad1985/ubx-gpsd-bridge.git
   cd ubx-gpsd-bridge
   ```

2. Install the required Go packages:
   ```
   go mod tidy
   ```

3. Build the application:
   ```
   go build -o ubx-gpsd-bridge ./cmd/ubxgpsdbridge
   ```

## Usage

1. Connect your UBX-compatible GPS module to the I2C bus of your system.

2. Run the UBX-GPSD Bridge:
   ```
   ./ubx-gpsd-bridge
   ```

   By default, the bridge will listen on `127.0.0.1:49000` for GPSD client connections.

3. Configure GPSD to connect to the bridge:
   - Edit your GPSD configuration file (usually `/etc/default/gpsd` or `/etc/gpsd/gpsd.conf`)
   - Set the device to: `tcp://127.0.0.1:49000`

4. Restart GPSD service:
   ```
   sudo systemctl restart gpsd
   ```

5. You can now use GPSD clients (like `gpsmon`, `cgps`, or custom applications) to access the GPS data.

## Configuration

The main configuration options are in `cmd/ubxgpsdbridge/main.go`:

- I2C bus number (default: 1)
- I2C device address (default: 0x42)
- TCP listen address (default: "127.0.0.1:49000")

Modify these values as needed for your specific setup.

## Project Structure

```
project/
├── cmd/
│   └── ubxgpsdbridge/
│       └── main.go
├── pkg/
│   ├── ubx/
│   │   ├── packet.go
│   │   └── i2cdevice.go
│   ├── gpsd/
│   │   └── feeder.go
│   └── tcp/
│       └── handler.go
├── internal/
│   └── utils/
│       └── utils.go
└── go.mod
```

- `cmd/ubxgpsdbridge`: Contains the main application entry point.
- `pkg/ubx`: Handles UBX protocol parsing and I2C communication.
- `pkg/gpsd`: Manages the GPSD feeder functionality.
- `pkg/tcp`: Handles TCP connections from GPSD clients.
- `internal/utils`: Contains utility functions used across the project.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
