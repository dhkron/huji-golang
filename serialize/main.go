package main

import(
	"fmt"
	"path/filepath"
	"flag"
	"os"
	"bufio"
	"time"
	"github.com/dhkron/huji/utils"
	"github.com/dhkron/huji/dataformats"
)

func getOutputname(pt string) string {
	//Go 2 up, change to serialized
	a, _ := filepath.Abs(pt + "/../..")
	b := filepath.Base(pt)
	return a+"/serialized/"+b
}

func main() {
	//Flagging around
	segmentP := flag.Int64("segment", 100000, "Segment size")
	inputFilenameP := flag.String("input","","Input file path to be serialized, preferably in /raw directory")
	inputFormatP := flag.String("format","","Format of data")
	flag.Parse()

	if *inputFilenameP == "" {
		fmt.Println("What about an input file, sir?")
		return
	}
	if *inputFormatP == "" {
		fmt.Println("Format please. Guessing could be hazardous.")
		return
	}
	inputFilename := *inputFilenameP
	outputFilename := getOutputname(*inputFilenameP);
	inputFormat :=  *inputFormatP
	segment := *segmentP
	//Print some stuff
	fmt.Printf("Parsing file %s\r\nFormat %s\r\nSegment size %d\r\nOutput file %s\r\n", inputFilename, inputFormat, segment, outputFilename)

	//Load sizes
	chrPos, sum := utils.GetChrPos("chrom.sizes.hg19.txt")

	//Calculate dimensions
	sum = sum /segment
	fmt.Printf("Total dim: %v\n",sum)
	fmt.Printf("Total Pixels: %v\n",sum*sum)

	//Allocate lots of ram
	var s []int
	dim := sum
	pixels := dim*dim
	s = make([]int,pixels)

	//Start parsing
	fmt.Printf("Openning huge file\n")
	var chrIndex1, chrIndex2, inc int64
	var pos1, pos2 int64
	var realPos1, realPos2 int64
	var realIndex int64
	var linecount int64 = 0
	//Go!
	f, _ := os.Open(inputFilename)
	scanner := bufio.NewScanner(f)
	t0 := time.Now()
	for scanner.Scan() {
		chrIndex1, pos1, chrIndex2, pos2, inc = dataformats.Format1(scanner.Bytes())
		realPos1 = chrPos[chrIndex1]+pos1 //Danger!!! chrIndex=0 or 24 will return -1!!!!!
		realPos2 = chrPos[chrIndex2]+pos2
		if chrIndex1 > 0 && chrIndex1<24 { //1 to X
			realIndex = (realPos1 / segment) + dim*(realPos2 / segment)
			s[realIndex] = s[realIndex] + int(inc)
		}
		linecount+=1
		if linecount % 1000000 == 0 {
			t1:= time.Now()
			fmt.Printf("%v : %d so far\n", t1.Sub(t0), linecount)
		}
	}
	t1 := time.Now()
	fmt.Printf("Reading took %v\n", t1.Sub(t0))

	t0 = time.Now()
	utils.EncodeArrayIntoFile(outputFilename, &s)
	t1 = time.Now()
	fmt.Printf("Serializing took %v\n", t1.Sub(t0))
}
