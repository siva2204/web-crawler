package trie

import (
	"unicode/utf8"
)

// Node represents each node in Trie.
type Node struct {
	children map[rune]*Node // map children nodes
	isLeaf   []string       // current node value
}

// NewNode creates a new Trie node with initialized
// children map.
func NewNode() *Node {
	n := &Node{}
	n.children = make(map[rune]*Node)
	n.isLeaf = []string{}
	return n
}

// insert a single word at a Trie node.
func (n *Node) insert(s string, url string) {
	curr := n
	for _, c := range s {
		next, ok := curr.children[c]
		if !ok {
			next = NewNode()
			curr.children[c] = next
		}
		curr = next
	}
	curr.isLeaf = append(curr.isLeaf, url)
}

// Insert zero, one or more words at a Trie node.
func (n *Node) Insert(s []string, url string) {
	for _, ss := range s {
		n.insert(ss, url)
	}
}

// Find  words at a Trie node.
func (n *Node) Find(s string) []string {
	next, ok := n, false
	for _, c := range s {
		next, ok = next.children[c]
		if !ok {
			return []string{}
		}
	}
	return next.isLeaf
}

// Capacity returns the number of nodes in the Trie
func (n *Node) Capacity() int {
	r := 0
	for _, c := range n.children {
		r += c.Capacity()
	}
	return 1 + r
}

// Size returns the number of words in the Trie
func (n *Node) Size() int {
	r := 0
	for _, c := range n.children {
		r += c.Size()
	}
	if n.isLeaf != nil {
		r++
	}
	return r
}

// remove lazily a word from the Trie node, no node is actually removed.
// func (n *Node) remove(s string) {
// 	if len(s) == 0 {
// 		return
// 	}
// 	next, ok := n, false
// 	for _, c := range s {
// 		next, ok = next.children[c]
// 		if !ok {
// 			// word cannot be found - we're done !
// 			return
// 		}
// 	}
// 	next.isLeaf =
// }

// // Remove zero, one or more words lazily from the Trie, no node is actually removed.
// func (n *Node) Remove(s ...string) {
// 	for _, ss := range s {
// 		n.remove(ss)
// 	}
// }

// // Compact will remove unecessay nodes, reducing the capacity, returning true if node n itself should be removed.
// func (n *Node) Compact() (remove bool) {

// 	for r, c := range n.children {
// 		if c.Compact() {
// 			delete(n.children, r)
// 		}
// 	}
// 	return !n.isLeaf && len(n.children) == 0
// }

func FindSubTrieByPrefix(root *Node, prefix string) *Node {
	if prefix == "" {
		return root
	}
	if root == nil {
		return nil
	}
	key := []rune(prefix)[0]
	if _, exists := root.children[key]; !exists {
		return root
	}
	_, i := utf8.DecodeRuneInString(prefix)
	return FindSubTrieByPrefix(root.children[key], prefix[i:])
}

func SearchSubTree(root *Node, texts *[]string, prefix string) {
	if root == nil {
		return
	}
	for k, node := range root.children {
		var newPrefix string = prefix + string(k)
		if node.isLeaf != nil {
			*texts = append(*texts, newPrefix)
		}
		SearchSubTree(node, texts, newPrefix)
	}
}

func AutoCompletePrefix(root *Node, prefix string) []string {
	subTrie := FindSubTrieByPrefix(root, prefix)
	var prefixes []string
	SearchSubTree(subTrie, &prefixes, prefix)
	return prefixes
}
