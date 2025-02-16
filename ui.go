package main

import (
	"fmt"
	"time"

	"fyne.io/fyne/v2"

	"strconv"

	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type ingerdientEntry struct {
	nameField  widget.Entry
	checkFeild widget.Check
	unitFeild  widget.Select
}

func MakeIngEntries(ings []ingredient) []fyne.CanvasObject {

	cnvs := make([]fyne.CanvasObject, 0)

	for _, i := range ings {

		unitSel := widget.NewSelect(unitVals, func(v string) {})
		nameEntry := widget.NewLabel(fmt.Sprintf("%v %c", i.Name, i.Emoji))
		unitEntry := widget.NewEntry()
		unitEntry.SetText(fmt.Sprintf("%v", i.Amount))
		checkBox := widget.NewCheck("", func(b bool) {})

		cont := container.New(
			layout.NewCustomPaddedHBoxLayout(40),
			nameEntry,
			unitEntry,
			unitSel,
			checkBox,
		)

		cnvs = append(cnvs, cont)
	}
	return cnvs
}

func drawIngEntries(e []fyne.CanvasObject, c *fyne.Container) {
	c.RemoveAll()
	for _, entry := range e {
		c.Add(entry)
	}
}

func setIngredientEntry(l fyne.CanvasObject, i *ingredient) {

	i.Name = l.(*fyne.Container).Objects[0].(*widget.Label).Text
	n, err := strconv.ParseFloat(l.(*fyne.Container).Objects[1].(*widget.Entry).Text, 64)
	if err != nil {
		panic(err)
	}
	i.Amount = n
	i.Unit = unitVals[l.(*fyne.Container).Objects[2].(*widget.Select).SelectedIndex()]

	//l.(*fyne.Container).Objects[0].(*widget.Label).SetText(fmt.Sprintf("%v %c", i.Name, i.Emoji))
	//l.(*fyne.Container).Objects[1].(*widget.Entry).SetText(fmt.Sprintf("%v", i.Amount))
	//l.(*fyne.Container).Objects[2].(*widget.Select).SetSelectedIndex(0)

}

func AddIngredientEntry(i *[]ingredient, c *fyne.Container, w fyne.Window) {

	textInput := widget.NewEntry()

	textInput.SetPlaceHolder("eh?")

	forms := make([]*widget.FormItem, 0)
	formItem := widget.NewFormItem("foodName", textInput)

	forms = append(forms, formItem)

	dialog := dialog.NewForm(
		"add new ingredient",
		"OK",
		"CANCEL",
		forms,
		func(bool) {
			fmt.Println(formItem.Widget.(*widget.Entry).Text)

			*i = AddIngredients(formItem.Widget.(*widget.Entry).Text, 0.0, rune(127816), unt, *i)
			ingsEntries := MakeIngEntries(*i)
			drawIngEntries(ingsEntries, c)

		},
		w,
	)

	dialog.Show()

}

func buildUI(a fyne.App, w fyne.Window, i []ingredient, r []recipe, d fyne.URI) *fyne.Container {

	w.Resize(fyne.NewSize(500, 500))
	dataUR := GetDataURI(a)
	clock := widget.NewLabel("time: ")

	go func() {
		for t := range time.Tick(time.Second) {

			h, m, s := t.Clock()
			clock.SetText(fmt.Sprintf("time : %v : %v : %v", h, m, s))
		}
	}()

	ingContainer := container.NewVBox()
	addIngsBtn := widget.NewToolbarAction(theme.ContentAddIcon(), func() { AddIngredientEntry(&i, ingContainer, w) })

	top := widget.NewLabel(fmt.Sprintf("%v", dataUR.Path()))
	ingsEntries := MakeIngEntries(i)

	drawIngEntries(ingsEntries, ingContainer)

	toolbar := widget.NewToolbar(addIngsBtn)
	cont := container.NewVBox(
		toolbar,
		top,
		clock,
		ingContainer,
	)

	/*
		list := widget.NewList(
			func() int { return len(i) },
			func() fyne.CanvasObject { return createEmptyEntry(&ingind, i) },
			func(ind widget.ListItemID, o fyne.CanvasObject) {
				setIngredientEntry(o, &i[ind])
			})

	*/
	//list.Resize(fyne.NewSize(500, 100))
	//cont.Add(list)
	//saveBtn := widget.NewButton("save", func() { Testsave(list, i, r, d) })
	//cont.Add(saveBtn)

	w.SetContent(cont)
	return cont

}

func GetDataURI(a fyne.App) fyne.URI {
	dataURI, urErr := storage.Child(a.Storage().RootURI(), "data.json")
	if urErr != nil {
		panic(urErr)
	}

	return dataURI

}

func SaveData(ur fyne.URI, i []ingredient, r []recipe) {
	jsonOBJ := CreateJson(i, r)

	writer, err := storage.Writer(ur)
	if err != nil {
		panic(err)
	}

	writer.Write(jsonOBJ)

}

func ReadData(ur fyne.URI) []byte {
	data, rErr := storage.LoadResourceFromURI(ur)

	if rErr != nil {
		fmt.Println("no data file found using default")
		data = resourceDataJson
	}

	return data.Content()
}

func Testsave(c *widget.List, i []ingredient, r []recipe, d fyne.URI) {
	fmt.Printf("saving data at %v", d.Path())

	SaveData(d, i, r)
}
