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

type GrocBaseItem struct {
	Name  string
	Check bool

	Highlighted bool
}

type ingredient struct {
	GrocBaseItem
	Amount     float64
	Unit       unitType
	Referenced []*ingredient
}

type ingredients struct {
	Ingredients []*ingredient
	Update      bool
}

type recipe struct {
	GrocBaseItem
	RecipeIngs ingredients
}

type recipes struct {
	Recipes []recipe
	Update  bool
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

func (i *ingredient) RemoveReferenced(ind int) {
	newRef := make([]*ingredient, 0)

	if ind > 0 {
		newRef = append(newRef, i.Referenced[0:ind]...)
	}

	newRef = append(newRef, i.Referenced[ind+1:int(len(i.Referenced))]...)

	i.Referenced = newRef

}

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
	i.Ingredients = append(i.Ingredients, &ingredient{
		GrocBaseItem: GrocBaseItem{
			Name:        name,
			Check:       false,
			Highlighted: false,
		},
		Unit:   unt,
		Amount: 1.0,
	})
}

func (r *recipes) Add(name string) {
	r.Recipes = append(r.Recipes, recipe{
		GrocBaseItem: GrocBaseItem{
			Name:        name,
			Check:       false,
			Highlighted: false,
		},
		RecipeIngs: ingredients{},
	})
}

func (i *ingredients) Read(data []byte) {

	i.Ingredients = []*ingredient{}

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

		if rec["RecipeIngs"] == nil {
			continue
		}

		recIngs := rec["RecipeIngs"].(map[string]interface{})

		if recIngs["Ingredients"] == nil {
			continue
		}

		for ingInd, ing := range recIngs["Ingredients"].([]interface{}) {

			//fmt.Println(ing)
			ingMap := ing.(map[string]interface{})

			r.Recipes[ind].RecipeIngs.Add(ingMap["Name"].(string))
			r.Recipes[ind].RecipeIngs.Ingredients[ingInd].Amount = ingMap["Amount"].(float64)
			r.Recipes[ind].RecipeIngs.Ingredients[ingInd].Unit = ingMap["Unit"].(unitType)
			r.Recipes[ind].RecipeIngs.Ingredients[ingInd].Check = ingMap["Check"].(bool)

		}

	}
}

func (i *ingredients) CheckAll(check *bool) {

	for ind := range i.Ingredients {
		i.Ingredients[ind].Check = *check

	}
}

func (i *ingredient) AddReference(ing *ingredient) {
	i.Referenced = append(i.Referenced, ing)
}

func (i *ingredient) CheckReferenced() {
	//fmt.Printf("reeferenced %v", i.Referenced)
	for iter := range i.Referenced {

		i.Referenced[iter].Check = true
		i.RemoveReferenced(iter)
	}
}

func (ings *ingredients) TransferIngredients(i *ingredients) {

	for tInd, tIng := range i.Ingredients {
		found := false
		for ind, ing := range ings.Ingredients {
			if sanatize_string(tIng.Name) == sanatize_string(ing.Name) && tIng.Unit == ing.Unit {
				if ing.Check {
					ings.Ingredients[ind].Check = false
					ings.Ingredients[ind].Amount = tIng.Amount

				} else {
					ings.Ingredients[ind].Amount += tIng.Amount
				}
				ings.Ingredients[ind].AddReference(i.Ingredients[tInd])
				//fmt.Printf("%v has these referenced %v", ing.Name, ings.Ingredients[ind].Referenced)
				found = true
				break
			}
		}

		if !found {
			ings.Add(tIng.Name)
			ings.Ingredients[len(ings.Ingredients)-1].Unit = tIng.Unit
			ings.Ingredients[len(ings.Ingredients)-1].Amount = tIng.Amount
			ings.Ingredients[len(ings.Ingredients)-1].Check = tIng.Check
		}
	}
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
