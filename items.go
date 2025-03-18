package main

import (
	"encoding/json"
	"fmt"
	"os"
	"unicode"
)

func sanatize_string(s string) string {
	var newString string
	for _, chr := range s {

		if unicode.IsLetter(chr) || unicode.IsDigit(chr) || unicode.IsSpace(chr) {
			newString += string(chr)
		}

	}
	return newString
}

type Groceitem interface {
	Add(name string)
	Remove(ind int)
	Read(data []byte)
}

type ingredient struct {
	Amount      float64
	Name        string
	Unit        unitType
	Check       bool
	Highlighted bool
}

type ingredients struct {
	Ingredients []ingredient
	Update      bool
}

type recipe struct {
	RecipeIngs ingredients
	Name       string
	Check      bool
}

type recipes struct {
	Recipes []recipe
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

func (i *ingredients) Remove(ind int) {

	newIngs := ingredients{}
	//fmt.Println(ind)

	if ind > 0 {
		newIngs.Ingredients = append(newIngs.Ingredients, i.Ingredients[0:ind]...)
	}

	newIngs.Ingredients = append(newIngs.Ingredients, i.Ingredients[ind+1:int(len(i.Ingredients))]...)

	newIngs.Ingredients, i.Ingredients = i.Ingredients, newIngs.Ingredients

	//fmt.Println(i.Ingredients)

}

func (r *recipes) Remove(ind int) {
	newRec := recipes{}

	if ind > 0 {
		newRec.Recipes = append(newRec.Recipes, r.Recipes[0:ind]...)

	}
	newRec.Recipes = append(newRec.Recipes, r.Recipes[ind+1:int(len(r.Recipes))]...)

	newRec.Recipes, r.Recipes = r.Recipes, newRec.Recipes
}

func (i *ingredients) Add(name string) {
	i.Ingredients = append(i.Ingredients, ingredient{
		Name:        name,
		Amount:      1.0,
		Unit:        unt,
		Check:       false,
		Highlighted: false,
	})
}

func (r *recipes) Add(name string) {
	r.Recipes = append(r.Recipes, recipe{
		Name:       name,
		RecipeIngs: ingredients{},
		Check:      false,
	})
}

func (i *ingredients) Read(data []byte) {

	i.Ingredients = []ingredient{}

	jsnMap := make(map[string][]map[string]interface{})
	json.Unmarshal(data, &jsnMap)

	for ind, ing := range jsnMap["i"] {

		i.Add(ing["Name"].(string))
		i.Ingredients[ind].Amount = ing["Amount"].(float64)
		i.Ingredients[ind].Unit = ing["Unit"].(unitType)
		i.Ingredients[ind].Check = ing["Check"].(bool)
	}

}

func (r *recipes) Read(data []byte) {
	r.Recipes = []recipe{}

	jsnMap := make(map[string][]map[string]interface{})
	json.Unmarshal(data, &jsnMap)

	for ind, rec := range jsnMap["r"] {

		r.Add(rec["Name"].(string))
		r.Recipes[ind].Check = rec["Check"].(bool)
		r.Recipes[ind].RecipeIngs = ingredients{}

		if rec["Ingredients"] != nil {
			for ingInd, ing := range rec["Ingredients"].([]interface{}) {
				ingMap := ing.(map[string]interface{})

				r.Recipes[ind].RecipeIngs.Add(ingMap["Name"].(string))
				r.Recipes[ind].RecipeIngs.Ingredients[ingInd].Amount = ingMap["Amount"].(float64)
				r.Recipes[ind].RecipeIngs.Ingredients[ingInd].Unit = ingMap["Unit"].(unitType)
				r.Recipes[ind].RecipeIngs.Ingredients[ingInd].Check = ingMap["Check"].(bool)

			}
		}
	}
}

func (i *ingredients) CheckSort() {
	checked := make([]ingredient, 0)
	nonChecked := make([]ingredient, 0)

	for _, item := range i.Ingredients {
		if !item.Check {
			nonChecked = append(nonChecked, item)
		} else {
			checked = append(checked, item)
		}
	}

	nonChecked = append(nonChecked, checked...)
	i.Ingredients = nonChecked

}

func (i *ingredients) HighlightSort() {

	highlight := make([]ingredient, 0)
	nonHighlight := make([]ingredient, 0)

	for _, item := range i.Ingredients {
		if !item.Highlighted {
			nonHighlight = append(nonHighlight, item)
		} else {
			highlight = append(highlight, item)
		}
	}

	highlight = append(highlight, nonHighlight...)
	i.Ingredients = highlight

}

func CreateJson(i ingredients, r recipes) []byte {

	data := make(map[string]interface{}, 0)

	data["i"] = i.Ingredients
	data["r"] = r.Recipes

	dataJson, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}

	return dataJson

}

func WriteFile(file string, i ingredients, r recipes) {

	data := make(map[string]interface{}, 0)

	data["i"] = i.Ingredients
	data["r"] = r.Recipes

	dataJson, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}
	os.WriteFile(file, dataJson, 0644)

}
