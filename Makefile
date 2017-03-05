# Binary name to generate
BINARY_NAME			:= nebulo

# Used for generation and copy in archive
CONFIGURATION_FILE	:= config.sample.ini

# Used for go coverage tools
TEST_COVERAGE_MODE	?= count

# Overload this variable on make call 'make <function> CI=1' to add debug information
#	and remove terminal colors
CI					?= 0

# Used only on function 'release', generate one binary per couple os/arch
RELEASE_OS			?= linux windows darwin
RELEASE_ARCH		?= 386 amd64

# Temporary directories to use to generate binaries and documentation
DIR_PROJECT			:= $(dir $(abspath $(lastword $(MAKEFILE_LIST))))
DIR_BUILD			:= $(DIR_PROJECT)/build
DIR_COVERAGE		:= $(DIR_PROJECT)/coverage
DIR_RELEASE			:= $(DIR_PROJECT)/release
DIR_RELEASE_TMP		:= $(DIR_PROJECT)/.tmp/


# Used to fill the /version api endpoint
BUILD_VERSION		:= $(shell git describe --tags --always --dirty="-dev")
BUILD_TIME			:= $(shell date -u '+%Y-%m-%d-%H%M UTC')
BUILD_FLAGS			:= -ldflags='-X "main.BuildVersion=$(BUILD_VERSION)" -X "main.BuildTime=$(BUILD_TIME)"'

# CI variable define how informations are displayed on console
ifeq ($(CI),0)
	# Don't show command, and make everything 'pretty'
	Q				:= @
	COLOR_SUCCESS	:= \e[0;32m # green
	COLOR_FAIL		:= \e[0;31m # red
	COLOR_PRINT		:= \e[0;33m # orange
	COLOR_RESET		:= \e[0m
else
	# Show everything, make everything readable in a file
	Q				:= $(shell echo "")
	COLOR_SUCCESS	:= $Q
	COLOR_FAIL		:= $Q
	COLOR_PRINT		:= $Q
	COLOR_RESET		:= $Q
endif


all : $(BINARY_NAME)

# Compile for current os/arch and save binary in $DIR_BUILD folder
$(BINARY_NAME):
	$Q echo -e '$(COLOR_PRINT)Building $(BINARY_NAME)...$(COLOR_RESET)'
	$Q mkdir -p $(DIR_BUILD)/bin
	$Q go build -o $(DIR_BUILD)/bin/$(BINARY_NAME) $(BUILD_FLAGS)
	$Q echo -e '$(COLOR_SUCCESS)Compilation done without errors$(COLOR_RESET)'

build: $(BINARY_NAME)

config: $(BINARY_NAME)
	$Q [ "$(shell $(DIR_BUILD)/bin/$(BINARY_NAME) --config-dont-load-default --config-gen=$(CONFIGURATION_FILE) 2>&1 > /dev/null || echo $$? && false)" = "102" ]

# Compile for current os/arch and run binary
run: $(BINARY_NAME)
	$Q echo -e '$(COLOR_PRINT)Running $(BINARY_NAME):$(COLOR_RESET)'
	$Q $(DIR_BUILD)/bin/$(BINARY_NAME) ${ARGS}
	$Q echo -e '$(COLOR_PRINT)Terminated$(COLOR_RESET)'

# Remove all non-essentials directories and files
clean:
	$Q echo -e '$(COLOR_PRINT)Cleaning...$(COLOR_RESET)'
	$Q rm -rf $(DIR_BUILD) $(DIR_RELEASE) $(DIR_RELEASE_TMP)
	$Q echo -e '$(COLOR_SUCCESS)Cleaned$(COLOR_RESET)'

# Generate the API documentation
doc_api:
	$Q echo -e '$(COLOR_PRINT)Generating apidoc...$(COLOR_RESET)'
	$Q mkdir -p $(DIR_BUILD)/doc/apidoc
	$Q apidoc -i ./ -o $(DIR_BUILD)/doc/apidoc/ -f ".*\\.go$$"
	$Q echo -e '$(COLOR_SUCCESS)Generated$(COLOR_RESET)'

# Generate all kind of documentation
doc: doc_api

# Run code documentation server
godoc:
	$Q echo -e '$(COLOR_PRINT)Open a webbrowser and go on 127.0.0.1:6060 ...$(COLOR_RESET)'
	$Q godoc -http=:6060 -index
	$Q echo -e '$(COLOR_PRINT)Terminated$(COLOR_RESET)'

# Check for useless and missing dependencies
test_dependencies:
	$Q echo -e '$(COLOR_PRINT)Testing dependencies...$(COLOR_RESET)'
	$Q govendor list +unused +missing
	@[ "$(shell govendor list +unused +missing | wc -l)" = "0" ]
	$Q echo -e '$(COLOR_SUCCESS)Done$(COLOR_RESET)'

# Check syntax, format, useless, and non-optimized code
test_code:
	$Q echo -e '$(COLOR_PRINT)Testing code with linters...$(COLOR_RESET)'
	$Q gofmt -d .
	$Q [ $(shell gofmt -d . | wc -l) = 0 ]
	$Q gometalinter --config=.gometalinter.json ./...
	$Q echo -e '$(COLOR_SUCCESS)Done$(COLOR_RESET)'

# Check unit tests
test_unit:
	$Q echo -e '$(COLOR_PRINT)Testing code with unit tests...$(COLOR_RESET)'
	$Q go test -v -timeout 5s $(shell go list ./... | grep -v /vendor/)
	$Q echo -e '$(COLOR_SUCCESS)Done$(COLOR_RESET)'

# TODOs should not exist
test_todo:
	$Q echo -e '$(COLOR_PRINT)Testing presence of TODOs in code...$(COLOR_RESET)'
	$Q find . -name vendor -prune -o -name "*.go" -exec grep -Hn "//TODO:" {} \;
	@[ "$(shell find . -name vendor -prune -o -name "*.go" -exec grep -Hn "//TODO:" {} \; | wc -l)" = "0" ]
	$Q echo -e '$(COLOR_SUCCESS)Done$(COLOR_RESET)'

# Check all kind of tests
test: test_dependencies test_code test_unit test_todo

# Compute coverage and create coverage files
coverage:
	$Q rm -rf $(DIR_COVERAGE)
	$Q mkdir -p $(DIR_COVERAGE)
	$Q echo "mode: $(TEST_COVERAGE_MODE)" > $(DIR_COVERAGE)/coverage.out
	$Q for pkg in $(shell go list ./... | grep -v /vendor/); do \
		go test -covermode="$(TEST_COVERAGE_MODE)" -coverprofile="$(DIR_COVERAGE)/coverage.tmp" "$$pkg" 2>&1 > /dev/null; \
		grep -h -v "^mode:" $(DIR_COVERAGE)/coverage.tmp >> $(DIR_COVERAGE)/coverage.out 2> /dev/null; \
	done
	$Q rm -f $(DIR_COVERAGE)/coverage.tmp
	$Q go tool cover -func=$(DIR_COVERAGE)/coverage.out

coverage_show: coverage
	$Q go tool cover -html=$(DIR_COVERAGE)/coverage.out

# Compile and save all necessaries file for each couple os/arch
release_build: test coverage config
	$Q echo -e '$(COLOR_PRINT)List of files beeing compiled:$(COLOR_RESET)'
	$Q go list -f '{{.GoFiles}}' ./...
	$Q mkdir -p $(DIR_BUILD)/bin
	$Q echo -e '$(COLOR_PRINT)Building...$(COLOR_RESET)'
	$Q for os in $(RELEASE_OS); do \
		for arch in $(RELEASE_ARCH); do \
			echo "Building $(BINARY_NAME) for $$os/$$arch..."; \
			GOOS=$$os GOARCH=$$arch go build -o $(DIR_BUILD)/bin/$(BINARY_NAME)-$$os$$arch $(BUILD_FLAGS); \
		done; \
	done
	$Q echo -e '$(COLOR_SUCCESS)Compilation done without errors$(COLOR_RESET)'

# Generate a release ; multiple files are going to be generated:
#	a documentation archive, a runnable environment archive for each couple os/arch
release: config clean doc release_build
	$Q [ -n "$(TAG)" ] || (echo "Please add the release tag with the TAG=x.x.x environment variable" && false)
	$Q mkdir -p $(DIR_RELEASE) $(DIR_RELEASE_TMP)/$(BINARY_NAME)
	$Q cp $(DIR_BUILD)/bin/* $(DIR_RELEASE_TMP)/$(BINARY_NAME)
	$Q cp $(CONFIGURATION_FILE) $(DIR_RELEASE_TMP)/$(BINARY_NAME)/config.ini
	$Q tar --create --gzip --file=$(DIR_RELEASE)/$(BINARY_NAME)-doc-$(TAG).tar.gz -C $(DIR_BUILD)/ doc/
	$Q for os in $(RELEASE_OS); do \
		echo "Creating $(BINARY_NAME)-$$os archives..."; \
		cd $(DIR_RELEASE_TMP) && tar --verbose --create --gzip --file=$(DIR_RELEASE)/$(BINARY_NAME)-bin-$(TAG)-$$os.tar.gz \
			$(BINARY_NAME)/config.ini $(BINARY_NAME)/$(BINARY_NAME)-$$os*; \
	done


.PHONY: all $(BINARY_NAME) build config run clean doc_api doc godoc test_dependencies test_code test_unit test_todo test coverage coverage_show release_build release
