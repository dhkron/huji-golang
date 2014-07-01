package dataformats

import(
	"bytes"
	"strconv"
	"github.com/dhkron/huji/utils"
)

//
func Format1(line []byte) (chrIndex1, pos1, chrIndex2, pos2, inc int64){
	splt := bytes.Split(line,[]byte("\t")) //Copying converts tabs to spaces, watch out
	rsplt := make([][]byte,7)
	rindex := 0
	//Rebuild splt array, skipping empty slices
	for _,v := range splt {
		if len(v) != 0 {
			rsplt[rindex] = v
			rindex+=1
		}
	}
	chrIndex1 = utils.Chr2Int(rsplt[1])
	pos1, _ = strconv.ParseInt(string(rsplt[2]),10,64)
	chrIndex2 = utils.Chr2Int(rsplt[4])
	pos2, _ = strconv.ParseInt(string(rsplt[5]),10,64)
	inc = 1
	return
}

//read name, chromosome1, position1, strand1, restrictionfragment1, chromosome2, position2, strand2, restrictionfragment2 
func Format2(line []byte) (chrIndex1, pos1, chrIndex2, pos2, inc int64){
	splt := bytes.Split(line,[]byte(" "))
	rsplt := make([][]byte,9)
	rindex := 0
	//Rebuild splt array, skipping empty slices
	for _,v := range splt {
		if len(v) != 0 {
			rsplt[rindex] = v
			rindex+=1
		}
	}
	chrIndex1 = utils.Chr2Int(rsplt[1])
	pos1, _ = strconv.ParseInt(string(rsplt[2]),10,64)
	chrIndex2 = utils.Chr2Int(rsplt[4])
	pos2, _ = strconv.ParseInt(string(rsplt[5]),10,64)
	inc = 1
	return
}
