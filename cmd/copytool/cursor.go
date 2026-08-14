package main

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
)

// pointerLayer sits above a control and shows a hand cursor while forwarding taps.
type pointerLayer struct {
	widget.BaseWidget
	target fyne.CanvasObject
}

func stackPointer(obj fyne.CanvasObject) fyne.CanvasObject {
	layer := &pointerLayer{target: obj}
	layer.ExtendBaseWidget(layer)
	return container.NewStack(obj, layer)
}

func (p *pointerLayer) CreateRenderer() fyne.WidgetRenderer {
	hit := canvas.NewRectangle(color.NRGBA{0, 0, 0, 0})
	return widget.NewSimpleRenderer(hit)
}

func (p *pointerLayer) Cursor() desktop.Cursor {
	if btn, ok := p.target.(*widget.Button); ok && btn.Disabled() {
		return desktop.DefaultCursor
	}
	return desktop.PointerCursor
}

func (p *pointerLayer) Tapped(ev *fyne.PointEvent) {
	if t, ok := p.target.(fyne.Tappable); ok {
		t.Tapped(ev)
	}
}

func (p *pointerLayer) TappedSecondary(ev *fyne.PointEvent) {
	if t, ok := p.target.(fyne.SecondaryTappable); ok {
		t.TappedSecondary(ev)
	}
}

// wrapPointer adds a hand cursor for non-button widgets such as RadioGroup.
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
