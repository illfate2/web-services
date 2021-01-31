# DNS-server
 
## Description
Implementation of DNS server with support of CNAME, A and NS types.

## Run

Build locally:
```make build
./bin/dns --help
```

Run with test config:
```
./bin/dns --config=config-example.csv
```
Run in debug mode:
```
./bin/dns --config=config-example.csv --debug=true
```

## Test

To test, build and run with default config:
`./bin/dns --config=config-example.csv`

Use `nslookup` or other cli tool to test:
```bash
nslookup -port=8090 -type=A domain.con localhost
nslookup -port=8090 -type=CNAME sub.domain.com localhost
nslookup -port=8090 -type=NS random.com localhost
```
