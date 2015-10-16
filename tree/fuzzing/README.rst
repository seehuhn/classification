To get the fuzzing framework::

    go get github.com/dvyukov/go-fuzz/go-fuzz
    go get github.com/dvyukov/go-fuzz/go-fuzz-build

To run the fuzzing code::

    go-fuzz-build github.com/seehuhn/classification/tree
    go-fuzz -bin=./tree-fuzz.zip -workdir=.
