## 1 层级依赖关系
- router: 路由
- controller：处理参数，响应
- handler: 处理业务逻辑
- service: 通用的工具、业务逻辑
- repo: 数据库操作
- model: 数据模型
```text
         engine
            |
         router
            |
        controller   
            |
          handler
           /    \     
       service   |
          |      |
           \    /
            repo
              |
            model              
```

## 2 编码规则
- NewXXX：创建controller、handler、service、repo对象，需使用构造函数。
- server/router.go：此处编写router规则。
- server/wire.go：此处编写wire注入规则。
- wire gen：每次更改wire注入规则后，需要重新运行wire。

## 3 配置
- mysql：默认使用读写分离配置。
- config.yaml：可放于workpwd，或workpwd/config/config.yaml。

## 4 功能特性
- [x] logrous+file-rotatelogs
  - [x] traceId、spanId
  - [x] caller信息
  - [x] 按日分表、文件留存时间
  - [x] 支持console/file输出切换
- [x] gorm
  - [x] 读写分离
  - [x] 按日分表、按mode分表
- [x] redis
- [x] wire