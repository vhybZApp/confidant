# Confidant

Confidant is a lightweight automation tool designed to bridge the gap between static software interfaces and the dynamic real world. It enables applications to execute actions autonomously, interacting with UI elements, peripherals, and external devices.

## Demo

![image](https://drive.google.com/uc?export=view&id=1daQctDIyf0xhRMkDn5de5wvO8shN0D5G)

- Automates UI interactions (mouse, keyboard, and more)
- Executes actions based on natural language commands
- Works across various platforms, from desktop environments to edge devices
- Designed for minimal resource usage, written in Go

## Installation

### Prerequisites

- Go 1.18+ installed on your system
- make

### Clone and Build

```sh
git clone https://github.com/TeleportAgents/confidant.git
cd confidant
make build
```

## Usage

Start AI Services:

```sh
make up
```

Run your request:

```sh
./bin/confidant [request]
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
