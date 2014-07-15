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
	"github.com/dhkron/huji/utils"
)

type Image struct {
	data []int
	dim int
	x int
	y int
	w int
	h int
}

func (Image) ColorModel() color.Model {
	return color.RGBAModel
}

func (m Image) Bounds() image.Rectangle {
	return image.Rect(m.x, m.y, m.x+m.w, m.y+m.h)
}

func (m Image) At(x,y int) color.Color {
	var v uint8 = uint8(m.data[x+m.dim*y])
	return color.RGBA{v,v,255,255}
}

func main() {
	inputFilenameP := flag.String("input","","A serialized input filename.")
	outputFilenameP := flag.String("output","image.jpeg","Output filename")
	targetChr1P := flag.Int("chr1",1,"First target chromosome. (X=23, Y=24, M=0)")
	targetChr2P := flag.Int("chr2",1,"Second target chromosome. (X=23, Y=24, M=0)")
	chrSizesFileP := flag.String("chrsize","","Path of chromosome sizes")
	segmentSizeP := flag.Int("segment",100000,"Segment size")
	magicNumberP := flag.Int("magic",1,"Magic number")
	thersholdP :=  flag.Int("ther",100,"Magic number")
	flag.Parse()

	chrPos, chrLen := utils.GetChrPosAndLens(*chrSizesFileP)
	segment := int64(*segmentSizeP)

	var a []int
	fmt.Printf("Deserializing... ")
	t0 := time.Now()
	utils.DecodeFileIntoArray(*inputFilenameP,&a)
	t1 := time.Now()
	fmt.Printf("took %v\r\n", t1.Sub(t0))

	dim := int(math.Sqrt(float64(len(a)))) //Should be initeger
	if dim*dim != len(a) {
		panic("Matrix size is not a square integer")
	}

	domain1pos := int(chrPos[*targetChr1P]/segment)
	domain2pos := int(chrPos[*targetChr2P]/segment)
	domain1len := int(chrLen[*targetChr1P]/segment)
	domain2len := int(chrLen[*targetChr2P]/segment)

	min := 0
	max := int(math.Inf(1))
	for x:=domain1pos; x<domain1pos+domain1len; x++ {
		for y:=domain2pos; y<domain2pos+domain2len; y++{
			key := x + y*dim
			value := a[key]
			if ((x-y)<*magicNumberP && (y-x)<*magicNumberP) || value>*thersholdP {
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

	m := Image{a, dim, domain1pos, domain2pos, domain1len, domain2len}
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
