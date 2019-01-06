ANSI_RED="\033[0;31m"
ANSI_GREEN="\033[0;32m"
ANSI_BLUE="\033[0;34m"
ANSI_BOLD="\033[1m"
ANSI_RESET="\033[0m"

ifneq ("$(wildcard /usr/local/opt/coreutils/libexec/gnubin/echo)","")
	ECHO="/usr/local/opt/coreutils/libexec/gnubin/echo"
else
	ECHO="/bin/echo"
endif

.PHONY: all test test-html docs dep-ensure install help

help:
	@$(ECHO) -e "Available targets"
	@$(ECHO) -e "\t"$(ANSI_BOLD)"all"$(ANSI_RESET)
	@$(ECHO) -e "\t"$(ANSI_BOLD)"docs"$(ANSI_RESET)" - opens the docs in your browser"
	@$(ECHO) -e "\t"$(ANSI_BOLD)"docs-server"$(ANSI_RESET)" - starts godoc local web server"
	@$(ECHO) -e "\t"$(ANSI_BOLD)"dep-ensure"$(ANSI_RESET)" - ensures the dependencies are installed"
	@$(ECHO) -e "\t"$(ANSI_BOLD)"install"$(ANSI_RESET)" - installs the animakit's cli tool"

all: AnimaKit.a animakit
docs:
	xdg-open "http://localhost:6060/pkg/github.com/gjvnq/AnimaKit/"
docs-server:
	godoc -http=:6060
test: coverage.out
dep-ensure:
ifneq ($(FAST_MAKE),1)
	@$(ECHO) -e $(ANSI_GREEN)"["$@"] Ensuring dependencies..."$(ANSI_RESET)
	dep ensure
endif
test-html: coverage.out
	@$(ECHO) -e $(ANSI_GREEN)"Generating coverage report..."$(ANSI_RESET)
	go tool cover -html=coverage.out
	@$(ECHO) -e $(ANSI_BLUE)"Finished target"$(ANSI_RESET)

bindata.go: res/* dep-ensure
	@$(ECHO) -e $(ANSI_GREEN)"["$@"] Packing bin-data..."$(ANSI_RESET)
	go-bindata -pkg AnimaKit res/
	
AnimaKit.a: bindata.go *.go dep-ensure
ifneq ($(FAST_MAKE),1)
	@$(ECHO) -e $(ANSI_GREEN)"["$@"] Fixing imports..."$(ANSI_RESET)
	goimports -w .
endif
	@$(ECHO) -e $(ANSI_GREEN)"["$@"] Formatting code..."$(ANSI_RESET)
	go fmt
	@$(ECHO) -e $(ANSI_GREEN)"["$@"] Compiling code..."$(ANSI_RESET)
	go build -o $@
	@$(ECHO) -e $(ANSI_BLUE)"["$@"] Finished target $@"$(ANSI_RESET)

coverage.out: *.go AnimaKit.a
	@$(ECHO) -e $(ANSI_GREEN)"["$@"] Testing code..."$(ANSI_RESET)
	go test -cover -coverprofile=coverage.out
	@$(ECHO) -e $(ANSI_BLUE)"["$@"] Finished target"$(ANSI_RESET)

animakit: AnimaKit.a cmd/animakit/*.go
ifneq ($(FAST_MAKE),1)
	@$(ECHO) -e $(ANSI_GREEN)"["$@"] Fixing imports..."$(ANSI_RESET)
	cd cmd/animakit && goimports -w .
endif
	@$(ECHO) -e $(ANSI_GREEN)"["$@"] Formatting code..."$(ANSI_RESET)
	cd cmd/animakit && go fmt
	@$(ECHO) -e $(ANSI_GREEN)"["$@"] Compiling code..."$(ANSI_RESET)
	cd cmd/animakit && go build -o ../../$@
	@$(ECHO) -e $(ANSI_BLUE)"["$@"] Finished target $@"$(ANSI_RESET)

install: animakit
	@$(ECHO) -e $(ANSI_GREEN)"["$@"] Installing..."$(ANSI_RESET)
	cd cmd/animakit && go install
	@$(ECHO) -e $(ANSI_BLUE)"["$@"] Finished target $@"$(ANSI_RESET)

clean:
	-rm AnimaKit.a animakit