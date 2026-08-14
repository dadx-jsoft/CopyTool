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
	"fyne.io/fyne/v2/widget"
)

const (
	labelColWidth = 52
	fieldGap      = 10
	sectionGap    = 14
)

var (
	colorStatus = color.NRGBA{R: 70, G: 70, B: 70, A: 255}
	colorMuted  = color.NRGBA{R: 130, G: 130, B: 130, A: 255}
	colorBlue   = color.NRGBA{R: 0, G: 90, B: 200, A: 255}
	colorRed    = color.NRGBA{R: 180, G: 30, B: 30, A: 255}
)

func main() {
	a := app.NewWithID("com.vn.copytool")
	w := a.NewWindow("Copy files tool")
	w.Resize(fyne.NewSize(500, 400))
	w.SetFixedSize(true)
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
	mode.Horizontal = false

	pickFolder := func(target *widget.Entry) {
		dialog.ShowFolderOpen(func(uri fyne.ListableURI, err error) {
			if err != nil || uri == nil {
				return
			}
			target.SetText(uriToPath(uri))
			resetStatus()
		}, w)
	}

	fromBrowseBtn := widget.NewButton("Browse", func() { pickFolder(fromEntry) })
	toBrowseBtn := widget.NewButton("Browse", func() { pickFolder(toEntry) })
	fromBrowse := wrapPointer(fromBrowseBtn)
	toBrowse := wrapPointer(toBrowseBtn)

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

	browseWidth := fromBrowseBtn.MinSize().Width
	browseSpacer := canvas.NewRectangle(color.Transparent)
	browseSpacer.SetMinSize(fyne.NewSize(browseWidth, 1))

	pathField := func(entry *widget.Entry, browse fyne.CanvasObject) fyne.CanvasObject {
		return container.NewBorder(nil, nil, nil, browse, entry)
	}

	filterField := wrapPointer(mode)

	rows := container.NewVBox(
		formRow("From", pathField(fromEntry, fromBrowse)),
		vGap(fieldGap),
		formRow("To", pathField(toEntry, toBrowse)),
		vGap(sectionGap),
		formRow("Filter", filterField),
		vGap(fieldGap),
		formRow("Ext", container.NewBorder(nil, nil, nil, browseSpacer, extEntry)),
	)

	copyBtn.Resize(fyne.NewSize(140, 36))
	copyRow := container.NewCenter(wrapPointer(copyBtn))

	content := container.NewVBox(
		rows,
		vGap(sectionGap),
		widget.NewSeparator(),
		vGap(fieldGap),
		copyRow,
		vGap(fieldGap),
		status,
		vGap(sectionGap),
		widget.NewSeparator(),
		vGap(6),
		footer,
	)

	return container.NewPadded(content)
}

func formRow(label string, field fyne.CanvasObject) fyne.CanvasObject {
	lbl := widget.NewLabel(label)
	lbl.Alignment = fyne.TextAlignTrailing
	labelBox := container.NewGridWrap(fyne.NewSize(labelColWidth, lbl.MinSize().Height))
	labelBox.Add(lbl)
	return container.NewBorder(nil, nil, labelBox, nil, field)
}

func vGap(h float32) fyne.CanvasObject {
	gap := canvas.NewRectangle(color.Transparent)
	gap.SetMinSize(fyne.NewSize(1, h))
	return gap
}

func uriToPath(uri fyne.ListableURI) string {
	path := uri.Path()
	if len(path) >= 3 && path[0] == '/' && path[2] == ':' {
		return path[1:]
	}
	return path
}
