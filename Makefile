app = himawari8
src = main.go
dst = $(app)
ver = 0.0.3
ts = $(shell date +"%Y%m%d%H%M%S")
commit = $(shell git log -1 --pretty=%h)
gitstat = $(shell (git status --porcelain | grep -q .) && echo dirty || echo clean)

default: test

run:
	@go run -race . -stderrthreshold=INFO -logtostderr=true

test: $(src)
	@gofmt -w -s $(src)
	@goimports -w $(src)
	@go vet ./...
	@go test ./...
	@staticcheck ./...
	@govulncheck ./...

build_brew: $(src)
	@go build -o $(app) -ldflags "\
		-X 'github.com/pihao/himawari8/main.Version=$(ver)' \
		-X 'github.com/pihao/himawari8/main.Commit=$(commit) $(gitstat)' \
		-X 'github.com/pihao/himawari8/main.Ts=$(ts)'" \
		.

build_amd64: test
	@rm -f $(app) $(app)-darwin-amd64-*.tar.gz
	@GOOS=darwin GOARCH=amd64 go build -o $(app) -ldflags "\
		-X 'github.com/pihao/himawari8/main.Version=$(ver)' \
		-X 'github.com/pihao/himawari8/main.Commit=$(commit) $(gitstat)' \
		-X 'github.com/pihao/himawari8/main.Ts=$(ts)'" \
		.
	@tar czvf $(app)-darwin-amd64-v$(ver).tar.gz $(dst)

build_arm64: test
	@rm -f $(app) $(app)-darwin-arm64-*.tar.gz
	@GOOS=darwin GOARCH=arm64 go build -o $(app) -ldflags "\
		-X 'github.com/pihao/himawari8/main.Version=$(ver)' \
		-X 'github.com/pihao/himawari8/main.Commit=$(commit) $(gitstat)' \
		-X 'github.com/pihao/himawari8/main.Ts=$(ts)'" \
		.
	@tar czvf $(app)-darwin-arm64-v$(ver).tar.gz $(dst)
