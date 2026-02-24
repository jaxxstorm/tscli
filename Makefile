.PHONY: test test-unit test-integration coverage-gaps coverage-gaps-check

test:
	go test ./...

test-unit:
	go test ./test/coveragegaps ./pkg/...
	go test ./cmd/tscli -run 'Test(Leaf|Version|Config|Do|Load)'

test-integration:
	go test ./cmd/tscli -run 'TestGroupCommandsWithMockedAPI|TestListDevicesOutputModes|TestIntegrationFailsWithoutMockServer'

coverage-gaps:
	go run ./test/coveragegaps \
		--json-out coverage/coverage-gaps.json \
		--md-out coverage/coverage-gaps.md

coverage-gaps-check:
	go run ./test/coveragegaps \
		--json-out coverage/coverage-gaps.json \
		--md-out coverage/coverage-gaps.md \
		--diff-out coverage/coverage-gaps-diff.md \
		--baseline coverage/coverage-gaps-baseline.json \
		--fail-on-regression
