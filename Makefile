install:
	sudo chmod -R u+x ./scripts/bash
	@ ./scripts/bash/install.sh
	@ ./scripts/bash/prepare-githooks.sh

lint: 			# lint all the codebase
	golangci-lint run

lint-staged: 	# lint only the staged files
	golangci-lint run --new --new-from-rev=HEAD

run:
	export GOOGLE_APPLICATION_CREDENTIALS=assets/credentials/firebase.dev.json && \
	go run main.go

build:
	go build main.go

test:
	go test ./...