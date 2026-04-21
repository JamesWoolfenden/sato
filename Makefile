.PHONY:
TEST?=$$(go list ./... | grep -v 'vendor'| grep -v 'scripts'| grep -v 'version')
HOSTNAME=jameswoolfenden
FULL_PKG_NAME=github.com/jameswoolfenden/sato
VERSION_PLACEHOLDER=version.ProviderVersion
# Repository owner/organization
HOSTNAME=jameswoolfenden
# Output binary name
BINARY=sato
# Version - use git tag if not set, default to dev
VERSION?=$(shell git describe --tags 2>/dev/null || echo "dev")
OS_ARCH=darwin_amd64
# Build flags for version injection
LDFLAGS=-ldflags "-X sato/src/version.Version=$(VERSION)"
TERRAFORM=./terraform/
TF_TEST=./terraform_test/

default:

build:
	go build $(LDFLAGS) -o ${BINARY}

release:
	GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o ./bin/${BINARY}_${VERSION}_darwin_amd64
	GOOS=freebsd GOARCH=386 go build $(LDFLAGS) -o ./bin/${BINARY}_${VERSION}_freebsd_386
	GOOS=freebsd GOARCH=amd64 go build $(LDFLAGS) -o ./bin/${BINARY}_${VERSION}_freebsd_amd64
	GOOS=freebsd GOARCH=arm go build $(LDFLAGS) -o ./bin/${BINARY}_${VERSION}_freebsd_arm
	GOOS=linux GOARCH=386 go build $(LDFLAGS) -o ./bin/${BINARY}_${VERSION}_linux_386
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o ./bin/${BINARY}_${VERSION}_linux_amd64
	GOOS=linux GOARCH=arm go build $(LDFLAGS) -o ./bin/${BINARY}_${VERSION}_linux_arm
	GOOS=openbsd GOARCH=386 go build $(LDFLAGS) -o ./bin/${BINARY}_${VERSION}_openbsd_386
	GOOS=openbsd GOARCH=amd64 go build $(LDFLAGS) -o ./bin/${BINARY}_${VERSION}_openbsd_amd64
	GOOS=solaris GOARCH=amd64 go build $(LDFLAGS) -o ./bin/${BINARY}_${VERSION}_solaris_amd64
	GOOS=windows GOARCH=386 go build $(LDFLAGS) -o ./bin/${BINARY}_${VERSION}_windows_386
	GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o ./bin/${BINARY}_${VERSION}_windows_amd64

test:
	go test $(TEST) || exit 1
	echo $(TEST) | xargs -t -n4 go test $(TESTARGS) -timeout=30s -parallel=4

testacc:
	TF_ACC=1 go test $(TEST) -v $(TESTARGS) -timeout 120m

coverage: ## Generate coverage report (excluding main.go)
	go test ./src/... -coverprofile=cover.out -covermode=atomic
	go tool cover -func=cover.out | grep total | awk '{print "Total coverage: " $$3}'

coverage-html: ## Generate HTML coverage report (excluding main.go)
	go test ./src/... -coverprofile=cover.out -covermode=atomic
	go tool cover -html=cover.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

coverage-all: ## Generate coverage report (including all packages)
	go test ./... -coverprofile=cover.out -covermode=atomic
	go tool cover -func=cover.out | grep total | awk '{print "Total coverage: " $$3}'


destroy:
	cd $(TERRAFORM) && terraform destroy --auto-approve


BIN=$(CURDIR)/bin
$(BIN)/%:
	@echo "Installing tools from tools/tools.go"
	@cat tools/tools.go | grep _ | awk -F '"' '{print $$2}' | GOBIN=$(BIN) xargs -tI {} go install {}

generate-docs:
	echo "does nowt"

docs:


vet:
	go vet ./...

bump:
	git push
	$(eval VERSION=$(shell git describe --tags --abbrev=0 | awk -F. '{OFS="."; $$NF+=1; print $0}'))
	git tag -a $(VERSION) -m "new release"
	git push origin $(VERSION)

psbump:
	git push
	powershell -command "./bump.ps1"

update:
	go get -u
	go mod tidy
	go mod vendor
	pre-commit autoupdate

lint:
	golangci-lint run --fix

gci:
	gci -w .

fmt:
	gofumpt -l -w .

staticcheck: ## Run staticcheck
	@./scripts/run-staticcheck.sh

gosec: ## Run security scanner
	gosec -quiet -exclude-generated ./...

govulncheck: ## Check for known vulnerabilities
	govulncheck ./...

complexity: ## Check cyclomatic complexity
	gocyclo -over 15 -avg .

check-all: ## Run all checks (vet, staticcheck, gosec, govulncheck)
	@echo "Running all code checks..."
	@$(MAKE) vet
	@$(MAKE) staticcheck
	@$(MAKE) gosec
	@$(MAKE) govulncheck
	@echo "All checks passed!"

security-scan: ## Run security-focused checks
	@echo "Running security scans..."
	@$(MAKE) gosec
	@$(MAKE) govulncheck
	@echo "Security scan complete!"

.PHONY: schema
schema: ## Download latest CloudFormation schemas and regenerate resource mappings
	wget -qO /tmp/cloudformation-schema.zip https://schema.cloudformation.us-east-1.amazonaws.com/CloudformationSchema.zip && unzip -o /tmp/cloudformation-schema.zip -d ./schema && rm /tmp/cloudformation-schema.zip
	@$(MAKE) mappings

.PHONY: mappings
mappings: ## Regenerate src/see/resource_mapping.go from schema/
	go run ./src/see/gen
