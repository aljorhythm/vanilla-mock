setup:
	git config core.hooksPath .githooks
fmt: setup
	sh .lint.sh || (echo "formatting failed $$?"; exit 1)
unit_test: fmt
	go test ./... || (echo "unit test failed $$?"; exit 1)
build: unit_test
	go build -o vanmock
install: build
	go install