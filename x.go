package main

import(
	"fmt"
	"github.com/dhkron/huji/utils"
)

func main() {
	fmt.Println("Hello, world!")
	fmt.Printf("%v\r\n",utils.GetChrPos("chrom.sizes.hg19.txt"))
}