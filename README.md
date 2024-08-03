# GoCalc
GoCalc is a simple calculator written in Go, mainly for the purpose of learning more about Go, and basic theoretical aspects of automata, and theory of computing.

### Building and Running


To build, simply `cd` into the top-level source directory, and run
```sh
$ go build
```

This should produce a binary called `main` in this directory.

To run the binary, with support for arrow keys in the prompt, download and install `rlwrap`, and run
```sh
$ rlwrap ./main
```
