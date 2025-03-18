//go:generate fyne bundle -o bundled.go data.json

// fyne package -os android -appID .com.test.grocerii
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
	dataUR := GetDataURI(a)

	ings := ingredients{Update: false}
	recipes := recipes{}

	ings.Read(ReadData(dataUR))
	recipes.Read(ReadData(dataUR))

	//fmt.Println(ingSearch)

	BuildUI(a, w, &ings, &recipes, dataUR)

	w.ShowAndRun()

}

/*
--------------------------------- make default JSON  -------------------------------
itms := []string{
			"olives", "cherry toms", "orecchiette", "soft cheese", "spinach", "tofu", "red pepper",
			"noodles", "coconut cream", "coslaw mix", "cucumber", "carrot", "lettuce", "chikn", "bulger", "courgete",
			"onion", "capers", "aubergine", "cereal", "milk", "parmesan cheese"}

for _, itm := range itms {
	ings.Add(itm)
}

recs := []string{
	"sticky chilly tofu bowl bosh",
	"tomato olive orecchiette",
	"coconut curry noodles tofu",
	"bang bang chikn rice salad",
	"turkish style bulger",
	"aubergine caponato orzo",
}

for _, rcs := range recs {
	recipes.Add(rcs)
}

		WriteFile("data.json", ings, recipes)

--------------------------------- trie test  -------------------------------

test := Trie{}

test.Add("hello")
test.Add("hey")
test.Add("plum🍩 and bread")

words := test.AutoComplete("plum")

//fmt.Println(words)



--------------------------------- CLI -------------------------------
import (
	"fmt"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func main() {

	a := app.New()
	w := a.NewWindow("Hello")

	w.Resize(fyne.NewSize(500, 500))

	clock := widget.NewLabel("")

	go func() {
		for t := range time.Tick(time.Second) {

			h, m, s := t.Clock()
			clock.SetText(fmt.Sprintf("time : %v : %v : %v", h, m, s))
		}
	}()

	hello := widget.NewLabel("Hello Fyne!")
	w.SetContent(container.NewVBox(
		hello,
		clock,
		widget.NewButton("Hi!", func() {
			hello.SetText("MUSTARDDDDD")
		}),
	))

	w.ShowAndRun()

	/*

		import (
		"bufio"
		"fmt"
		"os"
		"strconv"
		"strings"
		)

		ings := MakeIngredients()

		recipes := make([]recipe, 0)

		rIngs := MakeIngredients()
		rIngs = AddIngredients("cofee", 200.0, rune('c'), grams, rIngs)

		recipes = append(recipes, MakeRecipe("test", rIngs))

		for {
			//fmt.Println("which command? -- add, view, save, read ")

			reader := bufio.NewReader(os.Stdin)
			text, _ := reader.ReadString('\n')
			text = strings.Replace(text, "\n", "", -1)

			switch text {
			case "add":

				//fmt.Println("name?")

				reader = bufio.NewReader(os.Stdin)
				name, _ := reader.ReadString('\n')
				name = strings.Replace(name, "\n", "", -1)

				//fmt.Println("amount?")

				reader = bufio.NewReader(os.Stdin)
				amountStr, _ := reader.ReadString('\n')
				amountStr = strings.Replace(amountStr, "\n", "", -1)
				amount, _ := strconv.ParseFloat(amountStr, 64)

				e := '♡'
				ings = AddIngredients(name, amount, e, unt, ings)

			case "view":

				//fmt.Println("ingredient list:")
				for _, i := range ings {
					//fmt.Printf("%v %c - amount: %v %v \n", i.Name, i.Emoji, i.Unit, i.Amount)

				}

			case "save":
				WriteFile("test.json", ings, recipes)

			case "read":
				//fmt.Println("file name?")
				reader = bufio.NewReader(os.Stdin)
				file, _ := reader.ReadString('\n')
				file = strings.Replace(file, "\n", "", -1)

				ings, recipes = ReadFile(file)
				//fmt.Printf("read in %v and %v", ings, recipes)

			default:
				return
			}

		}

}
*/
