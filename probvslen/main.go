package main

import(
	"fmt"
	"time"
	"math"
	"flag"
	"os"
	"bufio"
	"github.com/dhkron/huji/utils"
)

func main() {
	inputFilenameP := flag.String("input","","A serialized input filename.")
	outputFilenameP := flag.String("output","","A serialized input filename.")
	targetChrP := flag.Int("chr",1,"Target chromosome. (X=23, Y=24, M=0)")
	chrSizesFileP := flag.String("chrsize","","Path of chromosome sizes")
	segmentSizeP := flag.Int("segment",100000,"Segment size")
	flag.Parse()

	chrPos, chrLen := utils.GetChrPosAndLens(*chrSizesFileP)
	segment := int64(*segmentSizeP)

	if *outputFilenameP=="" || *inputFilenameP=="" {
		panic("Files, files I say!")
	}

	var a []int
	fmt.Printf("Deserializing... ")
	t0 := time.Now()
	utils.DecodeFileIntoArray(*inputFilenameP,&a)
	t1 := time.Now()
	fmt.Printf("took %v\r\n", t1.Sub(t0))

	dim := int64(math.Sqrt(float64(len(a)))) //Should be initeger
	if dim*dim != int64(len(a)) {
		panic("Matrix size is not a square integer")
	}

	chrStart := chrPos[*targetChrP]/segment
	chrEnd := (chrPos[*targetChrP]+chrLen[*targetChrP])/segment
	hitsVector := make([]int, chrEnd-chrStart)
	hitsVectorSquared := make([]int, chrEnd-chrStart)
	hitCount := 0
	for i:=chrStart; i<chrEnd; i++ {
		for j:=chrStart; j<chrEnd; j++ {
			if i-j>=0 {
				hitsVector[i-j] += a[i+j*dim]
				hitsVectorSquared[i-j] += a[i+j*dim]*a[i+j*dim]
				hitCount += a[i+j*dim]
			}
		}
	}
	probVector := make([]float64, chrEnd-chrStart) //prob is proportional to hit count
	varVector := make([]float64, chrEnd-chrStart) //variance of hit count

	fmt.Printf("Writing data to file\r\n")
	f,err := os.Create(*outputFilenameP)
	if err != nil {
		panic(err)
	}
	w := bufio.NewWriter(f)

	for i:=1; i<len(probVector); i++ { //i=1 because we are skipping self interactions
		//The number of possible contacts at distance X is N-X.
		//For example, the number of successive blocks (distance 1) is N-1
		normalization := 1/float64(len(probVector)-i)
		probVector[i] = float64(hitsVector[i])*normalization
		varVector[i] = float64(hitsVectorSquared[i])*normalization - probVector[i]*probVector[i]
		probVector[i] = probVector[i] / float64(hitCount);
		varVector[i] = varVector[i] / (float64(hitCount)*float64(hitCount))
		fmt.Fprintf(w,"%d,%e,%e\r\n",int64(i)*segment,probVector[i],varVector[i])
	}
	w.Flush()

}
