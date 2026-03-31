# MageneSync

自动将顽鹿的骑行数据同步到 Strava。

## 功能

- **顽鹿登录认证**: 使用 MD5 签名安全登录
- **活动筛选**: 自动筛选当天的骑行记录
- **FIT 文件处理**: 下载 FIT 格式的活动数据（骑行数据最丰富的格式）
- **Strava OAuth**: 自动刷新 Strava 访问令牌，无需手动维护
- **内置授权流程**: 一条命令完成 Strava OAuth 授权，无需手动复制 token
- **持久化配置**: 凭证和令牌保存在 `config.json` 中，自动更新

## 前置要求

- Go 1.21+
- 顽鹿账号
- Strava API 应用（[在此创建](https://www.strava.com/settings/api)）
  - **重要配置**：在 Strava 设置页面，将 **Authorization Callback Domain** 设置为 `localhost`。
  - 创建后获取 `Client ID` 和 `Client Secret`。


## 使用步骤

### 1. 克隆并构建

```bash
git clone <repo-url>
cd MageneSync
go mod tidy
go build -o MageneSync.exe .
```

### 2. 配置

将 `config.sample.json` 复制为 `config.json`，填入你的凭证：

```json
{
  "onelap": {
    "account": "你的顽鹿账号",
    "password": "你的顽鹿密码"
  },
  "strava": {
    "client_id": "你的Strava Client ID",
    "client_secret": "你的Strava Client Secret",
    "access_token": "",
    "refresh_token": "",
    "expires_at": 0
  }
}
```

> **说明**: Strava 部分只需要填写 `client_id` 和 `client_secret`，token 会在下一步自动获取。

### 3. 授权 Strava

运行内置的授权命令：

```bash
# 使用 go run
go run . auth

# 或使用编译后的二进制
./MageneSync.exe auth
```

程序会自动打开浏览器跳转到 Strava 授权页面。
**重要**：请确保在授权页面勾选 **「Upload your activities and posts to Strava」** 权限，否则程序无法上传数据。

授权成功后，Token 会自动保存到 `config.json`。此操作只需执行**一次**。


## 日常使用

```bash
# 直接运行
go run .

# 或使用编译后的二进制
./MageneSync.exe
```

程序会依次执行：

1. 登录顽鹿
2. 获取今天的骑行记录
3. 下载 FIT 文件到临时目录
4. 刷新 Strava 访问令牌（如已过期）
5. 将 FIT 文件上传到 Strava
6. 清理临时文件

## 定时执行（可选）

可以使用 **Windows 任务计划程序** 定时运行 `MageneSync.exe`（例如每天晚上），实现全自动同步。

## 项目结构

```
MageneSync/
├── main.go                     # 入口，子命令路由
├── config.json                 # 运行时配置（不提交到 git）
├── config.sample.json          # 配置模板
├── internal/
│   ├── config/config.go        # 配置加载与保存
│   ├── onelap/onelap.go        # 顽鹿 API 客户端
│   └── strava/
│       ├── strava.go           # Strava 上传与令牌刷新
│       └── auth.go             # Strava OAuth 授权流程
├── go.mod
└── go.sum
```

## 许可证

MIT
