# Confidant

Confidant is a lightweight automation tool designed to bridge the gap between static software interfaces and the dynamic real world. It enables applications to execute actions autonomously, interacting with UI elements, peripherals, and external devices.

- Automates UI interactions (mouse, keyboard, and more)
- Executes actions based on natural language commands
- Works across various platforms, from desktop environments to edge devices
- Designed for minimal resource usage, written in Go

## Demo

![Confidant Spotify](https://raw.githubusercontent.com/TeleportAgents/assets/main/ConfidantSpotify.gif)

## Installation

### Prerequisites

- Go 1.18+ installed on your system
- make
- docker
- docker-compose

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
go run ./cmd/main.go "search latest infected mushroom track using browser and then play it in spotify app"
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
