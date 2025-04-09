package main

import (
	"fmt"

	"image/color"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// ----------- Theme --------------
type Theme interface {
	Color(fyne.ThemeColorName, fyne.ThemeVariant) color.Color
	Font(fyne.TextStyle) fyne.Resource
	Icon(fyne.ThemeIconName) fyne.Resource
	Size(fyne.ThemeSizeName) float32
}

type CustomTheme struct{}

var _ fyne.Theme = (*CustomTheme)(nil)

type CustomColor struct {
	r, g, b, a uint32
}

func (col CustomColor) RGBA() (r, g, b, a uint32) {
	r = uint32(0xffff * (float64(col.r) / 255.0))
	g = uint32(0xffff * (float64(col.g) / 255.0))
	b = uint32(0xffff * (float64(col.b) / 255.0))
	a = uint32(0xffff * (float64(col.a) / 255.0))

	return r, g, b, a

}

func (m CustomTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	if name == theme.ColorNameBackground {

		return CustomColor{r: 130, g: 10, b: 255, a: 0}
	}

	return theme.DefaultTheme().Color(name, variant)
}
func (m CustomTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(name)
}

func (m CustomTheme) Font(style fyne.TextStyle) fyne.Resource {
	return theme.DefaultTheme().Font(style)
}

func (m CustomTheme) Size(name fyne.ThemeSizeName) float32 {
	return theme.DefaultTheme().Size(name)
}

// ----------- Search Bar --------------
type SearchEntry struct {
	widget.Entry
	Update bool
}

func (s *SearchEntry) TappedSecondary(_ *fyne.PointEvent) {

}

func (s *SearchEntry) HighlightSearch(objs *[]fyne.CanvasObject, items Groceitem, t *Trie) {

	found := t.AutoComplete(s.Text)

	//fmt.Println(found)
	list := *objs

	if ings, ok := items.(*ingredients); ok {
		for iter := range list {
			nameEntry := list[iter].(*fyne.Container).Objects[1].(*TappableLabel)
			ings.Ingredients[iter].Highlighted = false

			for _, f := range found {
				if f == sanatize_string(nameEntry.Segments[0].(*widget.TextSegment).Text) {

					ings.Ingredients[iter].Highlighted = true
				}

			}
			//nameEntry.Refresh()

		}
		ings.Update = true

	}

	if recs, ok := items.(*recipes); ok {

		for iter := range list {
			nameEntry := list[iter].(*fyne.Container).Objects[1].(*TappableLabel)
			recs.Recipes[iter].Highlighted = false

			for _, f := range found {
				if f == sanatize_string(nameEntry.Segments[0].(*widget.TextSegment).Text) {

					recs.Recipes[iter].Highlighted = true
				}

			}
			//Re.Refresh()

		}
		recs.Update = true
	}

}

func DrawHighlights(items Groceitem, objs *[]fyne.CanvasObject) {

	list := *objs
	if ings, ok := items.(*ingredients); ok {
		for iter := range list {
			nameEntry := list[iter].(*fyne.Container).Objects[1].(*TappableLabel)
			nameEntry.Segments[0].(*widget.TextSegment).Style.ColorName = nameEntry.Color

			if ings.Ingredients[iter].Highlighted {
				nameEntry.Segments[0].(*widget.TextSegment).Style.ColorName = theme.ColorNameError

			}

		}
		items.(*ingredients).Update = true

	}
	if recs, ok := items.(*recipes); ok {
		for iter := range list {
			nameEntry := list[iter].(*fyne.Container).Objects[1].(*TappableLabel)
			nameEntry.Segments[0].(*widget.TextSegment).Style.ColorName = nameEntry.Color

			if recs.Recipes[iter].Highlighted {
				nameEntry.Segments[0].(*widget.TextSegment).Style.ColorName = theme.ColorNameError

			}
			//nameEntry.Refresh()

		}
		items.(*recipes).Update = true

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
	//Win      fyne.Window
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
	//label.Win = fyne.CurrentApp().NewWindow(fmt.Sprintf("%v %v window", label.Segments[0].(*widget.TextSegment).Text, label.EntryInd))
	return label
}

func (t *TappableLabel) SetText(s string) {
	t.Segments[0].(*widget.TextSegment).Text = s
}

func (t *TappableLabel) GetText() string {
	return t.Segments[0].(*widget.TextSegment).Text
}

// ----------- entries --------------

func MakeRecEntries(recs *recipes, ings *ingredients, a fyne.App) []fyne.CanvasObject {
	cnvs := make([]fyne.CanvasObject, 0)
	for ind, r := range recs.Recipes {
		nameEntry := NewTabableLabel(r.Name, ind)

		checkBox := widget.NewCheck("", func(b bool) {

			if b {

				nameEntry.Color = theme.ColorNameDisabled
				recs.Recipes[ind].Check = true

			} else {
				nameEntry.Color = theme.ColorNameForeground
				recs.Recipes[ind].Check = false

			}

			nameEntry.Segments[0].(*widget.TextSegment).Style.ColorName = nameEntry.Color
			recs.Update = true

		})

		ingWindowBtn := widget.NewButtonWithIcon("", theme.MenuIcon(), func() { buildIngredientsWindow(a, &recs.Recipes[ind], ings) })

		nameEntry.CallBack = func() bool {

			recs.Remove(nameEntry.EntryInd)
			return true
		}
		checkBox.SetChecked(r.Check)
		cont := container.New(
			layout.NewCustomPaddedHBoxLayout(3),
			checkBox,
			nameEntry,
			ingWindowBtn,
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
		unitEntry := NewSearchEntry()
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

func MakeEmptyEntry(g Groceitem) fyne.CanvasObject {
	var cont *fyne.Container

	if _, ok := g.(*ingredients); ok {
		nameEntry := NewTabableLabel("", 0)
		checkBox := widget.NewCheck("", func(b bool) {})
		unitEntry := NewSearchEntry()

		unitSel := widget.NewSelect(unitVals, func(v string) {})
		unitSel.SetSelectedIndex(0)

		cont := container.New(
			layout.NewCustomPaddedHBoxLayout(3),
			checkBox,
			nameEntry,
			unitEntry,
			unitSel,
		)
		return cont
	}

	if _, ok := g.(*recipes); ok {

		nameEntry := NewTabableLabel("", 0)

		checkBox := widget.NewCheck("", func(b bool) {})

		cont = container.New(
			layout.NewCustomPaddedHBoxLayout(3),
			nameEntry,
			checkBox,
		)
	}

	if cont == nil {
		cont = container.NewVBox()
	}

	return cont

}

func UpdateEntry(i widget.ListItemID, o fyne.CanvasObject, g Groceitem, f []fyne.CanvasObject) {

	o.(*fyne.Container).Objects = f[i].(*fyne.Container).Objects

}

func DrawEntries(e []fyne.CanvasObject, c *fyne.Container, g Groceitem) {
	c.RemoveAll()

	list := widget.NewList(
		func() int { return len(e) },
		func() fyne.CanvasObject { return MakeEmptyEntry(g) },
		func(i widget.ListItemID, o fyne.CanvasObject) {
			UpdateEntry(i, o, g, e)
		})

	c.Add(list)

}

func AddEntry(g Groceitem, c *fyne.Container, w fyne.Window) bool {

	textInput := NewSearchEntry()

	textInput.SetPlaceHolder("eh?")

	if i, ok := g.(*ingredients); ok {

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

					i.Add(formItem.Widget.(*SearchEntry).Text)

				}
			},
			w,
		)

		dialg.Show()

	}

	if r, ok := g.(*recipes); ok {

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
					r.Add(formItem.Widget.(*SearchEntry).Text)

				}
			},
			w,
		)

		dialg.Show()

	}
	return true
}

func UpdateEntries(g Groceitem, o Groceitem, c *fyne.Container, e *[]fyne.CanvasObject, t *Trie, a fyne.App) {
	if i, ok := g.(*ingredients); ok {
		for {
			if len(i.Ingredients) != len(*e) || i.Update {
				t.build(i)
				i.CheckSort()
				i.HighlightSort()
				*e = MakeIngEntries(i)
				DrawEntries(*e, c, i)
				DrawHighlights(i, e)
				c.Refresh()
				i.Update = false
			}
		}
	}

	if r, ok := g.(*recipes); ok {
		for {
			if len(r.Recipes) != len(*e) || r.Update {
				//fmt.Println(r)
				t.build(r)
				r.CheckSort()
				r.HighlightSort()
				*e = MakeRecEntries(r, o.(*ingredients), a)
				DrawEntries(*e, c, r)
				DrawHighlights(r, e)
				c.Refresh()
				r.Update = false
			}
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

func GetEntryData(c []fyne.CanvasObject, g Groceitem) {
	if i, ok := g.(*ingredients); ok {
		for ind, con := range c {

			rCon := con.(*fyne.Container)

			i.Ingredients[ind].Name = rCon.Objects[1].(*TappableLabel).GetText()
			n, err := strconv.ParseFloat(rCon.Objects[2].(*SearchEntry).Text, 64)

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

	if r, ok := g.(*recipes); ok {
		for ind, con := range c {

			rCon := con.(*fyne.Container)

			r.Recipes[ind].Name = rCon.Objects[1].(*TappableLabel).GetText()

			r.Recipes[ind].Check = rCon.Objects[0].(*widget.Check).Checked

		}

	}
}

// ----------- Main --------------

func buildIngredientsWindow(a fyne.App, r *recipe, i *ingredients) fyne.Window {

	//fmt.Println(r)
	ingSearch := Trie{}
	ingSearch.build(&r.RecipeIngs)

	window := a.NewWindow(fmt.Sprintf("%v Ingrediants", r.Name))
	window.Resize(fyne.NewSize(WINSIZEX, WINSIZEY))

	ingContainer := container.NewStack()
	ingsEntries := MakeIngEntries(&r.RecipeIngs)
	DrawEntries(ingsEntries, ingContainer, &r.RecipeIngs)
	go UpdateEntries(&r.RecipeIngs, i, ingContainer, &ingsEntries, &ingSearch, a)

	addIngsBtn := widget.NewToolbarAction(theme.ContentAddIcon(), func() {

		AddEntry(&r.RecipeIngs, ingContainer, window)

	})

	transIngsBtn := widget.NewToolbarAction(theme.ContentUndoIcon(), func() {

		GetEntryData(ingsEntries, &r.RecipeIngs)
		i.TransferIngredients(&r.RecipeIngs)
		msg := dialog.NewInformation(
			fmt.Sprintf("%v ingredients transfer", r.Name),
			fmt.Sprintf("%v ingredients added to main list", r.Name),
			window)
		msg.Show()

	})

	ingToolbar := widget.NewToolbar(addIngsBtn, transIngsBtn)
	ingTopCont := container.NewVBox(ingToolbar)

	//ingToolbar.Move(fyne.NewPos(0, 30))
	ingContainer.Resize(fyne.NewSize(WINSIZEX-10, WINSIZEY-180))
	ingContainer.Move(fyne.NewPos(0, 100))

	ingMainCont := container.NewWithoutLayout(ingTopCont, ingContainer)

	window.SetContent(ingMainCont)
	window.Show()
	window.SetOnClosed(func() {
		GetEntryData(ingsEntries, &r.RecipeIngs)

	})
	window.RequestFocus()
	window.CenterOnScreen()

	return window
}

func BuildUI(a fyne.App, w fyne.Window, i *ingredients, r *recipes, d fyne.URI) {
	// ingredients
	ingSearch := Trie{}
	recSearch := Trie{}

	ingSearch.build(i)
	recSearch.build(r)

	//fmt.Println(recSearch)
	//listObj := &widget.List{}

	ingContainer := container.NewStack()
	//ingContainer.Resize(fyne.NewSize(WINSIZEX, WINSIZEY))

	ingsEntries := MakeIngEntries(i)
	DrawEntries(ingsEntries, ingContainer, i)
	ingContainer.Resize(fyne.NewSize(WINSIZEX, WINSIZEY))
	go UpdateEntries(i, r, ingContainer, &ingsEntries, &ingSearch, a)

	addIngsBtn := widget.NewToolbarAction(theme.ContentAddIcon(), func() { AddEntry(i, ingContainer, w) })
	ingSaveBtn := widget.NewToolbarAction(theme.DocumentSaveIcon(), func() {
		GetEntryData(ingsEntries, i)
		SaveData(d, i, r)
	})
	chekAllBtn := widget.NewToolbarAction(theme.CheckButtonCheckedIcon(), func() {
		GetEntryData(ingsEntries, i)
		i.CheckAll()
		i.Update = true
	})

	ingToolbar := widget.NewToolbar(addIngsBtn, ingSaveBtn, chekAllBtn)
	ingTopCont := container.NewVBox(ingToolbar)
	ingSearchBar := NewSearchEntry()

	ingSearchBar.OnChanged = func(string) {
		ingSearchBar.HighlightSearch(&ingsEntries, i, &ingSearch)
	}

	ingMainCont := container.NewWithoutLayout(ingTopCont, ingSearchBar, ingContainer)
	//ingToolbar.Move(fyne.NewPos(0, 30))
	ingSearchBar.Resize(fyne.NewSize(WINSIZEX-10, 50))
	ingSearchBar.Resize(fyne.NewSize(WINSIZEX-10, 50))
	ingSearchBar.Move(fyne.NewPos(0, 40))
	ingContainer.Resize(fyne.NewSize(WINSIZEX-10, WINSIZEY-180))
	ingContainer.Move(fyne.NewPos(0, 100))

	// recipes
	recContainer := container.NewStack()
	recEntries := MakeRecEntries(r, i, a)
	DrawEntries(recEntries, recContainer, r)
	go UpdateEntries(r, i, recContainer, &recEntries, &recSearch, a)

	addRecBtn := widget.NewToolbarAction(theme.ContentAddIcon(), func() { AddEntry(r, recContainer, w) })

	saveBtn := widget.NewToolbarAction(theme.DocumentSaveIcon(), func() {
		GetEntryData(recEntries, r)
		SaveData(d, i, r)
	})
	recToolbar := widget.NewToolbar(addRecBtn, saveBtn)

	recTopCont := container.NewVBox(
		recToolbar,
	)
	recSearchBar := NewSearchEntry()
	recSearchBar.OnChanged = func(string) {
		recSearchBar.HighlightSearch(&recEntries, r, &recSearch)
	}

	recMainCont := container.NewWithoutLayout(recTopCont, recSearchBar, recContainer)
	//ingToolbar.Move(fyne.NewPos(0, 30))

	recSearchBar.Resize(fyne.NewSize(WINSIZEX-10, 50))
	recSearchBar.Resize(fyne.NewSize(WINSIZEX-10, 50))
	recSearchBar.Move(fyne.NewPos(0, 40))
	recContainer.Resize(fyne.NewSize(WINSIZEX-10, WINSIZEY-180))
	recContainer.Move(fyne.NewPos(0, 100))

	Tab := container.NewAppTabs(
		container.NewTabItem("Ingredients", ingMainCont),
		container.NewTabItem("Reipes", recMainCont))

	Tab.OnSelected = func(t *container.TabItem) {
		i.Update = true
		r.Update = true

	}

	w.SetContent(Tab)

	a.Settings().SetTheme(&CustomTheme{})

}
