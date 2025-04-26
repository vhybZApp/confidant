package ui

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/farhoud/confidant/pkg/fact"
)

type NodeData struct {
	Text  string
	Color color.Color
}

type FlowComponent struct {
	nodeListState *fact.ReactiveListState[NodeData]
	root          *fyne.Container
}

func NewFlowComponent(state *fact.ReactiveListState[NodeData]) *FlowComponent {
	fc := &FlowComponent{nodeListState: state, root: container.New(layout.NewStackLayout())}

	fc.renderNodes(state.Get())

	state.Subscribe(func(nodes []NodeData) {
		fc.renderNodes(nodes)
	})

	return fc
}

func (fc *FlowComponent) renderNodes(nodes []NodeData) {
	fc.root.Objects = nil
	for _, data := range nodes {
		circle := canvas.NewCircle(data.Color)
		text := widget.NewLabel(data.Text)
		container := container.New(layout.NewStackLayout(), circle, text)
		fc.root.Add(container)
		circle.Resize(fyne.NewSize(30, 30))
	}
	fc.root.Refresh()
}

func (fc *FlowComponent) Render() fyne.CanvasObject {
	return fc.root
}
