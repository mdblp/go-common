#!/bin/sh
go install github.com/jstemmer/go-junit-report@latest
go install github.com/t-yuki/gocover-cobertura@latest

go test -v -race -cover -coverprofile=coverage.out ./... 2>&1 > test-result.txt
RET=$?
cat test-result.txt
cat test-result.txt | go-junit-report > test-report-v2.xml
cat test-report-v2.xml
rm -f /tmp/coverage-v2.html
go tool cover -html=coverage.out -o coverage-v2.html

exit $RET
