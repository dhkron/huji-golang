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
	modeSmoothP := flag.Bool("smooth",false,"Smooth the data. If true, must supply base and exponent")
	smoothBaseP := flag.Float64("smooth-base",1.0,"Base")
	smoothExpP := flag.Float64("smooth-exp",1.1,"Exponent")
	serTypeP := flag.Int("sertype",0,"Serialization type. 0 for normal, otherwise it is the distance from diagonal");
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

	skipStuff := (*serTypeP>0)

	var dim int64
	var hitsVector []int
	var hitsVectorSquared []int
	var hitCount int
	var diagonalCount int
	var chrStart int64
	var chrEnd int64

	if !skipStuff {
		dim = int64(math.Sqrt(float64(len(a)))) //Should be initeger
		if dim*dim != int64(len(a)) {
			panic("Matrix size is not a square integer")
		}
		chrStart = chrPos[*targetChrP]/segment
		chrEnd = (chrPos[*targetChrP]+chrLen[*targetChrP])/segment
		hitsVector = make([]int, chrEnd-chrStart)
		hitsVectorSquared = make([]int, chrEnd-chrStart)
		hitCount = 0
		diagonalCount = 0
		for i:=chrStart; i<chrEnd; i++ {
			for j:=chrStart; j<chrEnd; j++ {
				if i==j { //Skip self interactions!
					diagonalCount += a[i+j*dim];
				}
				if i-j>0 {
					hitsVector[i-j] += a[i+j*dim]
					hitsVectorSquared[i-j] += a[i+j*dim]*a[i+j*dim]
					hitCount += a[i+j*dim]
				}
			}
		}
	} else {
		dim_h := 2*(*serTypeP)
		dim_w := len(a)/dim_h
		if dim_h*dim_w != len(a) {
			panic("Not a rectangle! Ahhh!")
		}
		fmt.Printf("Dim: h=%v, w=%v\n",dim_h,dim_w)
		chrStart = chrPos[*targetChrP]/segment
		chrEnd = (chrPos[*targetChrP]+chrLen[*targetChrP])/segment
		hitsVector = make([]int, *serTypeP+1) //Make room for self interactions
		hitsVectorSquared = make([]int, *serTypeP+1) //Make room for self interactions
		hitCount = 0
		diagonalCount = 0
		for i:=chrStart; i<chrEnd; i++ {
			for j:=0; j<*serTypeP+1; j++ {
				k := a[int(i)+j*dim_w]
				if j==*serTypeP { //Do not add self interactions to vector
					diagonalCount += k
				} else {
					//j==serType -> self interactions
					//j==0 -> interactions at distance serType
					hitsVector[*serTypeP-j] += k
					hitsVectorSquared[*serTypeP-j] += k*k
					hitCount += k
				}
			}
		}
	}

	probVector := make([]float64, len(hitsVector)) //prob is proportional to hit count
	varVector := make([]float64, len(hitsVectorSquared)) //variance of hit count

	fmt.Printf("Hit count (no diag)\t%d\n",hitCount);
	fmt.Printf("Diagonal hit count\t%d\n",diagonalCount);

	fmt.Printf("Writing data to file\r\n")
	f,err := os.Create(*outputFilenameP)
	if err != nil {
		panic(err)
	}
	w := bufio.NewWriter(f)

	for i:=1; i<len(probVector); i++ { //i=1 because we are skipping self interactions
		//The number of possible contacts at distance X is N-X.
		//For example, the number of successive blocks (distance 1) is N-1
		normalization := 1/float64(int(chrEnd-chrStart)-i)
		probVector[i] = float64(hitsVector[i])*normalization
		varVector[i] = float64(hitsVectorSquared[i])*normalization - probVector[i]*probVector[i]
		probVector[i] = probVector[i] / float64(hitCount);
		varVector[i] = varVector[i] / (float64(hitCount)*float64(hitCount))
		if ! * modeSmoothP {
			fmt.Fprintf(w,"%d,%e,%e\r\n",int64(i)*segment,probVector[i],varVector[i])
		}
	}
	if * modeSmoothP {
		base := *smoothBaseP
		exp := *smoothExpP
		partialSum := 0.0
		partialCount := 0.0
		var newbase float64
		for base < float64(len(probVector)-1) {
			if int(exp*base-base) == 0 {
				newbase = base + 1
				fmt.Fprintf(w,"%d,%e\r\n",int64(base)*segment,probVector[int64(base)])
			} else {
				newbase = math.Min(base * exp, float64(len(probVector)-1))
				partialSum = 0.0;
				partialCount = 0.0;
				for _,val := range probVector[int(base):int(newbase)+1] {
					partialSum += val
					partialCount++
				}
				partialSum = partialSum / partialCount
				fmt.Fprintf(w,"%d,%e\r\n",int64(base)*segment,partialSum)
			}
			base = newbase;
		}
	}
	w.Flush()

}
