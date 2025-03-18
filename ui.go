package main

import (
	"fmt"

	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// ----------- Search Bar --------------
type SearchEntry struct {
	widget.Entry
	Update bool
}

func (s *SearchEntry) HighlightSearch(objs *[]fyne.CanvasObject, items Groceitem, t *Trie) {

	found := t.AutoComplete(s.Text)
	list := *objs

	if ings, ok := items.(*ingredients); ok {
		for iter := range list {
			nameEntry := list[iter].(*fyne.Container).Objects[1].(*TappableLabel)
			ings.Ingredients[iter].Highlighted = false

			for _, f := range found {
				if f == nameEntry.Segments[0].(*widget.TextSegment).Text {

					ings.Ingredients[iter].Highlighted = true
				}

			}
			nameEntry.Refresh()

		}
		ings.Update = true

	}

}

func DrawHighlights(ings *ingredients, objs *[]fyne.CanvasObject) {

	list := *objs

	for iter := range list {
		nameEntry := list[iter].(*fyne.Container).Objects[1].(*TappableLabel)
		nameEntry.Segments[0].(*widget.TextSegment).Style.ColorName = nameEntry.Color

		if ings.Ingredients[iter].Highlighted {
			nameEntry.Segments[0].(*widget.TextSegment).Style.ColorName = theme.ColorNameError

		}
		nameEntry.Refresh()

	}

}

func NewSearchEntry() *SearchEntry {
	entry := &SearchEntry{}
	entry.ExtendBaseWidget(entry)
	return entry
}

// ----------- tapable label --------------

type TappableLabel struct {
	widget.RichText
	EntryInd int
	Win      fyne.Window
	CallBack func() bool
	Color    fyne.ThemeColorName
}

func (t *TappableLabel) TappedSecondary(_ *fyne.PointEvent) {
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

func NewTabableLabel(t string, i int) *TappableLabel {
	label := &TappableLabel{}
	label.ExtendBaseWidget(label)
	label.Wrapping = fyne.TextWrapOff
	label.Scroll = 3
	label.Segments = append(label.Segments, &widget.TextSegment{
		Style: widget.RichTextStyleEmphasis,
		Text:  t})

	label.EntryInd = i
	label.Color = theme.ColorNameForeground
	label.Win = fyne.CurrentApp().NewWindow(fmt.Sprintf("%v %v window", label.Segments[0].(*widget.TextSegment).Text, label.EntryInd))
	return label
}

func (t *TappableLabel) SetText(s string) {
	t.Segments[0].(*widget.TextSegment).Text = s
}

func (t *TappableLabel) GetText() string {
	return t.Segments[0].(*widget.TextSegment).Text
}

// ----------- entries --------------

func MakeRecEntries(recs *recipes) []fyne.CanvasObject {
	cnvs := make([]fyne.CanvasObject, 0)
	for ind, r := range recs.Recipes {
		nameEntry := NewTabableLabel(r.Name, ind)

		checkBox := widget.NewCheck("", func(b bool) {

		})

		nameEntry.CallBack = func() bool {

			recs.Remove(nameEntry.EntryInd)
			return true
		}
		checkBox.SetChecked(r.Check)
		cont := container.New(
			layout.NewCustomPaddedHBoxLayout(3),
			nameEntry,
			checkBox,
		)
		cnvs = append(cnvs, cont)
	}
	return cnvs

}

func MakeIngEntries(ings *ingredients) []fyne.CanvasObject {

	cnvs := make([]fyne.CanvasObject, 0)

	for ind, i := range ings.Ingredients {

		unitSel := widget.NewSelect(unitVals, func(v string) {})
		unitSel.SetSelectedIndex(getUnitInd(i.Unit))
		nameEntry := NewTabableLabel(fmt.Sprintf("%v", i.Name), ind)
		unitEntry := widget.NewEntry()
		unitEntry.SetText(fmt.Sprintf("%v", i.Amount))
		checkBox := widget.NewCheck("", func(b bool) {

			if b {
				nameEntry.Color = theme.ColorNameDisabled
				ings.Ingredients[ind].Check = true

			} else {
				nameEntry.Color = theme.ColorNameForeground
				ings.Ingredients[ind].Check = false

			}

			nameEntry.Segments[0].(*widget.TextSegment).Style.ColorName = nameEntry.Color
			ings.Update = true

		})
		checkBox.SetChecked(i.Check)

		nameEntry.CallBack = func() bool {

			ings.Remove(nameEntry.EntryInd)
			return true
		}

		cont := container.New(
			layout.NewCustomPaddedHBoxLayout(3),
			checkBox,
			nameEntry,
			unitEntry,
			unitSel,
		)

		cnvs = append(cnvs, cont)
	}
	return cnvs
}

func DrawEntries(e []fyne.CanvasObject, c *fyne.Container) {
	c.RemoveAll()
	/*
		list := widget.NewList(
			func() int { return len(e) },
			func() { return

			}
		)
	*/
	for _, entry := range e {
		c.Add(entry)
	}
}

func AddIngredientEntry(i *ingredients, c *fyne.Container, w fyne.Window) {

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
				////fmt.Println(formItem.Widget.(*widget.Entry).Text)
				i.Add(formItem.Widget.(*widget.Entry).Text)

			}
		},
		w,
	)

	dialg.Show()

}

func AddRecipeEntry(r *recipes, c *fyne.Container, w fyne.Window) {

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
				r.Add(formItem.Widget.(*widget.Entry).Text)

			}
		},
		w,
	)

	dialg.Show()

}

func UpdateIngEntries(i *ingredients, c *fyne.Container, e *[]fyne.CanvasObject, t *Trie) {

	for {
		if len(i.Ingredients) != len(*e) || i.Update {
			t.build(i)
			i.CheckSort()
			i.HighlightSort()
			*e = MakeIngEntries(i)
			DrawEntries(*e, c)
			DrawHighlights(i, e)
			c.Refresh()
			i.Update = false
		}
	}

}

func UpdateRecEntries(r *recipes, c *fyne.Container, e *[]fyne.CanvasObject) {

	for {
		if len(r.Recipes) != len(*e) {
			*e = MakeRecEntries(r)
			DrawEntries(*e, c)
			c.Refresh()
		}
	}

}

// ----------- data IO --------------

func GetDataURI(a fyne.App) fyne.URI {
	dataURI, urErr := storage.Child(a.Storage().RootURI(), "data.json")
	if urErr != nil {
		panic(urErr)
	}

	return dataURI

}

func SaveData(ur fyne.URI, i *ingredients, r *recipes) {
	jsonOBJ := CreateJson(*i, *r)

	writer, err := storage.Writer(ur)
	if err != nil {
		panic(err)
	}

	writer.Write(jsonOBJ)

}

func ReadData(ur fyne.URI) []byte {
	data, rErr := storage.LoadResourceFromURI(ur)

	if rErr != nil {
		//fmt.Println("no data file found using default")
		data = resourceDataJson
	}

	return data.Content()
}

func GetIngEntriesData(c []fyne.CanvasObject, i *ingredients, d fyne.URI) {
	//fmt.Printf("saving data at %v", d.Path())

	for ind, con := range c {

		rCon := con.(*fyne.Container)

		i.Ingredients[ind].Name = rCon.Objects[1].(*TappableLabel).GetText()
		n, err := strconv.ParseFloat(rCon.Objects[2].(*widget.Entry).Text, 64)

		if err != nil {
			panic(err)
		}
		i.Ingredients[ind].Amount = n

		unitInd := rCon.Objects[3].(*widget.Select).SelectedIndex()
		if unitInd == -1 {
			unitInd = 0
		}

		i.Ingredients[ind].Unit = unitVals[unitInd]

		i.Ingredients[ind].Check = rCon.Objects[0].(*widget.Check).Checked

	}

}

func GetRecEntriesData(c []fyne.CanvasObject, r *recipes, d fyne.URI) {
	//fmt.Printf("saving data at %v", d.Path())

	for ind, con := range c {

		rCon := con.(*fyne.Container)

		r.Recipes[ind].Name = rCon.Objects[0].(*TappableLabel).GetText()

		r.Recipes[ind].Check = rCon.Objects[1].(*widget.Check).Checked

	}

}

// ----------- Main --------------

func BuildUI(a fyne.App, w fyne.Window, i *ingredients, r *recipes, d fyne.URI) {
	// ingredients
	ingSearch := Trie{}
	recSearch := Trie{}

	ingSearch.build(i)
	recSearch.build(r)

	ingContainer := container.NewVBox()
	ingsEntries := MakeIngEntries(i)
	DrawEntries(ingsEntries, ingContainer)
	go UpdateIngEntries(i, ingContainer, &ingsEntries, &ingSearch)

	addIngsBtn := widget.NewToolbarAction(theme.ContentAddIcon(), func() { AddIngredientEntry(i, ingContainer, w) })

	ingSaveBtn := widget.NewToolbarAction(theme.DocumentSaveIcon(), func() {
		GetIngEntriesData(ingsEntries, i, d)
		SaveData(d, i, r)
	})
	ingToolbar := widget.NewToolbar(addIngsBtn, ingSaveBtn)

	ingSearchBar := NewSearchEntry()

	ingSearchBar.OnChanged = func(string) {
		ingSearchBar.HighlightSearch(&ingsEntries, i, &ingSearch)
	}

	ingMainCont := container.NewVBox(
		ingToolbar,
		ingSearchBar,
	)

	ingScroll := container.NewVScroll(ingContainer)
	ingScroll.SetMinSize(
		fyne.NewSize(WINSIZEX, WINSIZEY-120))

	ingMainCont.Add(ingScroll)

	// recipes
	recContainer := container.NewVBox()
	recEntries := MakeRecEntries(r)
	DrawEntries(recEntries, recContainer)
	go UpdateRecEntries(r, recContainer, &recEntries)

	addRecBtn := widget.NewToolbarAction(theme.ContentAddIcon(), func() { AddRecipeEntry(r, recContainer, w) })

	saveBtn := widget.NewToolbarAction(theme.DocumentSaveIcon(), func() {
		GetRecEntriesData(recEntries, r, d)
		SaveData(d, i, r)
	})
	recToolbar := widget.NewToolbar(addRecBtn, saveBtn)

	recMainCont := container.NewVBox(
		recToolbar,
		recContainer,
	)

	recScroll := container.NewVScroll(recMainCont)

	Tab := container.NewAppTabs(
		container.NewTabItem("Ingredients", ingMainCont),
		container.NewTabItem("Reipes", recScroll))

	w.SetContent(Tab)

}
