# geoip

A CLI tool that maps IPv4 addresses to countries and outputs a sorted list of countries by IP count.

## Usage

```bash
geoip <ip_list_file>
```

The input file should contain one IPv4 address per line. Basic validation is performed on the input.

## Features

- Batch processing of IPs (100 IPs per request) using ip-api.com's free API
- Basic IP validation and format correction
- Sorted output of countries by IP count

## API Notes

Uses ip-api.com free batch API which:
- Allows up to 100 IPs per batch request
- Has a rate limit of 45 requests/minute
- Implements a 1.5s delay between batches to respect rate limits

For processing large IP lists (>1000 IPs), consider:
- Using a local GeoIP database
- Switching to ip-api.com or other providers' paid API

## Example Output

```bash
China: 86
United States: 70
South Korea: 44
India: 41
Brazil: 25
```

## Installation

```bash
go install github.com/LixenWraith/geoip@latest
```

Or build from source:

```bash
git clone github.com/LixenWraith/geoip
cd geoip
go build
```

Requires Go 1.23.5 or later.

## Use Case

Useful for analyzing geographical distribution of banned IPs from server logs.
```