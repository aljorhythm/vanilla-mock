setup:
	git config core.hooksPath .githooks
	go get
fmt: setup
	sh .lint.sh || (echo "formatting failed $$?"; exit 1)
unit_test: fmt
	go test ./... || (echo "unit test failed $$?"; exit 1)
build:
	go build -o vanilla-mock
	cp vanilla-mock system-test
install:
	go install