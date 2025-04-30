// fyne package -os android -appID .com.test.groceriiTest
//
//go:generate fyne bundle -o bundled.go assets
package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
)

const (
	WINSIZEX float32 = 393
	WINSIZEY float32 = 750
)

func main() {

	a := app.NewWithID("Grocerii")

	w := a.NewWindow("Grocerii App")
	w.Resize(fyne.NewSize(WINSIZEX, WINSIZEY))
	w.CenterOnScreen()
	w.SetMaster()
	dataUR := GetDataURI(a)

	ings := ingredients{Update: false}
	recipes := recipes{}

	ings.Read(ReadData(dataUR))
	recipes.Read(ReadData(dataUR))

	a.Settings().SetTheme(&CustomTheme{})

	BuildUI(a, w, &ings, &recipes, dataUR)

	w.Show()

	a.Run()

}
