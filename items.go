package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type Trie struct {
	RootNode TrieNode
}

type TrieNode struct {
	char     rune
	Children map[rune]TrieNode
	End      bool
}

func NewTrieNode() TrieNode {
	return TrieNode{
		Children: make(map[rune]TrieNode),
		char:     '0',
		End:      false,
	}
}

func (t *Trie) Add(s string) {

	if t.RootNode.Children == nil {
		t.RootNode = NewTrieNode()
	}

	tempLevel := &t.RootNode

	for ind, chr := range s {
		level, exist := tempLevel.Children[chr]
		if !exist {
			level = NewTrieNode()
			level.char = chr
		}

		if ind == len(s)-1 {

			level.End = true
		}

		tempLevel.Children[rune(chr)] = level
		tempLevel = &level

	}
}

func (t *Trie) Find(s string) (*TrieNode, bool) {

	tempLevel := &t.RootNode
	for ind, chr := range s {
		level, exist := tempLevel.Children[chr]
		fmt.Println(tempLevel.Children[chr])

		if !exist && !level.End {
			return tempLevel, false
		}

		tempLevel = &level

		if ind == len(s)-1 {

			return tempLevel, true
		}

	}

	return tempLevel, false

}

func (t *Trie) AutoComplete(s string) []string {
	/*
		check if string exisit and get current level
		from current level
		call get complete that returns list of existing words possible from level
	*/

	trieLevel, found := t.Find(s)

	if !found {
		return nil
	}

	return FindWords(trieLevel, []string{}, s[:len(s)-1])

}

func FindWords(tn *TrieNode, s []string, cs string) []string {
	/*
		from currnet level check children and and call get complete to get possible from level
			if level is end add current cs + letter to return slice slice
			for each of children letter, concat retured list from called get complete
	*/

	cs += string(tn.char)

	lWords := make([]string, 0)
	for _, node := range tn.Children {

		lWords = append(lWords, FindWords(&node, s, cs)...)

	}
	s = append(s, lWords...)

	if tn.End {
		s = append(s, cs)
	}

	//fmt.Printf("at %v list is at %v \n", cs, s)
	return s
}

type Groceitem interface {
	Add(name string)
	Remove(ind int)
	Read(data []byte)
}

type ingredient struct {
	Amount float64
	Name   string
	Unit   unitType
	Check  bool
}

type ingredients struct {
	Ingredients []ingredient
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
	fmt.Println(ind)

	if ind > 0 {
		newIngs.Ingredients = append(newIngs.Ingredients, i.Ingredients[0:ind]...)
	}

	newIngs.Ingredients = append(newIngs.Ingredients, i.Ingredients[ind+1:int(len(i.Ingredients))]...)

	newIngs.Ingredients, i.Ingredients = i.Ingredients, newIngs.Ingredients

	fmt.Println(i.Ingredients)

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
		Name:   name,
		Amount: 1.0,
		Unit:   unt,
		Check:  false,
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
