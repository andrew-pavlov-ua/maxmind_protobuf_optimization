# Geo Data Structures Comparison

This repository is designed to compare two different approaches for representing geographic data: the MMDB (MaxMind Database) format and Protobuf (proto) schemas. It evaluates both the file sizes and the lookup performance of these formats. Additionally, two versions of proto schemas are provided:

- **Non-Optimised Proto Schema**
- **Optimised Proto Schema**

## Repository Contents

- **Test Files:**  
  - JSON files containing geo data.
  - MMDB files used for IP-based geographic lookups.
  - Proto files (both non-optimised and optimised versions).

- **Benchmark Tests:**  
  Go tests and benchmarks to measure and compare the performance of lookups using MMDB and Proto representations.

## File Size Comparison

### Non-Optimised Version
- `GeoLite2-Country-Test.mmdb` — 20K
- `GeoLite2-Country-Test.proto` — 84K
- `GeoLite2-Country-Test.json` — 328K



### After Optimisation
- `GeoLite2-Country-Test.proto` — 20K (Optimised proto schema)
  `GeoLite2-Country-Test.mmdb` — 20K
- `GeoLite2-Country-Test.json` — 328K

## Benchmark Results

The benchmarks compare the performance of country lookups using MMDB and Proto approaches.

### Before Optimisation

```bash
BenchmarkLookUpCountriesMmdb-12          1741753               706.8 ns/op            34 B/op          2 allocs/op
BenchmarkLookUpCountriesProto-12         8568613               139.8 ns/op             0 B/op          0 allocs/op
```
### After Optimisation

```bash
BenchmarkLookUpCountriesMmdb-12          1847648               645.6 ns/op            33 B/op          1 allocs/op
BenchmarkLookUpCountriesProto-12        97627989                11.81 ns/op            0 B/op          0 allocs/op
