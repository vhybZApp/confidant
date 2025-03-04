# Confidant

Confidant is a lightweight automation tool designed to bridge the gap between static software interfaces and the dynamic real world. It enables applications to execute actions autonomously, interacting with UI elements, peripherals, and external devices.

## Features

- Automates UI interactions (mouse, keyboard, and more)
- Executes actions based on natural language commands
- Works across various platforms, from desktop environments to edge devices
- Designed for minimal resource usage, written in Go

## Installation

### Prerequisites
- Go 1.18+ installed on your system

### Clone and Build
```sh
git clone https://github.com/your-org/confidant.git
cd confidant
go build -o confidant
```

## Usage

Run Confidant as a lightweight server/daemon:
```sh
./confidant
```

You can then send natural language commands via API or CLI.

## Configuration

Modify the configuration file (`config.yaml`) to customize behavior:
```yaml
server:
  port: 8080
  log_level: info
```

## API Example

To send an action request:
```sh
curl -X POST http://localhost:8080/action -d '{"command": "open browser"}' -H "Content-Type: application/json"
```

## Contributing

We welcome contributions! To get started:
1. Fork the repository
2. Create a new branch
3. Commit your changes
4. Open a pull request

## License

This project is licensed under the MIT License.

---

Confidant is part of the Teleport Agents ecosystem, integrating seamlessly with AI-driven automation and knowledge capture.
