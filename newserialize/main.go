package main

import(
	"fmt"
	"flag"
	"os"
	"bufio"
	"time"
	"github.com/dhkron/huji/utils"
	"github.com/dhkron/huji/dataformats"
)

func main() {
	//Flagging around
	//Split by chr's maybe?
	segmentP := flag.Int64("segment", 10000, "Segment size")
	inputFilenameP := flag.String("input","","Input file path to be serialized, preferably in /raw directory")
	outputFilenameP := flag.String("output","","Specify please")
	inputFormatP := flag.String("format","","Format of data")
	chrSizesP := flag.String("chrsize","","Path of chromosome sizes")
	distFromDiagP := flag.Int64("dist",20,"Distance from diagonal (in segement sized blocks)")
	//stopAfterP := flag.Int64("stop",0,"Stop serializing after N milion records")
	//splitByChrsP := flag.Bool("split",False,"Should split to chromosomes?")
	flag.Parse()

	if *inputFilenameP == "" || *outputFilenameP == ""  {
		fmt.Println("What about an input/output file, sir?")
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
	outputFilename := *outputFilenameP
	dataFormat :=  *inputFormatP
	chrSizesFile := *chrSizesP
	segment := *segmentP
	dist := *distFromDiagP
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

	//Allocate lots of ram
	var s []int
	w_dim := sum
	h_dim := dist*2
	pixels := w_dim*h_dim
	s = make([]int,pixels)

	fmt.Printf("Total dim: %v, %v\n",w_dim,h_dim)
	fmt.Printf("Total Pixels: %v\n",pixels)

	//Start parsing
	fmt.Printf("Openning huge file\n")
	var chrIndex1, chrIndex2, inc int64
	var pos1, pos2 int64
	var realPos1, realPos2 int64
	var realIndex int64
	var linecount int64 = 0
	var blockDist int64
	//Go!
	var myFmt dataformats.DataFormatter
	switch dataFormat {
	case "1":
		myFmt = dataformats.Format1{}
	case "2":
		myFmt = dataformats.Format2{}
	case "3":
		myFmt = dataformats.Format3{}
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
			blockDist = (realPos2-realPos1)/segment
			if blockDist < dist && chrIndex1==chrIndex2 { //Now allowing inter-chr interactions
				//Put in column by real pos, and row by blockdist
				//Maybe a problem with not including these twice?
				realIndex = (realPos1 / segment) + (w_dim * (blockDist+dist))
				if realIndex >= int64(len(s)) {
					panic("Panicking because I am afraid of index out of range")
				}
				if realIndex>=0 {
					s[realIndex] = s[realIndex] + int(inc)
				}
			}
		}
		linecount+=1
		if linecount % 1000000 == 0 {
			t1:= time.Now()
			fmt.Printf("%v : %d so far\n", t1.Sub(t0), linecount)
		}
		//if *stopAfterP > 0 && linecount > 1000000*(*stopAfterP) {
		//	break
		//}
	}
	t1 := time.Now()
	fmt.Printf("Reading took %v\n", t1.Sub(t0))

	t0 = time.Now()
	utils.EncodeArrayIntoFile(outputFilename, &s)
	t1 = time.Now()
	fmt.Printf("Serializing took %v\n", t1.Sub(t0))
}
