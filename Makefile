.PHONY: test test-unit test-integration coverage-gaps coverage-gaps-check docs-generate docs-check docs-serve

test:
	go test ./...

test-unit:
	go test ./coverage/coveragegaps ./pkg/...
	go test ./test/cli -run 'Test(Leaf|Version|Config|Do|Load)'

test-integration:
	go test ./test/cli -run 'TestGroupCommandsWithMockedAPI|TestListDevicesOutputModes|TestIntegrationFailsWithoutMockServer'

coverage-gaps:
	go run ./coverage/coveragegaps \
		--json-out coverage/coverage-gaps.json \
		--md-out coverage/coverage-gaps.md

coverage-gaps-check:
	go run ./coverage/coveragegaps \
		--json-out coverage/coverage-gaps.json \
		--md-out coverage/coverage-gaps.md \
		--diff-out coverage/coverage-gaps-diff.md \
		--baseline coverage/coverage-gaps-baseline.json \
		--fail-on-regression \
		--fail-on-gaps

docs-generate:
	go run ./coverage/docsgen --mode generate

docs-check:
	go run ./coverage/docsgen --mode check

docs-serve:
	docsify serve docs
