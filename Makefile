.PHONY: test test-unit test-integration openapi-refresh coverage-gaps coverage-gaps-check coverage-gaps-latest docs-generate docs-check docs-serve

OPENAPI_SOURCE_URL ?= https://api.tailscale.com/api/v2?outputOpenapiSchema=true
OPENAPI_DIR ?= $(CURDIR)/pkg/contract/openapi
OPENAPI_SCHEMA ?= $(OPENAPI_DIR)/tailscale-v2-openapi.yaml
OPENAPI_METADATA ?= $(OPENAPI_DIR)/snapshot-metadata.yaml
OPENAPI_COMMAND_MAP ?= $(OPENAPI_DIR)/command-operation-map.yaml
LEAF_COMMANDS ?= $(CURDIR)/test/cli/testdata/leaf_commands.txt

COVERAGE_DIR ?= $(CURDIR)/coverage
COVERAGE_JSON ?= $(COVERAGE_DIR)/coverage-gaps.json
COVERAGE_MD ?= $(COVERAGE_DIR)/coverage-gaps.md
COVERAGE_DIFF ?= $(COVERAGE_DIR)/coverage-gaps-diff.md
COVERAGE_BASELINE ?= $(COVERAGE_DIR)/coverage-gaps-baseline.json
COVERAGE_EXCLUSIONS ?= $(COVERAGE_DIR)/exclusions.yaml
COVERAGE_PROPERTY_MANIFEST ?= $(COVERAGE_DIR)/property-coverage.yaml
COVERAGE_PROPERTY_EXCLUSIONS ?= $(COVERAGE_DIR)/property-exclusions.yaml
GOCACHE ?= $(CURDIR)/.gocache
GO_RUN = GOCACHE=$(GOCACHE) go run
GO_TEST = GOCACHE=$(GOCACHE) go test
TOOLS_GO_RUN = GOCACHE=$(GOCACHE) go -C tools run
TOOLS_GO_TEST = GOCACHE=$(GOCACHE) go -C tools test

test:
	$(GO_TEST) ./...
	$(TOOLS_GO_TEST) ./...

test-unit:
	$(GO_TEST) ./coverage/coveragegaps ./pkg/...
	$(GO_TEST) ./test/cli -run 'Test(Leaf|Version|Config|Do|Load)'

test-integration:
	$(GO_TEST) ./test/cli -run 'TestGroupCommandsWithMockedAPI|TestListDevicesOutputModes|TestIntegrationFailsWithoutMockServer'

openapi-refresh:
	$(TOOLS_GO_RUN) ./cmd/openapirefresh \
		--source-url $(OPENAPI_SOURCE_URL) \
		--schema-out $(OPENAPI_SCHEMA) \
		--metadata-out $(OPENAPI_METADATA)

coverage-gaps:
	$(TOOLS_GO_RUN) ./cmd/coveragegaps \
		--openapi $(OPENAPI_SCHEMA) \
		--mapping $(OPENAPI_COMMAND_MAP) \
		--manifest $(LEAF_COMMANDS) \
		--exclusions $(COVERAGE_EXCLUSIONS) \
		--property-coverage $(COVERAGE_PROPERTY_MANIFEST) \
		--property-exclusions $(COVERAGE_PROPERTY_EXCLUSIONS) \
		--json-out $(COVERAGE_JSON) \
		--md-out $(COVERAGE_MD)

coverage-gaps-check:
	$(TOOLS_GO_RUN) ./cmd/coveragegaps \
		--openapi $(OPENAPI_SCHEMA) \
		--mapping $(OPENAPI_COMMAND_MAP) \
		--manifest $(LEAF_COMMANDS) \
		--exclusions $(COVERAGE_EXCLUSIONS) \
		--property-coverage $(COVERAGE_PROPERTY_MANIFEST) \
		--property-exclusions $(COVERAGE_PROPERTY_EXCLUSIONS) \
		--json-out $(COVERAGE_JSON) \
		--md-out $(COVERAGE_MD) \
		--diff-out $(COVERAGE_DIFF) \
		--baseline $(COVERAGE_BASELINE) \
		--fail-on-regression \
		--fail-on-gaps

coverage-gaps-latest: openapi-refresh
	$(MAKE) coverage-gaps-check

docs-generate:
	$(GO_RUN) ./coverage/docsgen --mode generate

docs-check:
	$(GO_RUN) ./coverage/docsgen --mode check

docs-serve:
	docsify serve docs
