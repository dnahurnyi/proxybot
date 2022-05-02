### linting options below
### lint: for local environment

OSNAME = $(shell go env GOOS)
ARCH = $(shell go env GOARCH)
TMP_DIR = /tmp/dnahurnyi

LINT_VERSION = 1.45.0

LINT_DIR = $(TMP_DIR)/golangci-lint/$(LINT_VERSION)
LINT_BIN = $(LINT_DIR)/golangci-lint

CMD = golangci-lint run --allow-parallel-runners -c .golangci.yml
ifeq ($(ARCH), arm64)
	CMD = golangci-lint run --allow-parallel-runners -c .golangci.yml --build-tags dynamic
endif

LINT_ARCHIVE = golangci-lint-$(LINT_VERSION)-$(OSNAME)-$(ARCH).tar.gz
LINT_ARCHIVE_DEST = $(TMP_DIR)/$(LINT_ARCHIVE)

# Run this on local machine.
# It downloads a version of golangci-lint and execute it locally.
.PHONY: lint
lint: $(LINT_BIN)
	$(LINT_DIR)/$(CMD)

# install a local golangci-lint if not found.
$(LINT_BIN):
	curl -L --create-dirs \
		https://github.com/golangci/golangci-lint/releases/download/v$(LINT_VERSION)/$(LINT_ARCHIVE) \
		--output $(LINT_ARCHIVE_DEST)
	mkdir -p $(LINT_DIR)/
	tar -xf $(LINT_ARCHIVE_DEST) --strip-components=1 -C $(LINT_DIR)/
	chmod +x $(LINT_BIN)
	rm -f $(LINT_ARCHIVE_DEST)
