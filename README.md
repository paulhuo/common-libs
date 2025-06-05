# common-go

**common-go** 是公司内部维护的 Golang 通用基础库，旨在统一各项目的基础能力、提升代码复用性，涵盖工具函数、HTTP 中间件、gRPC 拦截器等模块。

---

## 📁 模块结构

### 🔧 `utils/`
通用工具函数模块，**不依赖具体业务逻辑**，适用于所有 Go 项目：

- `stringx`：字符串处理（大小写转换、随机字符串、截断）
```go
// 示例
snake := stringx.ToSnake("HelloWorld") // hello_world
randStr := stringx.Random(8)           // 随机字符串
```

- `timeutil`：时间处理（格式化、解析、区间计算）
- `ctxx`：请求上下文扩展（trace-id、app-id、user-id 提取与注入）
- `httpx`：简洁的 HTTP 客户端封装（支持 JSON/form/multipart）
- `logx`：日志封装（支持 trace-id 注入，zap/logrus 可插拔）
- `idgen`：唯一 ID 生成（UUID、雪花ID、NanoID）
- `jsonx`：容错 JSON 编解码（结构体/map 互转、空值兼容）
- `envx`：环境变量读取工具（支持默认值和类型转换）
- `validate`：字段校验工具（手机号、邮箱、URL、IP 等）
- `mathx`：数字计算工具（保留小数、百分比、四舍五入）
- `filex`：文件读写和路径检查（是否存在、临时文件）

---

### 🌐 `middlewares/`
通用 HTTP 中间件集合（适配 gin / go-zero 等框架）：

- `traceid`：请求链路追踪 ID 注入与传播
- `logger`：统一日志记录中间件
- `recover`：统一 panic 恢复与错误响应封装

---

### ⚙️ `interceptors/`
gRPC 拦截器集合，适用于统一日志、错误处理、监控等场景：

- `traceid`：gRPC trace-id 注入与传播
- `logger`：统一记录请求日志
- `metrics`：Prometheus 指标采集拦截器（可选）

---

## 开发规范
- 每个模块应自带单元测试 (*_test.go)
- 所有工具函数应无副作用、线程安全
- 避免将业务逻辑混入公共库
- 所有可导出函数需加上注释说明
- gRPC 拦截器统一支持 context.Context 中传递 trace-id
