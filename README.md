`bsixfour` is a simple [base64](https://en.wikipedia.org/wiki/Base64) encoder/decoder.

When used from the command line, it reads input from stdin and prints to stdout. Its exported functions read from any `io.Reader` and return a `string` (for encoding) or a `[]byte` (for decoding).


## Usage

```bash
$ go build

$ echo -n 'hello, world!' | go ./bsixfour encode
aGVsbG8sIHdvcmxkIQ==

# or skip the build step
$ echo -n aGVsbG8sIHdvcmxkIQ== | go run . decode
hello, world!
```

Use `--append-newline=false` to prevent the newline character from being appended.
(This is useful for piping output, especially when decoding)


## Testing

Run unit tests with
```
$ go test
```
Add the `-v` flag to see each test.

## Benchmarks

Run benchmarks with
```
$ go test -bench .
```

Currently this takes multiple gigabytes of RAM, so you may consider running only encode or decode tests at a time:
```
$ go test -bench Encode && go test -bench Decode
```