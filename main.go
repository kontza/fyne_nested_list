package main

import (
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

const CREATE_NEW_ITEM = -1

type buttonWidgets map[Widget]*widget.Button

type WidgetManager struct {
	buttons    buttonWidgets
	listWidget *widget.List
}

type listItem struct {
	name    string
	subList []string
}

type guiApp struct {
	listBinding          binding.ExternalUntypedList
	currentListSelection widget.ListItemID
	listWidget           *widget.List
	items                []listItem
	widgetManager        *WidgetManager
}

//go:generate stringer -type Widget

type Widget int64

const (
	ItemAdd Widget = iota
	ItemRemove
)

type buttons map[Widget]*widget.Button

func (w *WidgetManager) createButton(s Widget) *widget.Button {
	if w.buttons == nil {
		w.buttons = make(map[Widget]*widget.Button)
	}
	btn := widget.NewButton(s.String(), func() {
		log.Printf("Not implemented: %s.OnTapped", s.String())
	})
	w.buttons[s] = btn
	return btn
}

func (w *WidgetManager) getButton(s Widget) *widget.Button {
	return w.buttons[s]
}

func (w *WidgetManager) createListWidget(listBinding binding.ExternalUntypedList) *widget.List {
	w.listWidget = widget.NewListWithData(listBinding,
		func() fyne.CanvasObject {
			w := widget.NewLabel("SPACE ALLOCATION")
			return w
		},
		func(i binding.DataItem, o fyne.CanvasObject) {
			if entry, err := i.(binding.Untyped).Get(); err != nil {
				log.Println("Failed to Get item")
			} else {
				w := o.(*widget.Label)
				w.SetText(entry.(listItem).name)
			}
		})
	return w.listWidget
}

func (w *WidgetManager) getProfileList() *widget.List {
	return w.listWidget
}

func (ga *guiApp) addItem() {
	excludes := []string{"eka", "toka", "kola"}
	newItem := listItem{name: "New Profile", subList: excludes}
	ga.listBinding.Append(newItem)
}

func (ga *guiApp) removeItem() {
	if ga.currentListSelection == CREATE_NEW_ITEM {
		return
	}
	boundItems, err := ga.listBinding.Get()
	if err != nil {
		log.Println("Failed to Get:", err)
	}
	boundItems = append(boundItems[:ga.currentListSelection], boundItems[(1+ga.currentListSelection):]...)
	ga.listBinding.Set(boundItems)
}

func (ga *guiApp) createWindowAndRun() {
	theApp := app.New()
	appWindow := theApp.NewWindow("Add items")
	items := make([]interface{}, len(ga.items))
	for i, p := range ga.items {
		items[i] = p
	}
	ga.listBinding = binding.BindUntypedList(&items)
	ga.listWidget = ga.widgetManager.createListWidget(ga.listBinding)

	itemAdd := ga.widgetManager.createButton(ItemAdd)
	itemAdd.OnTapped = ga.addItem
	itemRemove := ga.widgetManager.createButton(ItemRemove)
	itemRemove.OnTapped = ga.removeItem

	buttonContainer := container.New(layout.NewHBoxLayout(), itemAdd, itemRemove)
	listContainer := container.New(layout.NewMaxLayout(), ga.listWidget)
	mainContainer := container.New(
		layout.NewBorderLayout(nil, buttonContainer, nil, nil),
		listContainer,
		buttonContainer,
	)

	// Combined
	appWindow.SetContent(mainContainer)
	appWindow.Resize(fyne.Size{Height: 320, Width: 480})

	// Handle selection change.
	ga.listWidget.OnSelected = func(i widget.ListItemID) {
		ga.currentListSelection = i
	}
	if len(ga.items) > 0 {
		ga.listWidget.Select(0)
	}
	appWindow.ShowAndRun()
}

func main() {
	ga := guiApp{
		currentListSelection: CREATE_NEW_ITEM,
		items: []listItem{
			{
				name:    "first",
				subList: []string{"sub-1-1", "sub-1-2"},
			},
			{
				name:    "second",
				subList: []string{"sub-2-1", "sub-2-2"},
			},
		},
		widgetManager: &WidgetManager{},
	}
	ga.createWindowAndRun()
}
