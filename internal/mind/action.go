package mind

import (
	"errors"
	"time"

	"github.com/go-vgo/robotgo"
)

func ExecAction(a Action, ai AnnotatedImage) error {
	switch a.NextAction {
	case "left_click":
		x, y := ai.BoundedBox(a.BoxID)
		robotgo.MoveSmooth(x, y)
		robotgo.Click()
	case "right_click":
		x, y := ai.BoundedBox(a.BoxID)
		robotgo.MoveSmooth(x, y)
		robotgo.Click("right")
	case "double_click":
		x, y := ai.BoundedBox(a.BoxID)
		robotgo.MoveSmooth(x, y)
		robotgo.Click("left", true)
	case "hover":
		x, y := ai.BoundedBox(a.BoxID)
		robotgo.MoveSmooth(x, y)
	case "scroll_up":
		robotgo.ScrollDir(10, "up")
	case "scroll_down":
		robotgo.ScrollDir(10, "down")
	case "type":
		x, y := ai.BoundedBox(a.BoxID)
		robotgo.MoveSmooth(x, y)
		robotgo.Click()
		robotgo.TypeStr(a.Value)
		// robotgo.KeyPress("enter")
	case "wait":
		time.Sleep(time.Second * 3)
	default:
		return errors.New("action not found")
	}
	return nil
}
