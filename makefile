.DEFAULT_GOAL := help

GO_TEST_FLAGS?=-count=1 -p=4

# target: help - Display available recipes.
.PHONY: help
help:
	@egrep "^# target:" [Mm]akefile

# .PHONY: compose-build
# compose-build:
# 	@docker-compose -f docker-compose.yml build

# .PHONY: generate-mocks
# generate-mocks: compose-build
# 	@docker-compose -f docker-compose.yml run go mockery --all --inpackage

# # target: go-generate-mocks - Generate mocks.
# .PHONY: go-generate-mocks
# go-generate-mocks:
# 	@mockery --all --inpackage

# target: test-local - Run golang tests.
.PHONY: test
test:
	@go test $(GO_TEST_FLAGS) ./...

.PHONY: lint
lint:
	@yamllint -c .yamllint.yml

.PHONY: test-with-lint
test-with-lint: | test lint
