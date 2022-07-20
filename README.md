`bsixfour` is a simple [base64](https://en.wikipedia.org/wiki/Base64) encoder/decoder.

It reads input from stdin but its exported functions read from any `io.Reader` and return a `string` (for encoding) or a `[]byte` (for decoding).


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