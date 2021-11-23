bin := git-add-commit
buildFolder := ./release
buildPath := $(buildFolder)/$(bin)

build:
	rm -rf build
	go build -o $(buildPath)

install:
	go build -o $(buildPath)
	mv $(buildPath) /usr/local/bin/

tar:
	tar -czf $(bin).tar.gz --directory=$(buildFolder) $(bin)
	shasum -a 256 $(bin).tar.gz

upload-assets:
	gh release upload $(version)  $(buildPath) $(bin).tar.gz
	gh release view $(version) --json assets -q ".assets[1].url"

publish: build tar tag upload-assets

tag:
	gh release create $(version) -t "v$(version)"

