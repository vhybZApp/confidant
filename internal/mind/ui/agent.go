package ui

import (
	"fyne.io/fyne/v2/widget"
)

type AgentComponent struct {
	msgs []*MessageComponent
}

type MessageComponent struct {
	from    *widget.Label
	content *widget.Label
}
