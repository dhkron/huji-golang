package utils

import(
	"encoding/gob"
	"bufio"
	"os"
)

func encodeArrayIntoFile(filename string, array *[]int) {
	f, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	w := bufio.NewWriter(f)
	enc := gob.NewEncoder(w);
	
	err = enc.Encode(*array)
	if err != nil {
		panic(err)
	}
	w.Flush()
}

func decodeFileIntoArray(filename string, array *[]int) {
	f, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	r := bufio.NewReader(f)
	dec := gob.NewDecoder(r);
	
	err = dec.Decode(array)
	if err != nil {
		panic(err)
	}
}