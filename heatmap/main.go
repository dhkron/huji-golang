package main

import(
	"fmt"
	"time"
	"math"
	"flag"
	"image"
	"image/color"
	"image/jpeg"
	"os"
	"bufio"
	"encoding/csv"
	"strconv"
	"github.com/dhkron/huji/utils"
)

type Image struct {
	data []int
	dim int
	x int
	y int
	w int
	h int
	tilt int
}

func (Image) ColorModel() color.Model {
	return color.RGBAModel
}

func (m Image) Bounds() image.Rectangle {
	return image.Rect(m.x, m.y, m.x+m.w, m.y+m.h)
}

func (m Image) At(x,y int) color.Color {
	var v uint8
	v = uint8(m.data[x+m.dim*y])
	return color.RGBA{v,v,255,255}
}

func main() {
	inputFilenameP := flag.String("input","","A serialized input filename.")
	serTypeP := flag.Int("sertype",0,"Serialization type. 0 is square, 1 is rectangle")
	outputFilenameP := flag.String("output","image.jpeg","Output filename")
	targetChr1P := flag.Int("chr1",1,"First target chromosome. (X=23, Y=24, M=0)")
	targetChr2P := flag.Int("chr2",1,"Second target chromosome. (X=23, Y=24, M=0)")
	chrSizesFileP := flag.String("chrsize","","Path of chromosome sizes")
	segmentSizeP := flag.Int("segment",100000,"Segment size")
	magicNumberP := flag.Int("magic",1,"Magic number")
	thersholdP := flag.Int("ther",100,"Hit counts higher than this are removed")
	normalizeP := flag.Float64("norm",0,"Power law normalization")
	normalizeWithSourceP := flag.String("normsource","","Normalize using source file")
	flag.Parse()

	chrPos, chrLen := utils.GetChrPosAndLens(*chrSizesFileP)
	segment := int64(*segmentSizeP)

	var a []int
	fmt.Printf("Deserializing... ")
	t0 := time.Now()
	utils.DecodeFileIntoArray(*inputFilenameP,&a)
	t1 := time.Now()
	fmt.Printf("took %v\r\n", t1.Sub(t0))

	skipStuff := (*serTypeP>0)

	var dim int
	var domain1pos int
	var domain2pos int
	var domain1len int
	var domain2len int

	if !skipStuff {
		dim = int(math.Sqrt(float64(len(a)))) //Should be initeger
		if dim*dim != len(a) {
			panic("Matrix size is not a square integer")
		}
		domain1pos = int(chrPos[*targetChr1P]/segment)
		domain2pos = int(chrPos[*targetChr2P]/segment)
		domain1len = int(chrLen[*targetChr1P]/segment)
		domain2len = int(chrLen[*targetChr2P]/segment)
	} else {
		dim_h := *serTypeP
		dim_w := len(a)/dim_h
		dim = dim_w
		if dim_h*dim_w != len(a) {
			panic("Rectnangle height doesn't match it's length")
		}
		fmt.Printf("Dim: h=%v, w=%v\n",dim_h,dim_w)
		domain1pos = int(chrPos[*targetChr1P]/segment)
		domain1len = int(chrLen[*targetChr1P]/segment)
		domain2pos = 0
		domain2len = dim_h
	}

	//Initiate normalization map
	normMap := make([]float64, domain1len, domain1len)
	if *normalizeWithSourceP != "" {
		fmt.Println("Normalizing with source file")
		fcsv,err := os.Open(*normalizeWithSourceP)
		if err != nil {
			panic(err)
		}
		csvreader := csv.NewReader(bufio.NewReader(fcsv))
		normVector,err := csvreader.ReadAll();
		var key uint64
		var val float64
		for _,row := range normVector {
			key,_ = strconv.ParseUint(row[0],10,64)
			val,_ = strconv.ParseFloat(row[1],64)
			normMap[key/uint64(segment)]=val
		}
		for index,value := range normMap {
			if value == 0 && index>0 {
				normMap[index]=normMap[index-1]
			}
		}
	}

	//Enforce chr1=chr2 while normalizing
	if *targetChr1P == *targetChr2P && (*normalizeP>0 || *normalizeWithSourceP != "") {
		for x:=domain1pos; x<domain1pos+domain1len; x++ {
			for y:=domain2pos; y<domain2pos+domain2len; y++{
				key := x + y*dim
				dist := x-y
				if skipStuff {
					dist = y+(*serTypeP)/2
				}
				if dist<0 {
					dist = -dist
				}
				if dist == 0 {
					dist = 1
				}
				if * normalizeWithSourceP != "" {
					a[key] = int(float64(a[key]) / (normMap[dist]*100000))
				} else if *normalizeP > 0 {
					a[key] = a[key]*int(math.Pow(float64(dist),*normalizeP))
				}
			}
		}
	}

	//Do thersholding by value, and magicing by distance from diag
	min := 0
	max := int(math.Inf(1))
	for x:=domain1pos; x<domain1pos+domain1len; x++ {
		for y:=domain2pos; y<domain2pos+domain2len; y++{
			key := x + y*dim
			value := a[key]
			if !skipStuff && (((x-y)<*magicNumberP && (y-x)<*magicNumberP) || value>*thersholdP ) {
				a[key]=255
			} else {
				if value>max {
					max = value
				}
				if value<min {
					min = value
				}
			}
		}
	}
	fmt.Printf("Min %d Max %d\n",min,max)
	scalingFactor := 255.0/math.Log1p(float64(max))
	fmt.Printf("rescalling using %f... ",scalingFactor)
	t0 = time.Now()
	for x:=domain1pos; x<domain1pos+domain1len; x++ {
		for y:=domain2pos; y<domain2pos+domain2len; y++ {
			key := x+y*dim
			if a[key] == 255 {
				a[key] = 0
			} else {
				a[key] = 255-int(math.Log1p(float64(a[key]))*scalingFactor)
			}
		}
	}
	t1 = time.Now()
	fmt.Printf("Took %v\n",t1.Sub(t0))

	m := Image{a, dim, domain1pos, domain2pos, domain1len, domain2len, *serTypeP}
	fmt.Printf("Image domain: %v\n", m.Bounds())

	f,err := os.Create(*outputFilenameP)
	if err != nil {
		panic(err)
	}
	w := bufio.NewWriter(f)
	fmt.Printf("Encoding %d %d... ",*targetChr1P,*targetChr2P)
	t0 = time.Now()
	err = jpeg.Encode(w,m,nil)
	t1 = time.Now()
	fmt.Printf("took %v\n",t1.Sub(t0))
	if err!=nil {
		panic(err)
	}
	w.Flush()
}
