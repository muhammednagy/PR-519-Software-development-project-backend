PLATFORM?=linux
ARCHITECTURE?=amd64

BUILDTIME=`date "+%F %T%Z"`
VERSION=`git describe --tags`

build:
	GOOS=$(PLATFORM) GOARCH=$(ARCHITECTURE) go build -ldflags="-X 'github.com/muhammednagy/PR-519-Software-development-project-backend/api.buildTime=$(BUILDTIME)' -X 'github.com/muhammednagy/PR-519-Software-development-project-backend/api.version=$(VERSION)' -s -w -extldflags '-static'" -o bin/backend_api main.go

run:
	go run $(FILES)

clean:
	rm -rf bin
