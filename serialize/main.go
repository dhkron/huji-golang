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
	chrSizesP := flag.String("chrsize","","Path of chromosome sizes")
	flag.Parse()

	if *inputFilenameP == "" {
		fmt.Println("What about an input file, sir?")
		return
	}
	if *inputFormatP == "" {
		fmt.Println("Format please. Guessing could be hazardous.")
		return
	}
	if *chrSizesP == "" {
		fmt.Println("Hmm could you add a chromosome sizes file?")
		return
	}
	inputFilename := *inputFilenameP
	outputFilename := getOutputname(*inputFilenameP)
	dataFormat :=  *inputFormatP
	chrSizesFile := *chrSizesP
	segment := *segmentP
	//Print some stuff
	fmt.Printf("Input file\t%s\r\n",inputFilename)
	fmt.Printf("Output file\t%s\r\n",outputFilename)
	fmt.Printf("Data Format\t%s\r\n",dataFormat)
	fmt.Printf("Segment size\t%d\r\n",segment)
	fmt.Printf("Chromosizes\t%s\r\n",chrSizesFile)

	//Load sizes
	chrPos, chrLen := utils.GetChrPosAndLens(chrSizesFile)

	//Calculate dimensions
	sum := (chrPos[23]+chrLen[23])/segment + 1 //+1 For half segment. I could do explicit calc, but must of times it's ok
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
	var myFmt dataformats.DataFormatter
	switch dataFormat {
	case "1":
		myFmt = dataformats.Format1{}
	case "2":
		myFmt = dataformats.Format2{}
	default:
		panic("Unidentified data format! Exiting before running over this huge file!!!")
	}
	f, _ := os.Open(inputFilename)
	scanner := bufio.NewScanner(f)
	t0 := time.Now()
	for scanner.Scan() {
		chrIndex1, pos1, chrIndex2, pos2, inc = myFmt.Format(scanner.Bytes())
		if chrPos[chrIndex1] != -1 && chrPos[chrIndex2] != -1 { //Check if we are in valid chromosomes
			realPos1 = chrPos[chrIndex1]+pos1
			realPos2 = chrPos[chrIndex2]+pos2
			realIndex = (realPos1 / segment) + dim*(realPos2 / segment)
			if realIndex >= int64(len(s)) {
				fmt.Printf("%s ~~~ out of range!!!\r\n",scanner.Bytes())
				fmt.Printf("%d %d %d %d %d\r\n", chrIndex1, pos1, chrIndex2, pos2, inc)
				fmt.Printf("%d\r\n",realIndex)
				panic("Panicking because I am afraid of index out of range")
			}
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
