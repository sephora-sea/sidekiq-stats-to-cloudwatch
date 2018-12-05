.PHONY: build clean deploy

build:
	dep ensure -v
	env GOOS=linux go build -ldflags="-s -w" -o bin/main main/main.go

clean:
	rm -rf ./bin ./vendor Gopkg.lock

deploy: build
	/bin/bash -c 'serverless deploy --stage production'
