.PHONY: all
all: build fmt vet lint test
LIB=go-rabbit
GLIDE_NOVENDOR=$(shell glide novendor)
ALL_PACKAGES=$(shell go list ./... | grep -v "vendor")
UNIT_TEST_PACKAGES=$(shell glide novendor | grep -v "featuretests")
APP_EXECUTABLE="./out/$(LIB)"

setup:
	go get -u github.com/golang/lint/golint
	glide install
	go get -u github.com/360EntSecGroup-Skylar/goreporter
build-deps:
	glide install
update-deps:
	glide update
compile:
	go build -o $(APP_EXECUTABLE)
build: compile fmt vet
update-build: update-deps build
install:
	go install $(GLIDE_NOVENDOR)
fmt:
	go fmt $(GLIDE_NOVENDOR)
vet:
	go vet $(GLIDE_NOVENDOR)
lint:
	@for p in $(UNIT_TEST_PACKAGES); do \
		echo "==> Linting $$p"; \
		golint $$p | { grep -vwE "exported (var|function|method|type|const) \S+ should have comment" || true; } \
	done
test: compile
	go test $(GLIDE_NOVENDOR)
test-coverage: compile
	mkdir -p ./out
	@echo "mode: count" > out/coverage-all.out
	$(foreach pkg, $(ALL_PACKAGES),\
	ENVIRONMENT=test go test -coverprofile=out/coverage.out -covermode=count $(pkg);\
	tail -n +2 out/coverage.out >> out/coverage-all.out;)
	go tool cover -html=out/coverage-all.out -o out/coverage.html