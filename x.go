package main

import(
	"fmt"
	"time"
	"github.com/dhkron/huji/utils"
)

func main() {
	fmt.Println("Hello, world!")
	var a []int
	fmt.Printf("Deserializing... ")
	t0 := time.Now()
	utils.DecodeFileIntoArray("/cs/cbio/gil/serialized/GSM455134_30E0LAAXX.2.maq.hic.summary.binned.txt",&a)
	t1 := time.Now()
	fmt.Printf("took %v\r\n", t1.Sub(t0))
}
