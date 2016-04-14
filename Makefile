GIT_VER := $(shell git describe --tags)

.PHONY: all packages install clean

all:
	go build

install:
	go get github.com/fujiwara/go-dummer-simple

packages:
	gox -os="linux darwin" -arch="amd64 arm" -output "pkg/{{.Dir}}-${GIT_VER}-{{.OS}}-{{.Arch}}" -ldflags "-w -s"
	cd pkg && find . -name "*${GIT_VER}*" -type f -name "go-dummer-simple*" -exec zip {}.zip {} \;

clean:
	rm -f pkg/*
