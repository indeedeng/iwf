# get rid of default behaviors, they're just noise
MAKEFLAGS += --no-builtin-rules
.SUFFIXES:

default: help

# ###########################################
#                TL;DR DOCS:
# ###########################################
# - Targets should never, EVER be *actual source files*.
#   Always use book-keeping files in $(BUILD).
#   Otherwise e.g. changing git branches could confuse Make about what it needs to do.
# - Similarly, prerequisites should be those book-keeping files,
#   not source files that are prerequisites for book-keeping.
#   e.g. depend on .build/fmt, not $(ALL_SRC), and not both.
# - Be strict and explicit about prerequisites / order of execution / etc.
# - Test your changes with `-j 27 --output-sync` or something!
# - Test your changes with `make -d ...`!  It should be reasonable!

# temporary build products and book-keeping targets that are always good to / safe to clean.
BUILD := .build
# less-than-temporary build products, e.g. tools.
# usually unnecessary to clean, and may require downloads to restore, so this folder is not automatically cleaned.
BIN := .bin

# ====================================
# book-keeping files that are used to control sequencing.
#
# you should use these as prerequisites in almost all cases, not the source files themselves.
# these are defined in roughly the reverse order that they are executed, for easier reading.
#
# recipes and any other prerequisites are defined only once, further below.
# ====================================

# ====================================
# helper vars
# ====================================

# set a VERBOSE=1 env var for verbose output. VERBOSE=0 (or unset) disables.
# this is used to make verbose flags, suitable for `$(if $(test_v),...)`.
VERBOSE ?= 0
ifneq (0,$(VERBOSE))
test_v = 1
else
test_v =
endif

# a literal space value, for makefile purposes
SPACE :=
SPACE +=
COMMA := ,

# M1 macs may need to switch back to x86, until arm releases are available
EMULATE_X86 =
ifeq ($(shell uname -sm),Darwin arm64)
EMULATE_X86 = arch -x86_64
endif

# helper for executing bins that need other bins, just `$(BIN_PATH) the_command ...`
# I'd recommend not exporting this in general, to reduce the chance of accidentally using non-versioned tools.
BIN_PATH := PATH="$(abspath $(BIN)):$$PATH"

# version, git sha, etc flags.
# reasonable to make a :=, but it's only used in one place, so just leave it lazy or do it inline.
GO_BUILD_LDFLAGS = $(shell ./scripts/go-build-ldflags.sh LDFLAG)

# automatically gather all source files that currently exist.
# works by ignoring everything in the parens (and does not descend into matching folders) due to `-prune`,
# and everything else goes to the other side of the `-o` branch, which is `-print`ed.
# this is dramatically faster than a `find . | grep -v vendor` pipeline, and scales far better.
FRESH_ALL_SRC = $(shell \
	find . \
	\( \
		-path './vendor/*' \
		-o -path './idls/*' \
		-o -path './.build/*' \
		-o -path './.bin/*' \
	\) \
	-prune \
	-o -name '*.go' -print \
)
# most things can use a cached copy, e.g. all dependencies.
# this will not include any files that are created during a `make` run, e.g. via protoc,
# but that generally should not matter (e.g. dependencies are computed at parse time, so it
# won't affect behavior either way - choose the fast option).
#
# if you require a fully up-to-date list, e.g. for shell commands, use FRESH_ALL_SRC instead.
ALL_SRC := $(FRESH_ALL_SRC)
# as lint ignores generated code, it can use the cached copy in all cases
LINT_SRC := $(filter-out %_test.go ./.gen/%, $(ALL_SRC))

# ====================================
# $(BIN) targets
# ====================================

# downloads and builds a go-gettable tool, versioned by go.mod, and installs
# it into the build folder, named the same as the last portion of the URL.
define go_build_tool
@echo "building $(notdir $(1)) from $(1)..."
@go build -mod=readonly -o $(BIN)/$(notdir $(1)) $(1)
endef

# ====================================
# developer-oriented targets
#
# many of these share logic with other intermediates, but are useful to make .PHONY for output on demand.
# as the Makefile is fast, it's reasonable to just delete the book-keeping file recursively make.
# this way the effort is shared with future `make` runs.
# ====================================

# "re-make" a target by deleting and re-building book-keeping target(s).
# the + is necessary for parallelism flags to be propagated
define remake
@rm -f $(addprefix $(BUILD)/,$(1))
@+$(MAKE) --no-print-directory $(addprefix $(BUILD)/,$(1))
endef

# useful to actually re-run to get output again.
# reuse the intermediates for simplicity and consistency.
lint: ## (re)run the linter
	$(call remake,proto-lint lint)

# intentionally not re-making, goimports is slow and it's clear when it's unnecessary
fmt: $(BUILD)/fmt ## run goimports

# ====================================
# binaries to build
# ====================================
GOOS ?= $(shell go env GOOS)
GOARCH ?= $(shell go env GOARCH)

BINS =

BINS += iwf-server
iwf-server:
	@echo "compiling iwf-server with OS: $(GOOS), ARCH: $(GOARCH)"
	@go build -o $@ cmd/server/main.go

.PHONY: bins release clean

idl-code-gen: #generate/refresh go clent code for idl, do this after update the idl file
	rm -Rf ./gen ; true
	openapi-generator generate -i iwf-idl/iwf.yaml -g go -o gen/iwfidl/ -p packageName=iwfidl -p generateInterfaces=true -p isGoSubmodule=false --git-user-id indeedeng --git-repo-id iwf-idl
	rm ./gen/iwfidl/go.* ; rm -rf ./gen/iwfidl/test; true

bins: $(BINS)

clean: ## Clean binaries and build folder
	rm -f $(BINS)
	rm -Rf $(BUILD)
	$(if \
		$(filter $(BIN)/fake-codegen, $(wildcard $(BIN)/*)), \
		$(warning fake build tools may exist, delete the $(BIN) folder to get real ones if desired),)

deps: ## Check for dependency updates, for things that are directly imported
	@make --no-print-directory DEPS_FILTER='$(JQ_DEPS_ONLY_DIRECT)' deps-all

deps-all: ## Check for all dependency updates
	@go list -u -m -json all \
		| $(JQ_DEPS_AGE) \
		| sort -n

cleanTestCache:
	$Q go clean -testcache

integTests:
	$Q go test -v ./integ

ci:
	$Q go test -v ./integ -search=false -cadence=false -temporalHostPort=temporal:7233

integTestsNoSearch:
	$Q go test -v ./integ -search=false

stressTestsWithSearch:
	$Q go test -v ./integ -repeat=10 -intervalMs=100 -searchWaitMs=100 | tee test.log # TODO https://github.com/indeedeng/iwf/issues/134

stressTestsNoSearch:
	$Q go test -v ./integ -repeat=10 -search=false | tee test.log

stressTestsCadenceNoSearch:
	$Q go test -v ./integ -repeat=10 -temporal=false -search=false | tee test.log

stressTestsTemporalNoSearch:
	$Q go test -v ./integ -repeat=10 -cadence=false -search=false | tee test.log

help:
	@# print help first, so it's visible
	@printf "\033[36m%-20s\033[0m %s\n" 'help' 'Prints a help message showing any specially-commented targets'
	@# then everything matching "target: ## magic comments"
	@cat $(MAKEFILE_LIST) | grep -e "^[a-zA-Z_\-]*:.* ## .*" | awk 'BEGIN {FS = ":.*? ## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}' | sort
