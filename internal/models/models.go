package models

type Countries struct {
	Names map[string]string
}

type CountriesData struct {
	IP      string
	GeoData *Countries
}
