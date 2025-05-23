.PHONY: copy_db
# copy_db
copy_db:
	cp mhxyhelper.db ./bin
	cp mhxyhelper.db-shm ./bin
	cp mhxyhelper.db-wal ./bin

.PHONY: build
# build
build:
	GOOS=darwin GOARCH=amd64  # CGO_ENABLED=0
	mkdir -p bin/
	go build -ldflags "-X main.Version=v0.1" -o ./bin/mhxy_web ./sever/gin_cmd/main.go
	go build -ldflags "-X main.Version=v0.1" -o ./bin/mhxy_cmd ./sever/cmd/main.go
	cp ./mhxyhelper.db ./bin
	cp ./mhxyhelper.db-shm ./bin
	cp ./mhxyhelper.db-wal ./bin

.PHONY: build_linux
# build_linux
build_linux:
	GOOS=linux GOARCH=amd64 # CGO_ENABLED=0
	go build -ldflags "-X main.Version=v0.1" -o ./bin/mhxy_web_linux ./sever/gin_cmd/main.go
	go build -ldflags "-X main.Version=v0.1" -o ./bin/mhxy_cmd_linux ./sever/cmd/main.go

.PHONY: build_win
# build_win
build_win: # CGO_ENABLED=0
	GOOS=windows GOARCH=amd64
	go build -ldflags "-X main.Version=v0.1" -o ./bin/mhxy_web.exe ./sever/gin_cmd/main.go
	go build -ldflags "-X main.Version=v0.1" -o ./bin/mhxy_cmd.exe ./sever/cmd/main.go

.PHONY: docker_build
# docker_build
docker_build:
	docker build -f dockerfile -t mhxy_helper:v0.0.1 .

.PHONY: docker_run   # TODO CGO不允许跨平台编译 需要编译过程放到Docker中完成 https://blog.csdn.net/Maggie_up/article/details/83821386
# docker_run
docker_run:
	docker run --name mhxy-helper -d -p 8080:8080 mhxy_helper:v0.0.1