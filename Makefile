VERSION = $(shell git tag --points-at HEAD)
COMMIT = $(shell git rev-parse --short HEAD)

# Package whose vars are stamped at link time.
BUILDINFO = github.com/5-bare-bones/5bb__sphinx/internal/buildinfo
# Base64 Ed25519 public key used to verify `ascend` manifests.
# Override per build: make build-scholar SPHINX_PUBKEY=...
SPHINX_PUBKEY ?=

# Linker flags shared by every tier build: strip symbols, stamp version + pubkey.
# Each tier target appends its own -X ...Tier=<rank>.
TIER_LD = -s -w \
	-X $(BUILDINFO).Version=$(VERSION) \
	-X $(BUILDINFO).PubKeyB64=$(SPHINX_PUBKEY)

# Default dev build: tier 1 (apprentice), unstripped tier stamp left at default.
build:
	@go build -o sphinx_dev -ldflags="-s -w" .

install:
	@go install -ldflags="-s -w" .

# --- Tier builds -----------------------------------------------------------
# Build tags are CUMULATIVE: each tier passes every lower tier's tag too.
# tier 6 (hero) and 7 (tutor) layer onto the master base.

build-apprentice:
	@go build -tags '' -ldflags "$(TIER_LD) -X $(BUILDINFO).Tier=apprentice" -o sphinx_apprentice .

build-adept:
	@go build -tags 'adept' -ldflags "$(TIER_LD) -X $(BUILDINFO).Tier=adept" -o sphinx_adept .

build-scholar:
	@go build -tags 'adept scholar' -ldflags "$(TIER_LD) -X $(BUILDINFO).Tier=scholar" -o sphinx_scholar .

build-keeper:
	@go build -tags 'adept scholar keeper' -ldflags "$(TIER_LD) -X $(BUILDINFO).Tier=keeper" -o sphinx_keeper .

build-master:
	@go build -tags 'adept scholar keeper master' -ldflags "$(TIER_LD) -X $(BUILDINFO).Tier=master" -o sphinx_master .

build-dev:
	@go build -tags 'adept scholar keeper master dev' -ldflags "$(TIER_LD) -X $(BUILDINFO).Tier=hero" -o sphinx_dev_hero .

build-tutor:
	@go build -tags 'adept scholar keeper master tutor' -ldflags "$(TIER_LD) -X $(BUILDINFO).Tier=tutor" -o sphinx_tutor .

build-all: build-apprentice build-adept build-scholar build-keeper build-master build-dev build-tutor

# Assert lower tiers ship strictly less code. The script builds its own
# unstripped binaries (symbol inspection needs the symbols -s -w would strip).
verify-tiers:
	@./scripts/verify_tiers.sh

# --- Tests -----------------------------------------------------------------
test:
	go test ./...

test-race:
	go test ./... -race

# Tests under the full tag set (exercises every tier's command packages).
test-all-tags:
	go test -tags 'adept scholar keeper master dev tutor' ./... -race

proto:
	@cd pb && for type in card entry file totp ; do \
		protoc -I. --go_out=. $$type.proto ; \
	done

docker-build:
	docker build -t sphinx .

docker-run:
	docker run -it --rm sphinx sh

cmds:
	@cd docs && go build main.go && ./main --cmd all

summary:
	@cd docs && go build main.go && ./main --summary

.PHONY: build install build-apprentice build-adept build-scholar build-keeper \
	build-master build-dev build-tutor build-all verify-tiers test test-race \
	test-all-tags proto docker-build docker-run cmds summary
