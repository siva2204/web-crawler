package trie

import (
	"fmt"
	"unicode/utf8"
)

// Node represents each node in Trie.
type Node struct {
	children map[rune]*Node // map children nodes
	urls     []string       // current node value
	isleaf   bool           // is current node a leaf node
}

// NewNode creates a new Trie node with initialized
// children map.
func NewNode() *Node {
	n := &Node{}
	n.children = make(map[rune]*Node)
	n.urls = []string{}
	n.isleaf = false
	return n
}

// insert a single word at a Trie node.
func (n *Node) insert(s string, url string) {
	fmt.Println("inserting", s)
	curr := n
	for _, c := range s {
		next, ok := curr.children[c]
		if !ok {
			next = NewNode()
			curr.children[c] = next
		}
		curr = next
	}
	curr.urls = append(curr.urls, url)
	curr.isleaf = true
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
	return next.urls
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
	if n.isleaf {
		r++
	}
	return r
}

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
		if node.isleaf {
			*texts = append(*texts, newPrefix)
		}
		SearchSubTree(node, texts, newPrefix)
	}
}

func (root *Node) AutoCompletePrefix(prefix string) []string {
	subTrie := FindSubTrieByPrefix(root, prefix)
	var prefixes []string
	SearchSubTree(subTrie, &prefixes, prefix)
	return prefixes
}
