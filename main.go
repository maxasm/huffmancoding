package main

import (
	"fmt"
	"os"
)

type Node struct {
	lNode     *Node // The left child node
	rNode     *Node // The right child node
	symbol    byte  // The symbol 'held' by the node
	frequency int   // The freqeuncy of the symbol
}

// sort the nodes based on the frequency of the symbols
func sort(arr *[]Node) {
	aLen := len(*arr)
	for a := 0; a < aLen-1; a++ {
		for b := 0; b < aLen-1; b++ {
			temp := (*arr)[b]
			if (*arr)[b+1].frequency < temp.frequency {
				(*arr)[b] = (*arr)[b+1]
				(*arr)[b+1] = temp
			}
		}
	}
}

// Helper function to merger [0] and [1] using a common parent node
func merge(arr []Node) []Node {
	aLen := len(arr)
	_arr := make([]Node, aLen-1)

	// create a common parent for index 0 and index 1
	pNode := Node{
		lNode:     &arr[0],
		rNode:     &arr[1],
		symbol:    0x00,
		frequency: arr[0].frequency + arr[1].frequency,
	}

	// copy all the remaining nodes from index 2
	for a := 2; a < aLen; a++ {
		_arr[a-1] = arr[a]
	}

	// set index 0 to the new parent node
	_arr[0] = pNode

	// sort the []Node
	sort(&_arr)

	return _arr
}

// Helper function to get the frequncy table of symbols from a []byte
func getFrequencyTable(arr []byte) map[byte]int {
	mp := map[byte]int{}
	for a := range arr {
		bt := arr[a]
		if freq, ok := mp[bt]; ok {
			mp[bt] = freq + 1
		} else {
			mp[bt] = 1
		}
	}
	return mp
}

// Helper function to get the array of Nodes
func getNodeArray(ft map[byte]int) []Node {
	nArray := []Node{}
	for sym, freq := range ft {
		n := Node{
			lNode:     nil,
			rNode:     nil,
			symbol:    sym,
			frequency: freq,
		}
		nArray = append(nArray, n)
	}
	return nArray
}

func trav(n *Node, code int16, codeLen int16) {
	if n.lNode == nil && n.rNode == nil {
		fmt.Printf("%s -> %c\n", codeToString(code, codeLen), n.symbol)
	} else {
		trav(n.lNode, ((code << 1) | 1), codeLen+1)
		trav(n.rNode, ((code << 1) | 0), codeLen+1)
	}
}

func codeToString(code int16, codeLen int16) string {
	var str string
	str = fmt.Sprintf("%b", code)
	diff := int(codeLen) - len(str)
	for a := 0; a < diff; a++ {
		str = "0" + str
	}
	return str
}

func main() {
	// read bytes from a file
	if len(os.Args) < 2 {
		os.Exit(1)
	}

	filename := os.Args[1]
	f, err := os.Open(filename)
	if err != nil {
		os.Exit(1)
	}

	bArray := make([]byte, 100)
	f.Read(bArray)
	f.Close()
	// In the JPEG encoding just use the frequency tables for the dc symbols and the ac symbols
	//bArray := []byte{0, 1, 2}
	ft := getFrequencyTable(bArray)
	//fmt.Printf("%v\n", ft)
	nArray := getNodeArray(ft)
	// sort the array
	sort(&nArray)
	//	fmt.Printf("%v\n", nArray)
	max := len(nArray) - 1
	for a := 0; a < max; a++ {
		res := merge(nArray)
		//		fmt.Printf("%v\n", res)
		nArray = res
	}
	trav(&nArray[0], 0, 0)
}
