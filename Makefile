VERSION := $(shell cat ./VERSION)

github_release:
	git tag -a $(VERSION) -m "Release v$(VERSION)" || true
	git push origin $(VERSION)

