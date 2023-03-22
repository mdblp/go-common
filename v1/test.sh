#!/bin/sh
go install github.com/jstemmer/go-junit-report@latest
go install github.com/t-yuki/gocover-cobertura@latest

go test -v -race -cover -coverprofile=coverage.out ./... 2>&1 > test-result.txt
RET=$?
cat test-result.txt
cat test-result.txt | go-junit-report > test-report.xml
cat test-report.xml
rm -f /tmp/coverage.html
go tool cover -html=coverage.out -o coverage.html

exit $RET