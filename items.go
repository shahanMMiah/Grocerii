package main

import (
	"encoding/json"
	"os"

	"fyne.io/fyne/v2"
)

type recipe struct {
	Name        string
	Ingredients []ingredient
}

type ingredient struct {
	Name   string
	Amount float64
	Emoji  rune
	Unit   unitType
}

type unitType = string

const (
	grams unitType = "g"
	kg    unitType = "kg"
	ml    unitType = "ml"
	ltr   unitType = "ltr"
	unt   unitType = "unit"
)

var unitVals = []string{grams, kg, ml, ltr, unt}

func MakeIngredients() []ingredient {
	ingredients := make([]ingredient, 0)
	return ingredients
}

func AddIngredients(name string, amount float64, emoji rune, unit unitType, ingredients []ingredient) []ingredient {

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
	}
	ingredients = append(ingredients, i)
	return ingredients
}

func MakeRecipe(name string, ings []ingredient) recipe {
	return recipe{
		name,
		ings,
	}
}

func AddRecipeIngredients(r recipe, ings []ingredient) []ingredient {
	for _, i := range r.Ingredients {
		ings = AddIngredients(i.Name, i.Amount, i.Emoji, i.Unit, ings)
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
			i["Unit"].(unitType), ings)

	}

	for _, r := range data["r"] {
		rIng := MakeIngredients()
		for _, ing := range r["Ingredients"].([]interface{}) {
			ingMap := ing.(map[string]interface{})

			rIng = AddIngredients(
				ingMap["Name"].(string),
				ingMap["Amount"].(float64),
				rune(ingMap["Emoji"].(float64)),
				ingMap["Unit"].(unitType), rIng)

		}

		recs = append(recs, MakeRecipe(r["Name"].(string), rIng))

	}

	return ings, recs
}
