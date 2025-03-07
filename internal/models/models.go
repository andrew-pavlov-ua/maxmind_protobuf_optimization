package models

// Example usage with MaxMindDB
type Country struct {
	CIDR    string ` maxminddb:"cidr"`
	Country struct {
		ISOCode string `json:"iso_code" maxminddb:"iso_code"`
	} `json:"country" maxminddb:"country"`
}

type MMDBDataItem struct {
	Continent         *MMDBGeoContinent `maxminddb:"continent" json:"continent,omitempty"`
	Country           *MMDBGeoCountry   `maxminddb:"country" json:"country,omitempty"`
	RegisteredCountry *MMDBGeoCountry   `maxminddb:"registered_country" json:"registered_country,omitempty"`
}

type MMDBGeoContinent struct {
	Code      string        `maxminddb:"code" json:"code,omitempty"`
	GeonameId uint32        `maxminddb:"geoname_id" json:"geoname_id,omitempty"`
	Names     *MMDBGeoNames `maxminddb:"names" json:"names,omitempty"`
}

type MMDBGeoCountry struct {
	GeonameId         uint32        `maxminddb:"geoname_id" json:"geoname_id,omitempty"`
	IsInEuropeanUnion bool          `maxminddb:"is_in_european_union" json:"is_in_european_union,omitempty"`
	IsoCode           string        `maxminddb:"iso_code" json:"iso_code,omitempty"`
	Names             *MMDBGeoNames `maxminddb:"names" json:"names,omitempty"`
}

type MMDBGeoNames struct {
	De   string `maxminddb:"de" json:"de,omitempty"`
	En   string `maxminddb:"en" json:"en,omitempty"`
	Es   string `maxminddb:"es" json:"es,omitempty"`
	Fr   string `maxminddb:"fr" json:"fr,omitempty"`
	Ja   string `maxminddb:"ja" json:"ja,omitempty"`
	PtBR string `maxminddb:"pt-br" json:"pt-br,omitempty"`
	Ru   string `maxminddb:"ru" json:"ru,omitempty"`
	ZhCN string `maxminddb:"zh-cn" json:"zh-cn,omitempty"`
}
