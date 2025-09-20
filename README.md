# urlcheck

Fast concurrent URL checker written in Go.

## Installation

```
go build -o urlcheck cmd/urlcheck/main.go
```


## Usage
### Check URLs from command line
```
./urlcheck -urls "https://google.com,https://github.com"
```

### Check URLs from file
```
./urlcheck -file urls.txt
```

### Read from stdin
```
cat urls.txt | ./urlcheck -stdin
```

### Options
- urls string Comma-separated URLs
- file string File with URLs (one per line)
- stdin Read from stdin
- workers int Concurrent workers (default: 10)
- timeout duration Request timeout (default: 5s)
- quiet Show errors only
- color Colored output (default: true)

## Example Output
```
[1/2] ✓ https://google.com (200, 145ms)

[2/2] ✗ https://badurl.com (DNS failed)

Summary: 1 successful, 1 failed, 50.0% success rate
```

## License

[MIT](LICENSE)
