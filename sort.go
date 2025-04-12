package main

func (i *ingredients) CheckSort() {
	checked := make([]*ingredient, 0)
	nonChecked := make([]*ingredient, 0)

	for iter, item := range i.Ingredients {
		if !item.Check {
			nonChecked = append(nonChecked, i.Ingredients[iter])
		} else {
			checked = append(checked, i.Ingredients[iter])
		}
	}

	nonChecked = append(nonChecked, checked...)
	i.Ingredients = nonChecked

}

func (r *recipes) CheckSort() {

	checked := make([]recipe, 0)
	nonChecked := make([]recipe, 0)

	for _, item := range r.Recipes {
		if !item.Check {
			nonChecked = append(nonChecked, item)
		} else {
			checked = append(checked, item)
		}
	}

	nonChecked = append(nonChecked, checked...)
	r.Recipes = nonChecked

}

func (i *ingredients) HighlightSort() {

	highlight := make([]*ingredient, 0)
	nonHighlight := make([]*ingredient, 0)

	for iter, item := range i.Ingredients {
		if !item.Highlighted {
			nonHighlight = append(nonHighlight, i.Ingredients[iter])
		} else {
			highlight = append(highlight, i.Ingredients[iter])
		}
	}

	highlight = append(highlight, nonHighlight...)
	i.Ingredients = highlight

}

func (r *recipes) HighlightSort() {

	highlight := make([]recipe, 0)
	nonHighlight := make([]recipe, 0)

	for _, item := range r.Recipes {
		//fmt.Println(item)
		if !item.Highlighted {
			nonHighlight = append(nonHighlight, item)
		} else {
			highlight = append(highlight, item)
		}
	}

	highlight = append(highlight, nonHighlight...)
	r.Recipes = highlight

}
