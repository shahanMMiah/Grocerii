package main

import (
	"fmt"

	"fyne.io/fyne/v2"

	"strconv"

	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

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

func MakeRecEntries(recs *[]recipe) []fyne.CanvasObject {
	cnvs := make([]fyne.CanvasObject, 0)
	for ind, r := range *recs {
		nameEntry := NewTabableLabel(r.Name, ind)
		checkBox := widget.NewCheck("", func(b bool) {})

		nameEntry.CallBack = func() bool {

			*recs = RemoveRecipes(*recs, nameEntry.EntryInd)
			return true
		}

		cont := container.New(
			layout.NewCustomPaddedHBoxLayout(3),
			nameEntry,
			checkBox,
		)
		cnvs = append(cnvs, cont)
	}
	return cnvs

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

func DrawEntries(e []fyne.CanvasObject, c *fyne.Container) {
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

	dialg := dialog.NewForm(
		"add new ingredient",
		"OK",
		"CANCEL",
		forms,
		func(c bool) {
			if c {
				fmt.Println(formItem.Widget.(*widget.Entry).Text)
				*i = AddIngredients(formItem.Widget.(*widget.Entry).Text, 0.0, rune(127816), unt, false, *i)
			}
		},
		w,
	)

	dialg.Show()

}

func AddRecipeEntry(r *[]recipe, c *fyne.Container, w fyne.Window) {

	textInput := widget.NewEntry()

	textInput.SetPlaceHolder("eh?")

	forms := make([]*widget.FormItem, 0)
	formItem := widget.NewFormItem("Recipe Item", textInput)

	forms = append(forms, formItem)

	dialg := dialog.NewForm(
		"add new Recipe",
		"OK",
		"CANCEL",
		forms,
		func(c bool) {

			if c {
				*r = AddRecipe(formItem.Widget.(*widget.Entry).Text, *r)
			}
		},
		w,
	)

	dialg.Show()

}

func UpdateIngEntries(i *[]ingredient, c *fyne.Container, e *[]fyne.CanvasObject) {

	for {
		if len(*i) != len(*e) {
			*e = MakeIngEntries(i)
			DrawEntries(*e, c)
			c.Refresh()
		}
	}

}
func UpdateRecEntries(r *[]recipe, c *fyne.Container, e *[]fyne.CanvasObject) {

	for {
		if len(*r) != len(*e) {
			*e = MakeRecEntries(r)
			DrawEntries(*e, c)
			c.Refresh()
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

func GetIngEntriesData(c []fyne.CanvasObject, i *[]ingredient, d fyne.URI) {
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

}

func GetRecEntriesData(c []fyne.CanvasObject, r *[]recipe, d fyne.URI) {
	fmt.Printf("saving data at %v", d.Path())

	for ind, con := range c {

		rCon := con.(*fyne.Container)

		(*r)[ind].Name = rCon.Objects[0].(*tappableLabel).Text

		(*r)[ind].Check = rCon.Objects[1].(*widget.Check).Checked

	}

}

func BuildUI(a fyne.App, w fyne.Window, i *[]ingredient, r *[]recipe, d fyne.URI) {
	// ingredients
	ingContainer := container.NewVBox()
	ingsEntries := MakeIngEntries(i)
	DrawEntries(ingsEntries, ingContainer)
	go UpdateIngEntries(i, ingContainer, &ingsEntries)

	addIngsBtn := widget.NewToolbarAction(theme.ContentAddIcon(), func() { AddIngredientEntry(i, ingContainer, w) })

	ingSaveBtn := widget.NewToolbarAction(theme.DocumentSaveIcon(), func() {
		GetIngEntriesData(ingsEntries, i, d)
		SaveData(d, *i, *r)
	})
	ingToolbar := widget.NewToolbar(addIngsBtn, ingSaveBtn)
	ingMainCont := container.NewVBox(
		ingToolbar,

		ingContainer,
	)
	ingScroll := container.NewVScroll(ingMainCont)

	// recipes
	recContainer := container.NewVBox()
	recEntries := MakeRecEntries(r)
	DrawEntries(recEntries, recContainer)
	go UpdateRecEntries(r, recContainer, &recEntries)

	addRecBtn := widget.NewToolbarAction(theme.ContentAddIcon(), func() { AddRecipeEntry(r, recContainer, w) })

	saveBtn := widget.NewToolbarAction(theme.DocumentSaveIcon(), func() {
		GetRecEntriesData(recEntries, r, d)
		SaveData(d, *i, *r)
	})
	recToolbar := widget.NewToolbar(addRecBtn, saveBtn)

	recMainCont := container.NewVBox(
		recToolbar,
		recContainer,
	)

	recScroll := container.NewVScroll(recMainCont)

	Tab := container.NewAppTabs(
		container.NewTabItem("Ingredients", ingScroll),
		container.NewTabItem("Reipes", recScroll))

	w.SetContent(Tab)

}
