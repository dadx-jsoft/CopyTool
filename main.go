package main

import (
	"image/color"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

var (
	colorStatus = color.NRGBA{R: 0, G: 0, B: 0, A: 255}
	colorBlue   = color.NRGBA{R: 0, G: 0, B: 255, A: 255}
	colorRed    = color.NRGBA{R: 255, G: 0, B: 0, A: 255}
)

func main() {
	a := app.NewWithID("com.vn.copytool")
	w := a.NewWindow("Copy files tool")
	w.Resize(fyne.NewSize(520, 420))
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
	status.TextSize = 14

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
		// Native folder dialog on Windows, macOS, and Linux (Fyne + CGO).
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

	var copyBtn *widget.Button
	copyBtn = widget.NewButton("Copy", func() {
		sourceDir := strings.TrimSpace(fromEntry.Text)
		destDir := strings.TrimSpace(toEntry.Text)
		if sourceDir == "" || destDir == "" {
			setStatus("Please select folder", colorRed)
			return
		}

		filter := filterAll
		var exts []string
		switch mode.Selected {
		case "By extension":
			filter = filterInclude
			exts = parseExtensions(extEntry.Text)
			if len(exts) == 0 {
				setStatus("Please enter file extensions", colorRed)
				return
			}
		case "Exclude extension":
			filter = filterExclude
			exts = parseExtensions(extEntry.Text)
			if len(exts) == 0 {
				setStatus("Please enter file extensions", colorRed)
				return
			}
		}

		copyBtn.Disable()
		setStatus("Copying", colorStatus)

		go func() {
			n, err := copyFilesToDirectory(sourceDir, destDir, filter, exts)
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

	footer := canvas.NewText("Dương Xuân Đà - 0961010169 - dadx.jsoft@gmail.com", colorStatus)
	footer.Alignment = fyne.TextAlignCenter
	footer.TextSize = 12

	pathRow := func(label string, entry *widget.Entry, browse *widget.Button) fyne.CanvasObject {
		return container.NewBorder(nil, nil, widget.NewLabel(label), browse, entry)
	}

	filterRow := container.New(layout.NewFormLayout(),
		widget.NewLabel("Filter"), mode,
		widget.NewLabel("Ext"), extEntry,
	)

	return container.NewPadded(container.NewVBox(
		pathRow("From", fromEntry, fromBrowse),
		pathRow("To", toEntry, toBrowse),
		filterRow,
		container.NewCenter(copyBtn),
		status,
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
