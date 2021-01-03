VERSION := $(shell cat ./VERSION)


test:
	go test ./...

github_release:
	git tag -a v$(VERSION) -m "Release v$(VERSION)" || true
	git push origin v$(VERSION)

