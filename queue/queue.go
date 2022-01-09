package queue

import "sync"

// Node will be store the value and the next node as well
type Node struct {
	url  string
	next *Node
}

// Queue structure is tell us what our head is and what tail should be with length of the list
type Queue struct {
	head   *Node
	tail   *Node
	length int
	sync.RWMutex
}

// enqueue it will be added new value into queue
func (ll *Queue) Enqueue(n string) {
	ll.Lock()
	defer ll.Unlock()

	var newNode Node // create new Node
	newNode.url = n  // set the data

	if ll.tail != nil {
		ll.tail.next = &newNode
	}

	ll.tail = &newNode

	if ll.head == nil {
		ll.head = &newNode
	}
	ll.length++
}

// dequeue it will be removed the first value into queue (First In First Out)
func (ll *Queue) Dequeue() string {
	ll.Lock()
	defer ll.Unlock()
	if ll.IsEmpty() {
		return "" // if is empty return -1
	}
	data := ll.head.url

	ll.head = ll.head.next

	if ll.head == nil {
		ll.tail = nil
	}

	ll.length--
	return data
}

// isEmpty it will check our list is empty or not
func (ll *Queue) IsEmpty() bool {
	ll.RLock()
	defer ll.RUnlock()
	return ll.length == 0
}

// len is return the length of queue
func (ll *Queue) Len() int {
	ll.RLock()
	defer ll.RUnlock()
	return ll.length
}

// frontQueue it will return the front data
func (ll *Queue) FrontQueue() string {
	ll.RLock()
	defer ll.RUnlock()
	return ll.head.url
}

// backQueue it will return the back data
func (ll *Queue) BackQueue() string {
	ll.RLock()
	defer ll.RUnlock()
	return ll.tail.url
}
