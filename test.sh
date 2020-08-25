#! /bin/bash -eu

GOPATH="$(pwd -P)/../../../.."

TESTDIRS=$(find . -name "*_test.go" -and ! -path "./src/*" | xargs dirname | sort | uniq)

for i in $TESTDIRS; do
    echo Testing in $i:
    (cd $i; go test -v)
done
