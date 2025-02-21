package main

import (
	"encoding/json"
	"fmt"
	"os"

	"fyne.io/fyne/v2"
)

type item interface{}

type recipe struct {
	Name        string
	Ingredients []ingredient
	Check       bool
}

type ingredient struct {
	Name   string
	Amount float64
	Emoji  rune
	Unit   unitType
	Check  bool
}

type unitType = string

const (
	grams unitType = "g"
	kg    unitType = "kg"
	ml    unitType = "ml"
	ltr   unitType = "ltr"
	unt   unitType = "unit"
)

func getUnitInd(s string) int {
	unitMap := map[string]int{
		"g":    0,
		"kg":   1,
		"ml":   2,
		"ltr":  3,
		"unit": 4,
	}

	if _, exist := unitMap[s]; !exist {
		panic(fmt.Errorf("unit %v does not exists", s))
	}

	return unitMap[s]
}

var unitVals = []string{grams, kg, ml, ltr, unt}

func MakeIngredients() []ingredient {
	ingredients := make([]ingredient, 0)
	return ingredients
}

func RemoveIngredient(i []ingredient, ind int) []ingredient {

	newIngs := make([]ingredient, 0)
	for nInd, nIng := range i {
		if nInd != ind {
			newIngs = append(newIngs, nIng)
		}
	}

	return newIngs

}

func RemoveRecipes(i []recipe, ind int) []recipe {

	newRec := make([]recipe, 0)
	for nInd, nIng := range i {
		if nInd != ind {
			newRec = append(newRec, nIng)
		}
	}

	return newRec

}

func AddIngredients(name string, amount float64, emoji rune, unit unitType, check bool, ingredients []ingredient) []ingredient {

	for i, val := range ingredients {
		if val.Name == name {
			ingredients[i].Amount += amount
			return ingredients
		}
	}

	i := ingredient{
		name,
		amount,
		emoji,
		unit,
		check,
	}
	ingredients = append(ingredients, i)
	return ingredients
}

func AddRecipe(name string, recipes []recipe) []recipe {

	i := recipe{
		Name:        name,
		Ingredients: nil,
	}

	recipes = append(recipes, i)
	return recipes
}

func MakeRecipe(name string, ings []ingredient, check bool) recipe {
	return recipe{
		name,
		ings,
		check,
	}
}

func AddRecipeIngredients(r recipe, ings []ingredient) []ingredient {
	for _, i := range r.Ingredients {
		ings = AddIngredients(i.Name, i.Amount, i.Emoji, i.Unit, i.Check, ings)
	}

	return ings
}

func CreateJson(i []ingredient, r []recipe) []byte {

	data := make(map[string]interface{}, 0)

	data["i"] = i
	data["r"] = r

	dataJson, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}

	return dataJson

}

func WriteFile(file string, i []ingredient, r []recipe) {

	data := make(map[string]interface{}, 0)

	data["i"] = i
	data["r"] = r

	dataJson, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}
	os.WriteFile(file, dataJson, 0644)

}

func ReadFile(ur fyne.URI) ([]ingredient, []recipe) {

	dataConent := ReadData(ur)
	data := make(map[string][]map[string]interface{})
	json.Unmarshal(dataConent, &data)

	ings := MakeIngredients()
	recs := make([]recipe, 0)

	for _, i := range data["i"] {

		ings = AddIngredients(
			i["Name"].(string),
			i["Amount"].(float64),
			rune(i["Emoji"].(float64)),
			i["Unit"].(unitType),
			i["Check"].(bool),
			ings)

	}

	for _, r := range data["r"] {
		rIng := MakeIngredients()

		if r["Ingredients"] != nil {

			for _, ing := range r["Ingredients"].([]interface{}) {
				ingMap := ing.(map[string]interface{})

				rIng = AddIngredients(
					ingMap["Name"].(string),
					ingMap["Amount"].(float64),
					rune(ingMap["Emoji"].(float64)),
					ingMap["Unit"].(unitType),
					ingMap["Check"].(bool),
					rIng,
				)

			}
		}
		recs = append(recs, MakeRecipe(r["Name"].(string), rIng, r["Check"].(bool)))

	}

	return ings, recs
}
