NAME			= nebulo
MAIN			= app.go

BUILD_DIR		= ./build
RELEASE_DIR		= ./release
TMP_DIR			= ./.tmp

COLOR_SUCCESS	= \e[0;32m
COLOR_FAIL		= \e[0;31m
COLOR_PRINT		= \e[0;33m
COLOR_RESET		= \e[0m

all : $(NAME)

$(NAME):
	@echo -e '$(COLOR_PRINT)Building $(NAME)...$(COLOR_RESET)'
	@echo -e '$(COLOR_PRINT)List of files beeing compiled:$(COLOR_RESET)'
	@go list -f '{{.GoFiles}}' ./...
	@mkdir -p $(BUILD_DIR)/bin
	@sed -i "s/\tversion = \".*\"/\tversion = \"$(shell [ $(shell git tag) ] && git tag || git rev-parse --verify HEAD)\"/" handler/version.go
	@go build -o $(BUILD_DIR)/bin/$(NAME) $(MAIN)
	@sed -i "s/\tversion = \".*\"/\tversion = \"\"/" handler/version.go
	@echo -e '$(COLOR_SUCCESS)Compilation done without errors$(COLOR_RESET)'

run: $(NAME)
	@./$(BUILD_DIR)/bin/$(NAME)

build: $(NAME)

clean:
	@echo -e '$(COLOR_PRINT)Cleaning...$(COLOR_RESET)'
	@rm -rf $(BUILD_DIR) $(RELEASE_DIR) $(TMP_DIR)
	@echo -e '$(COLOR_SUCCESS)Cleaned$(COLOR_RESET)'

doc_api:
	@echo -e '$(COLOR_PRINT)Generating apidoc...$(COLOR_RESET)'
	@mkdir -p $(BUILD_DIR)/doc/apidoc
	@apidoc -i ./ -o $(BUILD_DIR)/doc/apidoc/ -f ".*\\.go$$"
	@echo -e '$(COLOR_SUCCESS)Generated$(COLOR_RESET)'

doc_go:
	@echo -e '$(COLOR_PRINT)Generating godoc...$(COLOR_RESET)'
	@mkdir -p $(BUILD_DIR)/doc/godoc
	@./doc/godoc/save.sh
	@echo -e '$(COLOR_SUCCESS)Generated$(COLOR_RESET)'

doc_go_www:
	@echo -e '$(COLOR_PRINT)Open a webbrowser and go on 127.0.0.1:6060 ...$(COLOR_RESET)'
	@godoc -http=:6060 -index

doc: doc_api doc_go

test_dependencies:
	@echo -e '$(COLOR_PRINT)Testing dependencies...$(COLOR_RESET)'
	@govendor list +unused +missing
	@[ "$(shell govendor list +unused +missing | wc -l)" = "0" ]
	@echo -e '$(COLOR_SUCCESS)Done$(COLOR_RESET)'

test_code:
	@echo -e '$(COLOR_PRINT)Testing code with linters...$(COLOR_RESET)'
	@gofmt -d .
	@[ $(shell gofmt -d . | wc -l) = 0 ]
	@gometalinter --config=.gometalinter.json ./...
	@echo -e '$(COLOR_SUCCESS)Done$(COLOR_RESET)'

test_unit:
	@echo -e '$(COLOR_PRINT)Testing code with unit tests...$(COLOR_RESET)'
	@go test -v -timeout 5
	@echo -e '$(COLOR_SUCCESS)Done$(COLOR_RESET)'

test: test_dependencies test_code test_unit

release: test $(NAME) doc
	([ -z "$(TAG+x)" ] && echo "Please add the release tag with the TAG=x.x.x environment variable" && false) || true
	@mkdir -p $(RELEASE_DIR)
	rm -rf $(TMP_DIR)
	@mkdir -p $(TMP_DIR)/$(NAME)
	cp -r $(BUILD_DIR)/bin $(BUILD_DIR)/doc $(TMP_DIR)/$(NAME)
	@tar -cvzf $(RELEASE_DIR)/nebulo-full-$(TAG).tar.gz $(TMP_DIR)/
	@zip $(RELEASE_DIR)/nebulo-full-$(TAG).zip $(TMP_DIR)/
	@tar -cvzf $(RELEASE_DIR)/nebulo-doc-$(TAG).tar.gz $(TMP_DIR)/$(NAME)/doc/
	@zip $(RELEASE_DIR)/nebulo-doc-$(TAG).zip $(TMP_DIR)/$(NAME)/doc/
	@tar -cvzf $(RELEASE_DIR)/nebulo-bin-$(TAG).tar.gz $(TMP_DIR)/$(NAME)/bin/
	@zip $(RELEASE_DIR)/nebulo-bin-$(TAG).zip $(TMP_DIR)/$(NAME)/bin/
	rm -rf $(TMP_DIR)


.PHONY: $(NAME) run build clean doc_api doc_go doc_go_www doc test_dependencies test_code test_unit test release
