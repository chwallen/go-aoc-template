SHELL := /bin/sh
GO ?= go
GO_FILES := $(shell find . -name "*.go" ! -name "day_parts.go")
DAY_PARTS := cmd/run/day_parts.go
DAY_PARTS_TEMPLATE := internal/template/day_parts.tmpl

ifdef YEAR
YEAR_FLAG := --year $(YEAR)
TEST_FLAG := ./$(YEAR)
else
TEST_FLAG := ./*
endif

ifdef DAY
DAY_FLAG := --day $(DAY)
DAY_PADDED := $(shell echo "$(DAY)" | sed 's/^[1-9]$$/0&/')
TEST_FLAG := $(TEST_FLAG)/day$(DAY_PADDED)/
else
TEST_FLAG := $(TEST_FLAG)/day*/
endif

ifdef PART
PART_FLAG := --part $(PART)
endif

ifdef AOC_COOKIE
AOC_COOKIE_FLAG := --cookie $(AOC_COOKIE)
endif

.PHONY: run runall new description input test clean
.DEFAULT_GOAL := run

$(DAY_PARTS): $(GO_FILES) $(DAY_PARTS_TEMPLATE)
	@$(GO) generate ./cmd/run

run: $(DAY_PARTS)
	@$(GO) run ./cmd/run $(YEAR_FLAG) $(DAY_FLAG) $(PART_FLAG)

runall: $(DAY_PARTS)
	@$(GO) run ./cmd/run --all $(YEAR_FLAG)

new:
	@$(GO) run ./cmd/new $(YEAR_FLAG) $(DAY_FLAG) $(AOC_COOKIE_FLAG)

description:
	@$(GO) run ./cmd/description $(YEAR_FLAG) $(DAY_FLAG) $(AOC_COOKIE_FLAG)

input:
	@$(GO) run ./cmd/input $(YEAR_FLAG) $(DAY_FLAG) $(AOC_COOKIE_FLAG)

test:
	@$(GO) test --cover $(TEST_FLAG)

clean:
	$(RM) advent-of-code ./cmd/run/day_parts.go
