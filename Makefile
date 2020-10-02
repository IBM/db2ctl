build:
	go build -ldflags "-X github.com/IBM/db2ctl/cmd.gitCommitHash=`git rev-parse HEAD` -X github.com/IBM/db2ctl/cmd.buildTime=`date -u '+%Y-%m-%d--%H:%M:%S%p'` -X github.com/IBM/db2ctl/cmd.gitBranch=`git branch --show-current` -X github.com/IBM/db2ctl/cmd.tagVersion=`git describe --tags --long`" -o bin/db2ctl main.go
build-w-clean: clean build
build-linux: # example: make build-linux DB_PATH=/dir/to/db
	env GOOS=linux GOARCH=amd64 go build -ldflags "-X github.com/IBM/db2ctl/internal/command.stateDBPathFromEnv=/tmp -X github.com/IBM/db2ctl/internal/command.logDirPathFromEnv=/var/log/IBM/db2ctl -X github.com/IBM/db2ctl/cmd.gitCommitHash=`git rev-parse HEAD` -X github.com/IBM/db2ctl/cmd.buildTime=`date -u '+%Y-%m-%d--%H:%M:%S%p'` -X github.com/IBM/db2ctl/cmd.gitBranch=`git branch --show-current` -X github.com/IBM/db2ctl/cmd.tagVersion=`git describe --tags --long`" -o bin/db2ctl main.go
send-fyre: install build-linux
	scp bin/db2ctl root@db2-pc-test-1-01.fyre.ibm.com:/usr/local/bin
send-fyre2: install build-linux
	scp bin/db2ctl root@p-17.fyre.ibm.com:/usr/local/bin
clean:
	rm -f bin/db2ctl
	rm -rf generated
	rm -rf logs
	rm -f *.db
	rm -f *.log
install: add-static clean
	go install
add-static: #add static code to binary. if error: do 'go get github.com/rakyll/statik'
	statik -src static -f
run-help:
	go run main.go --help
run-server: install
	db2ctl server
server-live: # go get -u github.com/cosmtrek/air
	air -c .air.toml
# save github token in an environment variable export GITHUB_TOKEN="token"
add-tag:
ifeq (,$(findstring v,$(tag)))
	@echo "error : tag needs to be of format v0.x.x. Usage --> make upload tag=v0.x.x"
	@echo
	exit 1
endif
	git fetch
	git tag $(tag)
	git push origin --tags
upload: add-tag install build-linux #make upload tag=v0.x.x, install --> brew install goreleaser
	goreleaser --rm-dist
test:
	go test -v ./...
start-ws:
	go run ws/ws.go