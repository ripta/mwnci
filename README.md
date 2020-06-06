mwnci: an implementation of the Monkey programming language,
       which is designed by [Thorsten Ball](https://interpreterbook.com/) 

To install:

```
go install github.com/ripta/mwnci/cmd/mwnci@latest
```

To run REPL, just invoke `mwnci`; the REPL does not support multiline input,
e.g., function definitions have to be compressed into one line. The evaluation
result for each line is printed.

To run a file, pass it one or more files, e.g.:

```
❯ mwnci examples/math.mwn
46
```

When running in file mode, per-line results are suppressed. Multiple files are
run in separate environments; no sharing across source files yet.

Some differences compared to the reference implementation:

- file mode in addition to REPL mode
- multibyte runes
- 64-bit floats
- line numbers in errors (incomplete)

