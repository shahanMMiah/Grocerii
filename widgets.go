package main

import (
	"image/color"

	"fyne.io/fyne/v2"
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

	switch name {
	case theme.ColorNameBackground:
		return CustomColor{r: 240, g: 240, b: 191, a: 255}

	case theme.ColorNameForeground:
		return CustomColor{r: 0, g: 0, b: 0, a: 255}

	case theme.ColorNameInputBackground,
		theme.ColorNameMenuBackground,
		theme.ColorNameOverlayBackground:
		return CustomColor{r: 245, g: 245, b: 210, a: 255}

	case theme.ColorNameButton:
		return CustomColor{r: 245, g: 245, b: 210, a: 255}

	}

	return theme.DefaultTheme().Color(name, variant)
}
func (m CustomTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(name)
}

func (m CustomTheme) Font(style fyne.TextStyle) fyne.Resource {
	return resourceCheescakeMonolineTtf

}

func (m CustomTheme) Size(name fyne.ThemeSizeName) float32 {
	switch name {
	case theme.SizeNameText:
		return 28
	}
	return theme.DefaultTheme().Size(name)
}

// ----------- Search Bar --------------
type SearchEntry struct {
	widget.Entry
	Update bool
}

func NewSearchEntry() *SearchEntry {
	entry := &SearchEntry{}

	entry.ExtendBaseWidget(entry)

	return entry

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

	return label
}

func (t *TappableLabel) SetText(s string) {
	t.Segments[0].(*widget.TextSegment).Text = s
}

func (t *TappableLabel) GetText() string {
	return t.Segments[0].(*widget.TextSegment).Text
}
