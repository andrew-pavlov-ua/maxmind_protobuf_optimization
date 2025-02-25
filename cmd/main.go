package main

import (
	"cmd/internal/services"
	"fmt"
)

func main() {
	err := services.ConvertJSONToProtoFiles("./assets/City-Test.json", "./assets/City-Test.proto")
	if err != nil {
		panic(err)
	}

	result, err := services.ReadProtoFile("./assets/City-Test.proto")
	if err != nil {
		panic(err)
	}

	for _, p := range result {
		fmt.Printf("CIDR: %s, Country: %s\n",
			p.Cidr, p.Country)
	}
	fmt.Println("Countries count: ", len(result))
}
