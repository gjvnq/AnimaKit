ANSI_RED="\033[0;31m"
ANSI_GREEN="\033[0;32m"
ANSI_BLUE="\033[0;34m"
ANSI_RESET="\033[0m"

ifneq ("$(wildcard /usr/local/opt/coreutils/libexec/gnubin/echo)","")
	ECHO="/usr/local/opt/coreutils/libexec/gnubin/echo"
else
	ECHO="/bin/echo"
endif

.PHONY: all test test-html docs dep-ensure

all: AnimaKit.a cmd/AnimaKit
docs:
	xdg-open "http://localhost:6060/pkg/github.com/gjvnq/AnimaKit/"
docs-server:
	godoc -http=:6060
test: coverage.out
dep-ensure:
	@$(ECHO) -e $(ANSI_GREEN)"["$@"] Ensuring dependencies..."$(ANSI_RESET)
	dep ensure
test-html: coverage.out
	@$(ECHO) -e $(ANSI_GREEN)"Generating coverage report..."$(ANSI_RESET)
	go tool cover -html=coverage.out
	@$(ECHO) -e $(ANSI_BLUE)"Finished target"$(ANSI_RESET)

bindata.go: res/*
	@$(ECHO) -e $(ANSI_GREEN)"["$@"] Packing bin-data..."$(ANSI_RESET)
	go-bindata -pkg AnimaKit res/
	
AnimaKit.a: bindata.go *.go dep-ensure
	@$(ECHO) -e $(ANSI_GREEN)"["$@"] Fixing imports..."$(ANSI_RESET)
	goimports -w .
	@$(ECHO) -e $(ANSI_GREEN)"["$@"] Formatting code..."$(ANSI_RESET)
	go fmt
	@$(ECHO) -e $(ANSI_GREEN)"["$@"] Compiling code..."$(ANSI_RESET)
	go build -o $@
	@$(ECHO) -e $(ANSI_BLUE)"["$@"] Finished target $@"$(ANSI_RESET)

coverage.out: *.go AnimaKit.a
	@$(ECHO) -e $(ANSI_GREEN)"["$@"] Testing code..."$(ANSI_RESET)
	go test -cover -coverprofile=coverage.out
	@$(ECHO) -e $(ANSI_BLUE)"["$@"] Finished target"$(ANSI_RESET)

cmd/AnimaKit: AnimaKit.a cmd/*.go
	cd cmd && make AnimaKit

clean:
	-rm cmd/AnimaKit AnimaKit.a