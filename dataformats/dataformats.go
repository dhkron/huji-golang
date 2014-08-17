package dataformats

import(
	"bytes"
	"strconv"
	"github.com/dhkron/huji/utils"
)

type DataFormatter interface {
	Format([]byte) (a,b,c,d,e int64)
}

type Format1 struct {}
func (Format1) Format(line []byte) (chrIndex1, pos1, chrIndex2, pos2, inc int64){
	splt := bytes.Split(line,[]byte("\t")) //Watch out for tabs converting into spaces!!!
	/*rsplt := make([][]byte,7)
	rindex := 0
	//Rebuild splt array, skipping empty slices
	for _,v := range splt {
		if len(v) != 0 {
			rsplt[rindex] = v
			rindex+=1
		}
	}*/
	chrIndex1,_ = utils.Chr2Int(splt[1][3:])
	pos1, _ = strconv.ParseInt(string(splt[2]),10,64)
	chrIndex2,_ = utils.Chr2Int(splt[4][3:])
	pos2, _ = strconv.ParseInt(string(splt[5]),10,64)
	inc = 1
	return
}

type Format2 struct {}
//read name, chromosome1, position1, strand1, restrictionfragment1, chromosome2, position2, strand2, restrictionfragment2 
func (Format2) Format(line []byte) (chrIndex1, pos1, chrIndex2, pos2, inc int64){
	splt := bytes.Split(line,[]byte("\t"))
	/*rsplt := make([][]byte,9)
	rindex := 0
	//Rebuild splt array, skipping empty slices
	for _,v := range splt {
		if len(v) != 0 {
			rsplt[rindex] = v
			rindex+=1
		}
	}*/
	chrIndex1,_ = utils.Chr2Int(splt[1])
	pos1,_ = strconv.ParseInt(string(splt[2]),10,64)
	chrIndex2,_ = utils.Chr2Int(splt[5])
	pos2,_ = strconv.ParseInt(string(splt[6]),10,64)
	inc = 1
	return
}

type Format3 struct {}
//chrms1,chrms2,cuts1,cuts2,strands1,strands2
//chr1 is 0, chr23 is 22. No M/X/Y are present
func (Format3) Format(line []byte) (chrIndex1, pos1, chrIndex2, pos2, inc int64){
	splt := bytes.Split(line,[]byte(","))
	chrIndex1,_ = utils.Chr2Int(splt[0])
	chrIndex1++
	pos1,_ = strconv.ParseInt(string(splt[2]),10,64)
	chrIndex2,_ = utils.Chr2Int(splt[1])
	chrIndex2++
	pos2,_ = strconv.ParseInt(string(splt[3]),10,64)
	inc = 1
	return
}
