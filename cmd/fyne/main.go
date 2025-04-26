package main

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"github.com/farhoud/confidant/internal/mind/ui"
	"github.com/farhoud/confidant/pkg/fact"
)

func main() {
	a := app.New()
	w := a.NewWindow("SysTray")

	nodesState := fact.NewReactiveListState[ui.NodeData]()
	nodesState.Add(ui.NodeData{Text: "sege", Color: color.RGBA{100, 100, 244, 100}})
	nodesState.Add(ui.NodeData{Text: "soote", Color: color.RGBA{200, 100, 244, 100}})
	flow := ui.NewFlowComponent(nodesState)
	content := container.New(layout.NewCenterLayout(), flow.Render())
	// nodes.Add(ui.NodeData{Text: "root", color: color.RGBA{0, 0, 255, 255}})
	// nodes.Add(ui.NodeData{"sege", color.RGBA{0, 255, 255, 100}})
	w.SetContent(content)
	w.Resize(fyne.NewSize(100, 100))
	w.ShowAndRun()
}
