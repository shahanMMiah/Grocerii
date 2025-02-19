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

func (t *tappableLabel) Tapped(_ *fyne.PointEvent) {
	fmt.Println("I have been tapped")
}

func (t *tappableLabel) TappedSecondary(_ *fyne.PointEvent) {
	t.CallBack()
	/*
		dialog.NewConfirm(
			fmt.Sprintf("delete %v", t.Text),
			fmt.Sprintf("delete %v", t.Text),
			func(bool) {

			},
			t.Win)
	*/
}

type tappableLabel struct {
	widget.Label
	EntryInd int
	Win      fyne.Window
	CallBack func() bool
}

func NewTabableLabel(t string, i int) *tappableLabel {
	label := &tappableLabel{}
	label.ExtendBaseWidget(label)
	label.SetText(t)
	label.EntryInd = i
	label.Win = fyne.CurrentApp().NewWindow(fmt.Sprintf("%v %v window", label.Text, label.EntryInd))

	return label
}

func MakeIngEntries(ings *[]ingredient) []fyne.CanvasObject {

	cnvs := make([]fyne.CanvasObject, 0)

	for ind, i := range *ings {

		unitSel := widget.NewSelect(unitVals, func(v string) {})
		unitSel.SetSelectedIndex(getUnitInd(i.Unit))
		nameEntry := NewTabableLabel(fmt.Sprintf("%v %c", i.Name, i.Emoji), ind)
		unitEntry := widget.NewEntry()
		unitEntry.SetText(fmt.Sprintf("%v", i.Amount))
		checkBox := widget.NewCheck("", func(b bool) {})
		checkBox.SetChecked(i.Check)

		nameEntry.CallBack = func() bool {

			*ings = RemoveIngredient(*ings, nameEntry.EntryInd)
			return true
		}

		cont := container.New(
			layout.NewCustomPaddedHBoxLayout(3),
			nameEntry,
			unitEntry,
			unitSel,
			checkBox,
		)

		cnvs = append(cnvs, cont)
	}
	return cnvs
}

func DrawIngEntries(e []fyne.CanvasObject, c *fyne.Container) {
	c.RemoveAll()
	for _, entry := range e {
		c.Add(entry)
	}
}

/*
func SetupIngredientEntry(ings *[]ingredient, e []fyne.CanvasObject, c *fyne.Container) {

	for _, entry := range e {
		label := entry.(*fyne.Container).Objects[0].(*tappableLabel)

	}
}*/

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
			*i = AddIngredients(formItem.Widget.(*widget.Entry).Text, 0.0, rune(127816), unt, false, *i)
		},
		w,
	)

	dialog.Show()

}

func UpdateIngEntries(i *[]ingredient, c *fyne.Container, e *[]fyne.CanvasObject) {

	for {
		if len(*i) != len(*e) {
			*e = MakeIngEntries(i)
			DrawIngEntries(*e, c)
		}
	}

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

func SaveIngEntries(c []fyne.CanvasObject, i *[]ingredient, r []recipe, d fyne.URI) {
	fmt.Printf("saving data at %v", d.Path())

	for ind, con := range c {

		rCon := con.(*fyne.Container)

		(*i)[ind].Name = rCon.Objects[0].(*tappableLabel).Text
		n, err := strconv.ParseFloat(rCon.Objects[1].(*widget.Entry).Text, 64)

		if err != nil {
			panic(err)
		}
		(*i)[ind].Amount = n

		unitInd := rCon.Objects[2].(*widget.Select).SelectedIndex()
		if unitInd == -1 {
			unitInd = 0
		}

		(*i)[ind].Unit = unitVals[unitInd]

		(*i)[ind].Check = rCon.Objects[3].(*widget.Check).Checked

	}

	SaveData(d, *i, r)
}

func BuildUI(a fyne.App, w fyne.Window, i *[]ingredient, r []recipe, d fyne.URI) *fyne.Container {

	//w.Resize(fyne.NewSize(500, 500))
	clock := widget.NewLabel("time: ")

	go func() {
		for t := range time.Tick(time.Second) {

			h, m, s := t.Clock()
			clock.SetText(fmt.Sprintf("time : %v : %v : %v", h, m, s))
		}
	}()

	ingContainer := container.NewVBox()

	ingsEntries := MakeIngEntries(i)

	DrawIngEntries(ingsEntries, ingContainer)

	go UpdateIngEntries(i, ingContainer, &ingsEntries)

	addIngsBtn := widget.NewToolbarAction(theme.ContentAddIcon(), func() { AddIngredientEntry(i, ingContainer, w) })

	saveBtn := widget.NewToolbarAction(theme.DocumentSaveIcon(), func() { SaveIngEntries(ingsEntries, i, r, d) })
	toolbar := widget.NewToolbar(addIngsBtn, saveBtn)
	cont := container.NewVBox(
		toolbar,
		clock,
		ingContainer,
	)

	ingTab := container.NewAppTabs(
		container.NewTabItem("Ingredients", cont),
		container.NewTabItem("Reipes", widget.NewLabel("temp")))

	w.SetContent(ingTab)

	return cont

}
