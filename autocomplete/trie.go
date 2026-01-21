package autocomplete

type node struct {
	terminal bool
	word     string
	children map[byte]*node
}

type Trie struct {
	root *node
}

func NewTrie(words ...string) *Trie {
	trie := Trie{
		root: &node{children: make(map[byte]*node)},
	}
	for _, word := range words {
		trie.Insert(word)
	}
	return &trie
}

func (t *Trie) Insert(word string) {
	curr := t.root
	for i := range len(word) {
		c := word[i]
		if _, ok := curr.children[c]; !ok {
			curr.children[c] = &node{children: make(map[byte]*node)}
		}
		curr = curr.children[c]
	}
	curr.terminal = true
	curr.word = word
}

func (t *Trie) Search(word string) bool {
	node := t.findNode(word)
	return node != nil && node.terminal
}

func (t *Trie) StartsWith(prefix string) bool {
	node := t.findNode(prefix)
	return node != nil
}

func (t *Trie) FindAllWithPrefix(prefix string) []string {
	start := t.findNode(prefix)
	if start == nil {
		return nil
	}

	var results []string
	var dfs func(node *node)

	dfs = func(node *node) {
		if node.terminal {
			results = append(results, node.word)
		}
		for _, child := range node.children {
			dfs(child)
		}
	}

	dfs(start)
	return results
}

func (t *Trie) findNode(prefix string) *node {
	curr := t.root
	for i := range len(prefix) {
		next, ok := curr.children[prefix[i]]
		if !ok {
			return nil
		}
		curr = next
	}
	return curr
}
