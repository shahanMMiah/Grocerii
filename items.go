package main

type recipe struct {
	name        string
	ingredients []ingredient
}

type ingredient struct {
	name   string
	amount float64
	emoji  rune
}

func MakeIngredients() []ingredient {
	ingredients := make([]ingredient, 0)
	return ingredients
}

func AddIngredients(name string, amount float64, emoji rune, ingredients []ingredient) []ingredient {

	for i, val := range ingredients {
		if val.name == name {
			ingredients[i].amount += amount
			return ingredients
		}
	}

	i := ingredient{
		name,
		amount,
		emoji,
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
	for _, i := range r.ingredients {
		ings = AddIngredients(i.name, i.amount, i.emoji, ings)
	}

	return ings
}
