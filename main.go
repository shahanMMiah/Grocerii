package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	ings := MakeIngredients()

	recipes := make([]recipe, 0)

	rIngs := MakeIngredients()
	rIngs = AddIngredients("cofee", 200.0, rune('c'), grams, rIngs)

	recipes = append(recipes, MakeRecipe("test", rIngs))

	for {
		fmt.Println("which command? -- add, view, save, read ")

		reader := bufio.NewReader(os.Stdin)
		text, _ := reader.ReadString('\n')
		text = strings.Replace(text, "\n", "", -1)

		switch text {
		case "add":

			fmt.Println("name?")

			reader = bufio.NewReader(os.Stdin)
			name, _ := reader.ReadString('\n')
			name = strings.Replace(name, "\n", "", -1)

			fmt.Println("amount?")

			reader = bufio.NewReader(os.Stdin)
			amountStr, _ := reader.ReadString('\n')
			amountStr = strings.Replace(amountStr, "\n", "", -1)
			amount, _ := strconv.ParseFloat(amountStr, 64)

			e := '♡'
			ings = AddIngredients(name, amount, e, unt, ings)

		case "view":

			fmt.Println("ingredient list:")
			for _, i := range ings {
				fmt.Printf("%v %c - amount: %v %v \n", i.Name, i.Emoji, i.Unit, i.Amount)

			}

		case "save":
			WriteFile("test.json", ings, recipes)

		case "read":
			fmt.Println("file name?")
			reader = bufio.NewReader(os.Stdin)
			file, _ := reader.ReadString('\n')
			file = strings.Replace(file, "\n", "", -1)

			ings, recipes = ReadFile(file)
			fmt.Printf("read in %v and %v", ings, recipes)

		default:
			return
		}

	}
}
