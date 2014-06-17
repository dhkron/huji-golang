package utils

import(
	"encoding/gob"
	"bufio"
	"os"
)

func EncodeArrayIntoFile(filename string, array *[]int) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	w := bufio.NewWriter(f)
	enc := gob.NewEncoder(w);
	
	err = enc.Encode(*array)
	if err != nil {
		return err
	}
	w.Flush()
	return nil
}

func DecodeFileIntoArray(filename string, array *[]int) error {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	r := bufio.NewReader(f)
	dec := gob.NewDecoder(r);
	
	err = dec.Decode(array)
	if err != nil {
		return err
	}
	return nil
}