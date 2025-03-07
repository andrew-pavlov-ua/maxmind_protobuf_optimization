package test

import (
	"cmd/internal/models"
	"cmd/internal/services"
	"net"
	"os"
	"testing"

	"github.com/oschwald/maxminddb-golang"
	"github.com/stretchr/testify/require"
)

const (
	JSONCountryFilePath  = "../assets/GeoLite2-Country-Test.json"
	ProtoCountryFilePath = "../assets/GeoLite2-Country-Test.proto"
	MMDBCountryFilePath  = "../assets/GeoLite2-Country-Test.mmdb"
)

var testTable = []struct {
	CIDR    string
	ISOCode string
}{
	{CIDR: "2.125.160.216/29", ISOCode: "GB"},
	{CIDR: "50.114.0.0/22", ISOCode: "US"},
	{CIDR: "67.43.156.0/24", ISOCode: "BT"},
	{CIDR: "81.2.69.142/31", ISOCode: "GB"},
	{CIDR: "81.2.69.144/28", ISOCode: "GB"},
	{CIDR: "81.2.69.160/27", ISOCode: "GB"},
	{CIDR: "81.2.69.192/28", ISOCode: "GB"},
	{CIDR: "89.160.20.112/28", ISOCode: "SE"},
	{CIDR: "89.160.20.128/25", ISOCode: "SE"},
	{CIDR: "111.235.160.0/22", ISOCode: "CN"},
	{CIDR: "202.196.224.0/20", ISOCode: "PH"},
	{CIDR: "216.160.83.56/29", ISOCode: "US"},
	{CIDR: "217.65.48.0/29", ISOCode: "GI"},
	{CIDR: "2001:218::/32", ISOCode: "JP"},
	{CIDR: "2001:220::/32", ISOCode: "KR"},
	{CIDR: "2001:230::/32", ISOCode: "KR"},
	{CIDR: "2001:238::/32", ISOCode: "TW"},
	{CIDR: "2001:240::/32", ISOCode: "JP"},
	{CIDR: "2001:250::/31", ISOCode: "CN"},
	{CIDR: "2001:252::/32", ISOCode: "CN"},
	{CIDR: "2001:254::/32", ISOCode: "CN"},
	{CIDR: "2001:256::/32", ISOCode: "CN"},
	{CIDR: "2001:258::/32", ISOCode: "JP"},
	{CIDR: "2001:260::/32", ISOCode: "JP"},
	{CIDR: "2001:268::/32", ISOCode: "JP"},
	{CIDR: "2001:270::/32", ISOCode: "KR"},
	{CIDR: "2001:278::/32", ISOCode: "JP"},
	{CIDR: "2001:280::/32", ISOCode: "KR"},
	{CIDR: "2001:288::/32", ISOCode: "TW"},
	{CIDR: "2001:290::/32", ISOCode: "KR"},
	{CIDR: "2001:298::/32", ISOCode: "JP"},
	{CIDR: "2001:2a0::/32", ISOCode: "JP"},
	{CIDR: "2001:2a8::/32", ISOCode: "JP"},
	{CIDR: "2001:2b0::/32", ISOCode: "KR"},
	{CIDR: "2001:2b8::/32", ISOCode: "KR"},
	{CIDR: "2001:2c0::/32", ISOCode: "JP"},
	{CIDR: "2001:2c8::/32", ISOCode: "JP"},
	{CIDR: "2001:2d8::/32", ISOCode: "KR"},
	{CIDR: "2001:2e0::/32", ISOCode: "HK"},
	{CIDR: "2001:2e8::/32", ISOCode: "JP"},
	{CIDR: "2001:2f0::/32", ISOCode: "JP"},
	{CIDR: "2001:2f8::/32", ISOCode: "JP"},
	{CIDR: "2a02:cf40::/29", ISOCode: "NO"},
	{CIDR: "2a02:cf80::/29", ISOCode: "IL"},
	{CIDR: "2a02:cfc0::/29", ISOCode: "FR"},
	{CIDR: "2a02:d000::/29", ISOCode: "CH"},
	{CIDR: "2a02:d040::/29", ISOCode: "SE"},
	{CIDR: "2a02:d080::/29", ISOCode: "BH"},
	{CIDR: "2a02:d0c0::/29", ISOCode: "RU"},
	{CIDR: "2a02:d100::/29", ISOCode: "PL"},
	{CIDR: "2a02:d140::/29", ISOCode: "NO"},
	{CIDR: "2a02:d180::/29", ISOCode: "DE"},
	{CIDR: "2a02:d1c0::/29", ISOCode: "IT"},
	{CIDR: "2a02:d200::/29", ISOCode: "FI"},
	{CIDR: "2a02:d240::/29", ISOCode: "BY"},
	{CIDR: "2a02:d280::/29", ISOCode: "CZ"},
	{CIDR: "2a02:d2c0::/29", ISOCode: "IR"},
	{CIDR: "2a02:d300::/29", ISOCode: "UA"},
	{CIDR: "2a02:d340::/29", ISOCode: "FR"},
	{CIDR: "2a02:d380::/29", ISOCode: "IR"},
	{CIDR: "2a02:d3c0::/29", ISOCode: "GB"},
	{CIDR: "2a02:d400::/29", ISOCode: "HU"},
	{CIDR: "2a02:d440::/29", ISOCode: "SE"},
	{CIDR: "2a02:d480::/29", ISOCode: "DE"},
	{CIDR: "2a02:d4c0::/30", ISOCode: "FI"},
	{CIDR: "2a02:d4e0::/30", ISOCode: "DE"},
	{CIDR: "2a02:d500::/29", ISOCode: ""},
	{CIDR: "2a02:d540::/29", ISOCode: "GB"},
	{CIDR: "2a02:d580::/29", ISOCode: "FR"},
	{CIDR: "2a02:d5c0::/29", ISOCode: "ES"},
	{CIDR: "2a02:d600::/29", ISOCode: "DE"},
	{CIDR: "2a02:d640::/29", ISOCode: "FR"},
	{CIDR: "2a02:d680::/30", ISOCode: "GB"},
	{CIDR: "2a02:d6a0::/30", ISOCode: "DE"},
	{CIDR: "2a02:d6c0::/29", ISOCode: "BG"},
	{CIDR: "2a02:d700::/29", ISOCode: "DE"},
	{CIDR: "2a02:d740::/29", ISOCode: "CH"},
	{CIDR: "2a02:d780::/29", ISOCode: "IR"},
	{CIDR: "2a02:d7c0::/29", ISOCode: "FR"},
	{CIDR: "2a02:d800::/29", ISOCode: "RO"},
	{CIDR: "2a02:d840::/29", ISOCode: "RU"},
	{CIDR: "2a02:d880::/29", ISOCode: "RU"},
	{CIDR: "2a02:d8c0::/29", ISOCode: "NO"},
	{CIDR: "2a02:d900::/29", ISOCode: "SE"},
	{CIDR: "2a02:d940::/29", ISOCode: "BE"},
	{CIDR: "2a02:d980::/29", ISOCode: "TR"},
	{CIDR: "2a02:d9c0::/29", ISOCode: "TR"},
	{CIDR: "2a02:da00::/29", ISOCode: "DE"},
	{CIDR: "2a02:da40::/29", ISOCode: "GB"},
	{CIDR: "2a02:da80::/29", ISOCode: "AT"},
	{CIDR: "2a02:dac0::/29", ISOCode: "RU"},
	{CIDR: "2a02:db00::/29", ISOCode: "DE"},
	{CIDR: "2a02:db40::/29", ISOCode: "RO"},
	{CIDR: "2a02:db80::/29", ISOCode: "RU"},
	{CIDR: "2a02:dbc0::/29", ISOCode: "RU"},
	{CIDR: "2a02:dc00::/29", ISOCode: "RU"},
	{CIDR: "2a02:dc40::/29", ISOCode: "TR"},
	{CIDR: "2a02:dc80::/29", ISOCode: "RU"},
	{CIDR: "2a02:dcc0::/29", ISOCode: "UA"},
	{CIDR: "2a02:dd00::/29", ISOCode: "AL"},
	{CIDR: "2a02:dd40::/29", ISOCode: "GB"},
	{CIDR: "2a02:dd80::/29", ISOCode: "SE"},
	{CIDR: "2a02:ddc0::/29", ISOCode: "RU"},
	{CIDR: "2a02:de00::/29", ISOCode: "RU"},
	{CIDR: "2a02:de40::/29", ISOCode: "IL"},
	{CIDR: "2a02:de80::/29", ISOCode: "RU"},
	{CIDR: "2a02:dec0::/29", ISOCode: "LB"},
	{CIDR: "2a02:df00::/29", ISOCode: "IR"},
	{CIDR: "2a02:df40::/29", ISOCode: "TR"},
	{CIDR: "2a02:df80::/29", ISOCode: "GB"},
	{CIDR: "2a02:dfc0::/29", ISOCode: "IR"},
	{CIDR: "2a02:e000::/29", ISOCode: "FR"},
	{CIDR: "2a02:e040::/29", ISOCode: "NL"},
	{CIDR: "2a02:e080::/29", ISOCode: "KW"},
	{CIDR: "2a02:e0c0::/29", ISOCode: "CH"},
	{CIDR: "2a02:e100::/29", ISOCode: "GB"},
	{CIDR: "2a02:e140::/29", ISOCode: "PL"},
	{CIDR: "2a02:e180::/29", ISOCode: "GB"},
	{CIDR: "2a02:e1c0::/29", ISOCode: "NL"},
	{CIDR: "2a02:e200::/30", ISOCode: "AT"},
	{CIDR: "2a02:e220::/30", ISOCode: "SA"},
	{CIDR: "2a02:e240::/29", ISOCode: "DE"},
	{CIDR: "2a02:e280::/29", ISOCode: "DE"},
	{CIDR: "2a02:e2c0::/29", ISOCode: "IT"},
	{CIDR: "2a02:e300::/29", ISOCode: "BY"},
	{CIDR: "2a02:e340::/29", ISOCode: "NO"},
	{CIDR: "2a02:e380::/29", ISOCode: "IT"},
	{CIDR: "2a02:e3c0::/29", ISOCode: "FR"},
	{CIDR: "2a02:e400::/29", ISOCode: "SE"},
	{CIDR: "2a02:e440::/29", ISOCode: "DE"},
	{CIDR: "2a02:e480::/29", ISOCode: "RU"},
	{CIDR: "2a02:e4c0::/29", ISOCode: "NL"},
	{CIDR: "2a02:e500::/29", ISOCode: "FR"},
	{CIDR: "2a02:e540::/29", ISOCode: "RS"},
	{CIDR: "2a02:e580::/29", ISOCode: "NO"},
	{CIDR: "2a02:e5c0::/29", ISOCode: "RU"},
	{CIDR: "2a02:e600::/30", ISOCode: "FR"},
	{CIDR: "2a02:e620::/30", ISOCode: "RU"},
	{CIDR: "2a02:e640::/29", ISOCode: "FR"},
	{CIDR: "2a02:e680::/29", ISOCode: "JO"},
	{CIDR: "2a02:e6c0::/29", ISOCode: "RU"},
	{CIDR: "2a02:e700::/29", ISOCode: "LY"},
	{CIDR: "2a02:e740::/29", ISOCode: "DE"},
	{CIDR: "2a02:e780::/29", ISOCode: "TR"},
	{CIDR: "2a02:e7c0::/29", ISOCode: "PL"},
	{CIDR: "2a02:e800::/29", ISOCode: "DE"},
	{CIDR: "2a02:e840::/29", ISOCode: "RU"},
	{CIDR: "2a02:e880::/29", ISOCode: "RU"},
	{CIDR: "2a02:e900::/29", ISOCode: "IE"},
	{CIDR: "2a02:e940::/29", ISOCode: "RO"},
	{CIDR: "2a02:e980::/29", ISOCode: "IL"},
	{CIDR: "2a02:e9c0::/29", ISOCode: "SE"},
	{CIDR: "2a02:ea00::/29", ISOCode: "CH"},
	{CIDR: "2a02:ea40::/29", ISOCode: "GB"},
	{CIDR: "2a02:ea80::/29", ISOCode: "PL"},
	{CIDR: "2a02:eac0::/29", ISOCode: "PL"},
	{CIDR: "2a02:eb00::/29", ISOCode: "RU"},
	{CIDR: "2a02:eb40::/29", ISOCode: "GB"},
	{CIDR: "2a02:eb80::/29", ISOCode: "RU"},
	{CIDR: "2a02:ebc0::/29", ISOCode: "FR"},
	{CIDR: "2a02:ec00::/29", ISOCode: "FR"}}

func TestLookUpCountriesMmdb(t *testing.T) {
	content, err := os.ReadFile(MMDBCountryFilePath)
	require.NoError(t, err)

	db, err := maxminddb.FromBytes(content)
	require.NoError(t, err)
	defer db.Close()

	for _, i := range testTable {
		var record models.Country

		ip, _, err := net.ParseCIDR(i.CIDR)
		require.NoError(t, err)

		err = db.Lookup(ip, &record)
		t.Logf("Ip: %v", i.CIDR)
		require.NoError(t, err)

		t.Logf("Current case IP and Country: %v - %v, got: %v - %v", i.CIDR, i.ISOCode, ip, record.Country.ISOCode)

		require.Equal(t, i.ISOCode, record.Country.ISOCode)
	}
}

func TestLookUpCountriesProto(t *testing.T) {
	// Read the full Proto file and check for errors
	root, err := services.UnmarshalProtoFile(ProtoCountryFilePath)
	require.NoError(t, err)

	// Build a map from CIDR to ISO code for fast lookups.
	cidrToIso := make(map[string]string, len(root.CidrCountryPairs))
	for cidr, index := range root.CidrCountryPairs {
		var iso string
		if root.Geos[index].Country != nil {
			iso = root.Geos[index].Country.IsoCode
		} else if root.Geos[index].RegisteredCountry != nil {
			iso = root.Geos[index].RegisteredCountry.IsoCode
		}
		cidrToIso[cidr] = iso
	}

	// For each test case, check that the expected ISO code matches the one in the map
	for _, tc := range testTable {
		iso, ok := cidrToIso[tc.CIDR]
		require.True(t, ok, "CIDR %v not found in proto data", tc.CIDR)
		t.Logf("Current case IP and Country: %v - expected: %v, got: %v", tc.CIDR, tc.ISOCode, iso)
		require.Equal(t, tc.ISOCode, iso)
	}
}

var dummyCount int

// BenchmarkLookUpCountriesMmdb benchmarks the lookup performance using the MMDB file
func BenchmarkLookUpCountriesMmdb(b *testing.B) {
	dummyCount = 0
	db, err := maxminddb.Open(MMDBCountryFilePath)
	if err != nil {
		b.Fatal(err)
	}
	defer db.Close()

	// Pre-parse the IP addresses from the test table
	ips := make([]net.IP, len(testTable))
	for i, tc := range testTable {
		ip, _, err := net.ParseCIDR(tc.CIDR)
		if err != nil {
			b.Fatal(err)
		}
		ips[i] = ip
	}

	// Reset the timer before the benchmark loop
	for i := 0; b.Loop(); i++ {
		ip := ips[i%len(ips)]
		var record models.Country
		if err := db.Lookup(ip, &record); err != nil {
			b.Fatal(err)
		}
		dummyCount += len(record.Country.ISOCode)
	}
}

// BenchmarkLookUpCountriesProto benchmarks the lookup performance using the Proto file
func BenchmarkLookUpCountriesProto(b *testing.B) {
	dummyCount = 0
	root, err := services.UnmarshalProtoFile(ProtoCountryFilePath)
	if err != nil {
		b.Fatal(err)
	}
	// Build a map from CIDR to ISO code for O(1) lookups
	cidrToIso := make(map[string]string, len(root.CidrCountryPairs))
	for cidr, index := range root.CidrCountryPairs {
		var iso string
		if root.Geos[index].Country != nil {
			iso = root.Geos[index].Country.IsoCode
		} else if root.Geos[index].RegisteredCountry != nil {
			iso = root.Geos[index].RegisteredCountry.IsoCode
		}
		cidrToIso[cidr] = iso
	}

	// Pre-store the CIDR strings from testTable
	cidrs := make([]string, len(testTable))
	for i, tc := range testTable {
		cidrs[i] = tc.CIDR
	}

	// Reset the timer before the benchmark loop
	for i := 0; b.Loop(); i++ {
		cidr := cidrs[i%len(cidrs)]
		iso, ok := cidrToIso[cidr]
		if !ok {
			b.Fatalf("CIDR %v not found in proto data", cidr)
		}
		dummyCount += len(iso)
	}
}
