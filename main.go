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
	for {
		fmt.Println("which command -- add, view")

		reader := bufio.NewReader(os.Stdin)
		text, _ := reader.ReadString('\n')
		text = strings.Replace(text, "\n", "", -1)
		fmt.Println(text)

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
			ings = AddIngredients(name, amount, e, ings)

		case "view":

			fmt.Println("ingredient list:")
			for _, i := range ings {
				fmt.Printf("%v %c - amount: %v \n", i.name, i.emoji, i.amount)

			}

		default:
			return
		}

	}
}
