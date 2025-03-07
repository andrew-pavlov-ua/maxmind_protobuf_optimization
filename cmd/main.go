package main

import (
	"cmd/internal/models"
	"cmd/internal/services"
	"fmt"
)

const (
	jsonPath      = "./assets/GeoLite2-Country-Test.json"
	protoPath     = "./assets/GeoLite2-Country-Test.proto"
	tempProtoPath = "./assets/GeoLite2-Country-Test-temp.proto"
	mmdbPath      = "./assets/GeoLite2-Country-Test.mmdb"
)

// func main() {
// 	data, err := services.UnmarshalJSON("./assets/GeoLite2-Country-Test.json")
// 	if err != nil {
// 		panic(err)
// 	}

// 	PrintData(data)

// }

func main() {
	err := services.ConvertJSONToProtoFiles(jsonPath, protoPath)
	if err != nil {
		panic(err)
	}

	data, err := services.UnmarshalProtoFile(protoPath)
	if err != nil {
		panic(err)
	}

	PrintData(data)
}

// func main() {
// 	content, err := os.ReadFile(mmdbPath)
// 	if err != nil {
// 		panic(err)
// 	}

// 	db, err := maxminddb.FromBytes(content)
// 	if err != nil {
// 		panic(err)
// 	}
// 	defer db.Close()

// 	geoPairs, err := services.ReadFullProtoFile(protoPath)
// 	if err != nil {
// 		panic(err)
// 	}

// 	for i, pair := range geoPairs.Geos {
// 		ip, _, err := net.ParseCIDR(pair.CIDR)
// 		if err != nil {
// 			panic(err)
// 		}

// 		var result models.MMDBDataItem

// 		db.Lookup(ip, &result)
// 		fmt.Printf("%v) Result:%v \n---------------------------------------------\n", i, result)
// 	}
// }

// PrintData prints all parsed information
func PrintData(dataItems *models.Root) {
	if dataItems == nil || len(dataItems.Geos) == 0 {
		fmt.Println("No data available.")
		return
	}

	fmt.Println("Parsed Data:")
	for cidr, geoId := range dataItems.CidrCountryPairs {
		geo := dataItems.Geos[geoId]
		fmt.Printf("CIDR: %s\n", cidr)
		fmt.Printf("  Continent: %s\n", geo.Continent)
		fmt.Printf("  Country: %s\n", geo.Country)
		fmt.Printf("  Registered Country: %s\n", geo.RegisteredCountry)
		fmt.Println("-----------------------------")
	}
}
