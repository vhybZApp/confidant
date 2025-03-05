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

### Step 1: Configure API Keys  

1. Copy the default environment file:  

   ```sh
   cp infra/default.env infra/.env
   ```  

2. Open `infra/.env` and set your LLM model and API key.  

### Step 2: Set Up the Application  

1. Copy the default app environment file:  

   ```sh
   cp default.env .env
   ```  

2. Open `.env` and configure the following:  
   - **LLM_MODEL**: Set this to match the model you configured in `infra/.env`. Example:  

     ```env
     LLM_MODEL=gemini
     ```  

   - **DEVICE_TYPE**: Set this to match your device. Example:  

     ```env
     DEVICE_TYPE="Gnome Linux"
     ```  

     (For macOS, use `"Mac"`, and adjust accordingly for other systems.)  

### Step 3: Start AI Services  

Run the following command to start the required services:  

```sh
make up
```  

### Step 4: Run Your Request  

Execute the application with a sample request:  

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
