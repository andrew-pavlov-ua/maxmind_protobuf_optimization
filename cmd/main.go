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
		fmt.Printf("Original IP: %s, Country: %s\n",
			p.IP, p.Country)
	}
	fmt.Println("Countries count: ", len(result))
}
