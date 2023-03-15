VERSION ?= latest
IMGDEV = leocbs/httpmiddleware:$(VERSION)
RUN=docker run --rm $(IMGDEV)
RUNTI=docker run -ti --rm $(IMGDEV)
cov=coverage.out
covhtml=coverage.html

imagedev:
	docker build --target devimage . -t $(IMGDEV)

release:
	git tag -a $(VERSION) -m "Generated release "$(VERSION)
	git push origin $(VERSION)

check: imagedev
	$(RUN) go test -tags=unit -timeout 60s -race -coverprofile=$(cov) ./...

coverage: check
	$(RUN) go tool cover -html=$(cov) -o=$(covhtml)
	xdg-open coverage.html

static-analysis: imagedev
	$(RUN) golangci-lint run ./...

modtidy:
	go mod tidy

fmt: imagedev
	$(RUN) gofmt -w -s -l .

shell: imagedev
	$(RUNTI) sh
