package main

import (
	"fmt"
)

func bit_writer(data [][]byte) []byte {

	buffer := make([]byte, 8)
	var mask byte = 1
	var j, k, i int
	//put k++ here because i would have to increment after every loop
	for ; i < len(data); k++ {
		if k == len(data[i]) {
			// so that k++ makes it 0 again
			k = -1
			i++
			continue
		}
		if data[i][k] == '1' {
			buffer[j] |= mask
		}
		// fmt.Printf("buffer : %b %d,%d\n", buffer[i], i, k)
		// fmt.Printf(":%d\n", mask)
		if mask == 128 {
			mask = 1
			j++
			continue
		}
		mask <<= 1
	}
	return buffer
}
func main() {
	data := [][]byte{
		[]byte("10110101"),
		[]byte("1100101"),
		[]byte("1001"),
	}
	buffer := bit_writer(data)
	// 0123456789 --> stored as --> 76543210 , 00000098
	for _, v := range buffer {
		fmt.Printf("buffer : %b\n", v)
	}
}
