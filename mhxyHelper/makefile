.PHONY: copy_db
# copy_db
copy_db:
	cp mhxyhelper.db ./bin
	cp mhxyhelper.db-shm ./bin
	cp mhxyhelper.db-wal ./bin

.PHONY: build
# build
build:
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64
	mkdir -p bin/
	go build -ldflags "-X main.Version=v0.1" -o ./bin/mhxy_helper ./gin_cmd/main.go
	go build -ldflags "-X main.Version=v0.1" -o ./bin/mhxy_cmd_helper ./cmd/main.go
	go build -ldflags "-X main.Version=v0.1" -o ./bin/mhxy_account_helper ./account_book_cmd/main.go
	cp mhxyhelper.db ./bin
	cp mhxyhelper.db-shm ./bin
	cp mhxyhelper.db-wal ./bin

.PHONY: build_linux
# build_linux
build_linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64
	go build -ldflags "-X main.Version=v0.1" -o ./bin/mhxy_helper_linux ./gin_cmd/main.go
	go build -ldflags "-X main.Version=v0.1" -o ./bin/mhxy_cmd_helper_linux ./cmd/main.go
	go build -ldflags "-X main.Version=v0.1" -o ./bin/mhxy_account_helper_linux ./account_book_cmd/main.go

.PHONY: build_win
# build_win
build_win:
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64
	go build -ldflags "-X main.Version=v0.1" -o ./bin/mhxy_helper.exe ./gin_cmd/main.go
	go build -ldflags "-X main.Version=v0.1" -o ./bin/mhxy_cmd_helper.exe ./cmd/main.go
	go build -ldflags "-X main.Version=v0.1" -o ./bin/mhxy_account_helper.exe ./account_book_cmd/main.go

.PHONY: docker_build
# docker_build
docker_build:
	docker build -f dockerfile -t mhxy_helper:v0.0.1 .

.PHONY: docker_run   # TODO CGO不允许跨平台编译 需要编译过程放到Docker中完成 https://blog.csdn.net/Maggie_up/article/details/83821386
# docker_run
docker_run:
	docker run --name mhxy-helper -d -p 8080:8080 mhxy_helper:v0.0.1