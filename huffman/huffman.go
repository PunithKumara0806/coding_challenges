package main

import (
	"container/heap"
	"fmt"
	"os"
)

type HuffMan_Node struct {
	freq  int
	val   rune
	left  *HuffMan_Node
	right *HuffMan_Node
}

type PriorityQueue []*HuffMan_Node

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool { return pq[i].freq < pq[j].freq }

func (pq PriorityQueue) Swap(i, j int) {
	pq[j], pq[i] = pq[i], pq[j]
}

func (pq *PriorityQueue) Push(x any) {
	node := x.(*HuffMan_Node)
	*pq = append(*pq, node)
}

func (pq *PriorityQueue) Pop() any {
	old := *pq
	n := len(old)
	node := old[n-1]
	old[n-1] = nil
	*pq = old[0 : n-1]
	return node
}

type HuffMan_Tree struct {
	minQ   PriorityQueue
	output map[rune]string
	//used to record  current state of the tree traversal for FSM
	currentNode *HuffMan_Node
	//used to recent currentNode after outputing value
	root *HuffMan_Node
}

// Initialises the heap array
func (h *HuffMan_Tree) init(m map[rune]int) {
	h.minQ = make(PriorityQueue, len(m))
	i := 0
	for val, freq := range m {
		h.minQ[i] = &HuffMan_Node{
			freq:  freq,
			val:   val,
			left:  nil,
			right: nil,
		}
		i++
	}
	heap.Init(&(h.minQ))
}

// Creates the huffman tree
func (h *HuffMan_Tree) Create(m map[rune]int) {
	h.init(m)
	n := len(h.minQ)
	h.output = map[rune]string{}
	//algorithm for huffman tree
	for i := 1; i < n; i++ {
		var z HuffMan_Node
		z.left = heap.Pop(&h.minQ).(*HuffMan_Node)
		z.right = heap.Pop(&h.minQ).(*HuffMan_Node)
		z.freq = z.left.freq + z.right.freq
		heap.Push(&h.minQ, &z)
	}
	//root of the huffman tree
	h.root = heap.Pop(&h.minQ).(*HuffMan_Node)
	h.currentNode = h.root
	h.createOutput(h.root)
}

/*
function that traverse the huffman tree for given byte ( 0 or 1 )
the current state is recored in currentNode pointer, if the currentNode lands on a leaf
node it will output the rune, and reset currentNode to root of the tree
*/
func (h *HuffMan_Tree) Move(b byte) rune {
	if b == 1 {
		h.currentNode = h.currentNode.right
		fmt.Printf("move -right : %s\n", string(h.currentNode.val))
	} else {
		h.currentNode = h.currentNode.left
		fmt.Printf("move -left : %s\n", string(h.currentNode.val))
	}
	if h.currentNode.val != 0 {
		c := h.currentNode.val
		h.currentNode = h.root
		return c
	}
	return 0
}

// Creates the respective code for each rune and stores in h.output
func (h *HuffMan_Tree) createOutput(root *HuffMan_Node) {
	current := make([]byte, 0, 8)
	current = append(current, '0')
	h.recursive(root.left, current)
	current[len(current)-1] = '1'
	h.recursive(root.right, current)
}

// helper function that traverse the tree and creates code for each rune as it descends
func (h *HuffMan_Tree) recursive(root *HuffMan_Node, current []byte) {
	if root.val != 0 {
		h.output[root.val] = string(current)
		return
	}
	current = append(current, '0')
	h.recursive(root.left, current)
	current[len(current)-1] = '1'
	h.recursive(root.right, current)
}

// Given slice of byte, it creates a map of m[rune] => count_of_rune
func createCountMap(buf []byte) map[rune]int {
	data := []rune(string(buf))
	c_map := make(map[rune]int)
	for _, c := range data {
		c_map[c]++
	}
	return c_map
}

func bit_writer(h_map map[rune]string, buf []byte) []byte {

	data := []rune(string(buf))
	buffer := make([]byte, 0)
	var mask byte = 1
	var k, i int
	//put k++ here because i would have to increment after every loop
	code := h_map[data[0]]
	var byt byte
	for ; i < len(data); k++ {
		code = h_map[data[i]]
		if k == len(code) {
			// so that k++ makes it 0 again
			k = -1
			i++
			continue
		}
		if code[k] == '1' {
			byt |= mask
		}
		// fmt.Printf("buffer : %b %d,%d\n", buffer[i], i, k)
		// fmt.Printf(":%d\n", mask)
		if mask == 128 {
			mask = 1
			buffer = append(buffer, byt)
			byt = 0
			continue
		}
		mask <<= 1
	}
	return buffer
}
func main() {
	filename := "../test.txt"
	fileP, err := os.Open(filename)
	if err != nil {
		fmt.Println(err)
	}
	buffer := make([]byte, 1024*1024)
	count, err := fileP.Read(buffer)
	if err != nil {
		fmt.Println(err)
	}
	//creating count map
	c_map := createCountMap(buffer[:count])
	var h HuffMan_Tree
	//creating the huffman code
	h.Create(c_map)
	compressed_data, _ := os.Open("compressed.txt")
	read_buf := make([]byte, 100)
	compressed_data.Read(read_buf)
	result := make([]rune, 0, 8)
	var b byte
	var r rune
	for i := 0; i < len(read_buf); i++ {
		b = read_buf[i]
		fmt.Printf("byte : %b\n", b)
		r = 0
		for j := 0; j < 8; j++ {
			if b&(1<<j) > 0 {
				r = h.Move(1)
			} else {
				r = h.Move(0)
			}
			if r != 0 {
				fmt.Printf("output : %s\n", string(r))
				result = append(result, r)
			}
		}
	}
	fmt.Printf("result : %s\n", string(result))
	for key, val := range h.output {
		fmt.Printf("%c : %s\n", key, val)
	}
	// compressed_buf := bit_writer(h.output, buffer[:count])
	// outputP, err := os.OpenFile("compressed.txt", os.O_CREATE|os.O_WRONLY, 0644)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// n, err := outputP.Write(compressed_buf)
	// fmt.Printf("len : %d\n", n)

}
