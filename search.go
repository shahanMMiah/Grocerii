package main

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

	s = sanatize_string(s)

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

	s = sanatize_string(s)

	tempLevel := &t.RootNode
	for ind, chr := range s {
		level, exist := tempLevel.Children[chr]

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

	s = sanatize_string(s)
	trieLevel, found := t.Find(s)

	if !found {
		return nil
	}

	return FindWords(trieLevel, []string{}, s[:len(s)-1])

}

func FindWords(tn *TrieNode, s []string, cs string) []string {

	cs += string(tn.char)

	lWords := make([]string, 0)
	for _, node := range tn.Children {

		lWords = append(lWords, FindWords(&node, s, cs)...)

	}
	s = append(s, lWords...)

	if tn.End {
		s = append(s, cs)
	}

	return s
}

func (t *Trie) build(g Groceitem) {

	*t = Trie{}

	ings, ok := g.(*ingredients)
	if ok {
		for _, i := range ings.Ingredients {
			t.Add(i.Name)
		}
		return
	}

	recs, ok := g.(*recipes)
	if ok {
		for _, r := range recs.Recipes {
			t.Add(r.Name)
		}
		return
	}

}
