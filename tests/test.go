package main

import "fmt"
import "os"
import "github.com/ishiikurisu/edf"

func main() {
	source := os.Args[1]
	header, _ := edf.ReadFile(source)
	sampling := edf.GetSampling(header)

	fmt.Printf("# Reading EDF files to extract the sampling\n")
	fmt.Printf("---\n")
	fmt.Printf("file: %v\n", source)
	fmt.Printf("header: %v\n", header["samplesrecord"])
	fmt.Printf("sampling: %v\n", sampling)
}
