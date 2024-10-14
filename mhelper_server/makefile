.PHONY: echo build tidy run swagger bindata-config

echo:
	@echo "mhelper_server"

build:
	go mod tidy
	go build -ldflags "-X main.version=0.0.1" -o mhelper-server main.go

tidy:
	go mod tidy

run:
	./mhelper_server -port 8080 -mode debug

swagger:
	swag init

# 将 config.yaml 生成为静态文件供服务使用
# 源文件在 configs/config.yaml  生成到 configs/config.go 中
bindata-config:
	go-bindata --nocompress -pkg configs -o configs/config.go configs/...




#=====================================以下为接口测试相关=========================================

#测试Token获取
curl-auth: # 通过 form-data方式传递参数
	curl -X POST \
      'http://127.0.0.1:8080/auth' \
      -F app_key=tom \
      -F app_secret=tom-svc

curl-auth2:  #通过 Content-Type: application/x-www-form-urlencoded 方式传递参数
	curl -X POST \
        'http://127.0.0.1:8080/auth' \
     -d 'app_key=tom&app_secret=tom-svc'

#测试token校验
curl-auth-check:
	curl -X GET http://127.0.0.1:8000/api/v1/tags  # 未带token，报错
	curl -X GET http://127.0.0.1:8000/api/v1/tags -H 'token: eyJhbGciOiJIUzI1NiIsInRxxx' # 带Token，成功


#测试查看上传文件
curl-fileupload:
	curl -X POST http://localhost:8080/upload/file \
	  -F 'file=@/Users/zhangrui/Downloads/302004_青背龙.png' \   # TODO 待修改
	  -F type=1

#测试tags相关接口-[新建]
curl-tags-create:
	curl -X POST http://127.0.0.1:8080/api/v1/tags -F 'name=Go' -F created_by=tom
	curl -X POST http://127.0.0.1:8080/api/v1/tags -F 'name=.Net' -F created_by=tom
	curl -X POST http://127.0.0.1:8080/api/v1/tags -F 'name=Rust' -F created_by=tom

#测试tags相关接口-[查询]
curl-tags-get:
	curl -X GET 'http://127.0.0.1:8080/api/v1/tags?page=1&pageSize=2'
	curl -X GET 'http://127.0.0.1:8080/api/v1/tags?page=2&pageSize=2'

#测试tags相关接口-[修改]
curl-tags-put:
	curl -X PUT http://127.0.0.1:8080/api/v1/tags/1 -F state=1 -F modified_by=tom

#测试tags相关接口-[修改]
curl-tags-del:
	curl -X DELETE  http://127.0.0.1:8080/api/v1/tags/1

