package main

import (
	"fmt"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/getlantern/systray"
)

var (
	chatWindow fyne.Window
	a          fyne.App
)

func main() {
	a = app.NewWithID("confidant")
	w := a.NewWindow("SysTray")

	go systray.Run(onReady, onExit)

	w.SetContent(widget.NewLabel("Fyne System Tray"))
	w.SetCloseIntercept(func() {
		w.Hide()
	})
	w.ShowAndRun()
}

func onReady() {
	// systray.SetIcon(getIcon()) // Replace with your icon
	systray.SetTitle("CogentCore Chat")
	systray.SetTooltip("Click to open chat")

	mOpenChat := systray.AddMenuItem("Open Chat", "Open chat window")
	mQuit := systray.AddMenuItem("Quit", "Quit the application")

	go func() {
		for {
			select {
			case <-mOpenChat.ClickedCh:
				showChatWindow()
			case <-mQuit.ClickedCh:
				systray.Quit()
				os.Exit(0)
			}
		}
	}()
}

func onExit() {
	fmt.Println("Exiting application...")
}

func showChatWindow() {
	chatWindow = a.NewWindow("Chat")
	chatWindow.Resize(fyne.NewSize(400, 300))

	input := widget.NewEntry()
	sendButton := widget.NewButton("Send", func() {
		fmt.Println("Message Sent:", input.Text)
		input.SetText("") // Clear input after sending
	})

	chatWindow.SetContent(container.NewVBox(
		input,
		sendButton,
	))
	chatWindow.Show()
}

func getIcon() []byte {
	// Placeholder: Load your icon data here
	return []byte{}
}
