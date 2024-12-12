.PHONY: proto-format proto-lint proto-gen format build
all: proto-all format build

###############################################################################
###                                  Build                                  ###
###############################################################################

build:
	@echo "🤖 Building simd..."
	@cd simapp && make build 1> /dev/null
	@echo "✅ Completed build!"

###############################################################################
###                                 Tooling                                 ###
###############################################################################

gofumpt_cmd=mvdan.cc/gofumpt
goimports_reviser_cmd=github.com/incu6us/goimports-reviser/v3

format:
	@echo "🤖 Running formatter..."
	@go run $(gofumpt_cmd) -l -w keeper/*
	@go run $(goimports_reviser_cmd) keeper/* &> /dev/null
	@echo "✅ Completed formatting!"

###############################################################################
###                                Protobuf                                 ###
###############################################################################

BUF_VERSION=1.42
BUILDER_VERSION=0.15.1

proto-all: proto-format proto-lint proto-gen

proto-format:
	@echo "🤖 Running protobuf formatter..."
	@docker run --rm --volume "$(PWD)":/workspace --workdir /workspace \
		bufbuild/buf:$(BUF_VERSION) format --diff --write
	@echo "✅ Completed protobuf formatting!"

proto-gen:
	@echo "🤖 Generating code from protobuf..."
	@docker run --rm --volume "$(PWD)":/workspace --workdir /workspace \
		ghcr.io/cosmos/proto-builder:$(BUILDER_VERSION) sh ./proto/generate.sh
	@echo "✅ Completed code generation!"

proto-lint:
	@echo "🤖 Running protobuf linter..."
	@docker run --rm --volume "$(PWD)":/workspace --workdir /workspace \
		bufbuild/buf:$(BUF_VERSION) lint
	@echo "✅ Completed protobuf linting!"
