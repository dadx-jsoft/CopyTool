package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
)

// pointerWidget shows a hand cursor over clickable controls.
type pointerWidget struct {
	widget.BaseWidget
	child fyne.CanvasObject
}

func wrapPointer(obj fyne.CanvasObject) fyne.CanvasObject {
	w := &pointerWidget{child: obj}
	w.ExtendBaseWidget(w)
	return w
}

func (p *pointerWidget) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(p.child)
}

func (p *pointerWidget) Cursor() desktop.Cursor {
	return desktop.PointerCursor
}
