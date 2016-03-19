To get the fuzzing framework::

    go get -u github.com/dvyukov/go-fuzz/go-fuzz
    go get -u github.com/dvyukov/go-fuzz/go-fuzz-build

To run the fuzzing code::

    cd ~/go/src/github.com/seehuhn/classification/tree/fuzzing/
    go-fuzz-build github.com/seehuhn/classification/tree
    go-fuzz -bin=./tree-fuzz.zip -workdir=.
