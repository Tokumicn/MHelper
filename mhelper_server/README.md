

涉及到的gin关联组件的使用案例：
1. 配置文件组件：github.com/spf13/viper
2. 数据库访问组件：github.com/jinzhu/gorm
3. 日志组件：gopkg.in/natefinch/lumberjack.v2
4. 接口文档组件：swagger
5. 接口参数校验：github.com/go-playground/validator/v10
6. 国际化组件：go-playground/locales 、 go-playground/universal-translator
7. 访问控制(Token)：github.com/dgrijalva/jwt-go
8. 中间件：
    - 访问日志
    - 异常捕获处理
    - Recovery
    - 接口限流：github.com/juju/ratelimit 
    - 统一超时处理
9. 邮件报警处理：gopkg.in/gomail.v2
10. 链路追踪：TODO
11. 编译：TODO
12. 优雅启停
13. 热更新：不停机更新 TODO

 ```
gin-tlp
├── configs
├── docs
├── global
├── internal
│   ├── dao
│   ├── middleware
│   ├── model
│   ├── routers
│   └── service
├── pkg
├── storage
├── scripts
└── third_party
```   

- configs：配置文件。
- docs：文档集合。
- global：全局变量。
- internal：内部模块。
- dao：数据访问层（Database Access Object），所有与数据相关的操作都会在 dao 层进行，例如 MySQL、ElasticSearch 等。
- middleware：HTTP 中间件。
- model：模型层，用于存放 model 对象。
- routers：路由相关逻辑处理。
- service：项目核心业务逻辑。
- pkg：项目相关的模块包。
- storage：项目生成的临时文件。
- scripts：各类构建，安装，分析等操作的脚本。
- third_party：第三方的资源工具，例如 Swagger UI。


docker run for mysql:
```txt
# 使用 mysql:8.4 镜像  设置端口、账号、密码
docker run -itd --name mysql-test -p 3306:3306 -e MYSQL_ROOT_PASSWORD=rootroot mysql:8.4
```




