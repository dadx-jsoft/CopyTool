package main

import (
	"image/color"
	"strconv"
	"strings"

	"copytool/internal/copy"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

var (
	colorStatus = color.NRGBA{R: 70, G: 70, B: 70, A: 255}
	colorMuted  = color.NRGBA{R: 120, G: 120, B: 120, A: 255}
	colorBlue   = color.NRGBA{R: 0, G: 90, B: 200, A: 255}
	colorRed    = color.NRGBA{R: 180, G: 30, B: 30, A: 255}
)

func main() {
	a := app.NewWithID("com.vn.copytool")
	w := a.NewWindow("Copy files tool")
	w.Resize(fyne.NewSize(560, 400))
	w.SetFixedSize(false)
	w.CenterOnScreen()
	w.SetContent(buildUI(w))
	w.ShowAndRun()
}

func buildUI(w fyne.Window) fyne.CanvasObject {
	fromEntry := widget.NewEntry()
	fromEntry.SetPlaceHolder("Paste folder path or Browse")
	toEntry := widget.NewEntry()
	toEntry.SetPlaceHolder("Paste folder path or Browse")
	extEntry := widget.NewEntry()
	extEntry.SetPlaceHolder("e.g. jpg, png, pdf")
	extEntry.Disable()

	status := canvas.NewText("Status", colorStatus)
	status.Alignment = fyne.TextAlignCenter
	status.TextSize = 13

	setStatus := func(text string, c color.Color) {
		status.Text = text
		status.Color = c
		status.Refresh()
	}

	resetStatus := func() {
		setStatus("Status", colorStatus)
	}

	fromEntry.OnChanged = func(string) { resetStatus() }
	toEntry.OnChanged = func(string) { resetStatus() }
	extEntry.OnChanged = func(string) { resetStatus() }

	mode := widget.NewRadioGroup([]string{"All files", "By extension", "Exclude extension"}, func(selected string) {
		resetStatus()
		if selected == "By extension" || selected == "Exclude extension" {
			extEntry.Enable()
		} else {
			extEntry.Disable()
		}
	})
	mode.SetSelected("All files")
	mode.Horizontal = true

	pickFolder := func(target *widget.Entry) {
		dialog.ShowFolderOpen(func(uri fyne.ListableURI, err error) {
			if err != nil || uri == nil {
				return
			}
			target.SetText(uriToPath(uri))
			resetStatus()
		}, w)
	}

	fromBrowse := widget.NewButton("Browse", func() { pickFolder(fromEntry) })
	toBrowse := widget.NewButton("Browse", func() { pickFolder(toEntry) })
	fromBrowse.Importance = widget.MediumImportance
	toBrowse.Importance = widget.MediumImportance

	var copyBtn *widget.Button
	copyBtn = widget.NewButton("Copy", func() {
		sourceDir := strings.TrimSpace(fromEntry.Text)
		destDir := strings.TrimSpace(toEntry.Text)
		if sourceDir == "" || destDir == "" {
			setStatus("Please select folder", colorRed)
			return
		}

		filter := copy.FilterAll
		var exts []string
		switch mode.Selected {
		case "By extension":
			filter = copy.FilterInclude
			exts = copy.ParseExtensions(extEntry.Text)
			if len(exts) == 0 {
				setStatus("Please enter file extensions", colorRed)
				return
			}
		case "Exclude extension":
			filter = copy.FilterExclude
			exts = copy.ParseExtensions(extEntry.Text)
			if len(exts) == 0 {
				setStatus("Please enter file extensions", colorRed)
				return
			}
		}

		copyBtn.Disable()
		setStatus("Copying", colorStatus)

		go func() {
			n, err := copy.CopyFilesToDirectory(sourceDir, destDir, filter, exts)
			fyne.Do(func() {
				copyBtn.Enable()
				if err != nil || n < 0 {
					setStatus("Copy fail", colorRed)
					return
				}
				setStatus("Copy success "+strconv.FormatInt(n, 10)+" files", colorBlue)
			})
		}()
	})
	copyBtn.Importance = widget.HighImportance

	footer := canvas.NewText("Dương Xuân Đà - 0961010169 - dadx.jsoft@gmail.com", colorMuted)
	footer.Alignment = fyne.TextAlignCenter
	footer.TextSize = 11

	browseSpacer := canvas.NewRectangle(color.Transparent)
	browseSpacer.SetMinSize(fyne.NewSize(fromBrowse.MinSize().Width, 1))

	label := func(text string) *widget.Label {
		l := widget.NewLabel(text)
		l.Alignment = fyne.TextAlignTrailing
		return l
	}

	form := container.New(layout.NewFormLayout(),
		label("From"), container.NewBorder(nil, nil, nil, fromBrowse, fromEntry),
		label("To"), container.NewBorder(nil, nil, nil, toBrowse, toEntry),
		label("Filter"), mode,
		label("Ext"), container.NewBorder(nil, nil, nil, browseSpacer, extEntry),
	)

	sep := canvas.NewLine(color.NRGBA{R: 210, G: 210, B: 210, A: 255})
	sep.StrokeWidth = 1

	copyRow := container.NewGridWithColumns(3,
		layout.NewSpacer(),
		copyBtn,
		layout.NewSpacer(),
	)

	return container.NewPadded(container.NewVBox(
		form,
		widget.NewSeparator(),
		copyRow,
		status,
		sep,
		footer,
	))
}

func uriToPath(uri fyne.ListableURI) string {
	path := uri.Path()
	if len(path) >= 3 && path[0] == '/' && path[2] == ':' {
		return path[1:]
	}
	return path
}
