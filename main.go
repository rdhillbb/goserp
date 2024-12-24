package main

import "fmt"
import "github.com/rdhillbb/goserp"

func main() {

	results, err := goserp.SerpSearch("What is Garlic")
	if err != nil {
		fmt.Println("Error: ",err)
		return
	}
	fmt.Println("----- Light Search -----")
	fmt.Println(results)

	results, err = goserp.SerpExtensiveSearch("What is Garlic? Is there a benefit for eating onions")
	if err != nil {
		fmt.Println("Error: ",err)
		return
	}
	fmt.Println("\n\n\n", "----- Deep Search -----", results)

}
