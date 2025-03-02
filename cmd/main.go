package main

import (
	"cmd/internal/services"
	"fmt"
)

func main() {
	data, err := services.UnmarshalJSON("./assets/GeoLite2-Country-Test.json")
	if err != nil {
		panic(err)
	}

	for i, pair := range data.Pairs {
		fmt.Printf("%v) cidr - iso: %v - %v\n", i, pair.Cidr, pair.Isocode)
	}
}

// func main() {
// 	err := services.ConvertJSONToProtoFiles("./assets/GeoLite2-Country-Test.json", "./assets/GeoLite2-Country-Test.proto")
// 	if err != nil {
// 		panic(err)
// 	}

// 	pairs, err := services.ReadFullProtoFile("./assets/GeoLite2-Country-Test.proto")
// 	if err != nil {
// 		panic(err)
// 	}

// 	for _, pair := range pairs.Pairs {
// 		fmt.Printf("CIDR: %v, ISO: %v\n", pair.Cidr, pair.Country)
// 	}
// }

// func main() {
// 	db, err := maxminddb.Open("./assets/GeoLite2-Country-Test.mmdb")
// 	if err != nil {
// 		panic(err)
// 	}
// 	defer db.Close()
// 	var record models.Country

// 	addr := net.ParseIP("::27d:a0d8")
// 	err = db.Lookup(addr, &record)
// 	if err != nil {
// 		panic(err)
// 	}

// 	fmt.Printf("%s\n", record.Country.ISOCode)
// }
