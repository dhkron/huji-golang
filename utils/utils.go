package utils

import(
	"encoding/gob"
	"bufio"
	"os"
	"fmt"
	"strconv"
)

func EncodeArrayIntoFile(filename string, array *[]int) error
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

//Maps M to 0, Y to 24, X to 23, and all the rest to numbers. Invalid datas are -1
func Chr2Int(chrname []byte) (int64,error) {
	if len(chrname)==1 {
		switch chrname[0] {
			case 'm','M':
				return 0,nil
			case 'x','X':
				return 23,nil
			case 'y','Y':
				return 24,nil
		}
	}
	num,err := strconv.ParseInt(string(chrname),10,64)
	if err!=nil {
		return -1,err
	}
	return num,nil
}

//Returns an array with mapping of chr number into position on the genome
func GetChrPosAndLens(filename string) (chrPos, chrLen []int64) {
	chr_file, err := os.Open(filename)
	if err != nil {
		fmt.Printf("Something went wrong while reading %s\r\n",filename)
		return
	}
	//Data holders
	chrLen = make([]int64,25) //0 is M, 23 X, 24 Y
	chrPos = make([]int64,25) //0 is M, 23 X, 24 Y
	//Scanner state machine
	var wordState int = 0
	var cLength int64
	var prvChr int64
	var ln []byte
	//Initiate scanner
	scanner := bufio.NewScanner(chr_file)
	scanner.Split(bufio.ScanWords)
	for scanner.Scan() {
		ln = scanner.Bytes()
		if wordState == 0 {
			prvChr,err = Chr2Int(ln[3:])
		} else {
			cLength,err = strconv.ParseInt(string(ln),10,64)
			if prvChr>=0 {
				chrLen[prvChr] = int64(cLength)
			}
		}
		wordState = 1-wordState
	}
	chrPos[0] = -1 //M is invalid
	chrPos[1] = 0 //Chr1 starts from 0
	for i := 2; i<23; i++ {
		chrPos[i] = chrLen[i-1]+chrPos[i-1]
	}
	chrPos[23] = chrLen[22]+chrPos[22] //X comes right after 22
	chrPos[24] = -1 //Y is invalid
	return
}
