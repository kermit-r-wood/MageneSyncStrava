# OnelapSyncStrava

自动将顽鹿的骑行数据同步到 Strava。

## 项目背景

顽鹿运动（Onelap）此前支持将运动数据自动同步至 Strava，但该功能于 2026 年 3 月 19 日关闭。本项目旨在恢复这一功能，通过 API 自动将顽鹿的骑行数据同步到 Strava。


## 功能

- **顽鹿登录认证**: 使用 MD5 签名安全登录
- **活动筛选**: 自动提取并过滤尚未同步的当天骑行记录
- **FIT 文件处理**: 获取最完整的 FIT 格式的活动数据并上传
- **Strava OAuth**: 自动验证并刷新 Strava 访问令牌，无需手动维护
- **多种运行模式**: 提供 `sync`, `auth`, `check` 和 `status` 等多种子命令，方便配置和调试
- **本地状态管理**: 同步记录及访问令牌持久化保存在本地 (`state.json` 和 `config.json`)，避免重复上传
- **Agent 专属向导**: 内置 `sync_wizard` Skill，支持 AI 助手一步步引导完成环境和配置搭建

## 前置要求

- 顽鹿账号
- Strava API 应用（[在此创建](https://www.strava.com/settings/api)）
  - **重要配置**：在 Strava 设置页面，将 **Authorization Callback Domain** 设置为 `localhost`。
  - 创建后获取 `Client ID` 和 `Client Secret`。

## 下载与安装

### 使用预编译的二进制（推荐）

可以通过 Github Actions 下载自动构建好的最新二进制文件。支持 Windows, macOS 和 Linux 平台。
[前往下载 Releases](https://github.com/kermit-r-wood/OnelapSyncStrava/releases) 或者在 Actions 详情页面下载 Artifacts。

### 源码编译

需要 Go 1.21+ 环境:

```bash
git clone https://github.com/kermit-r-wood/OnelapSyncStrava.git
cd OnelapSyncStrava
make build
# 或者直接使用 go 命令编译
go build -o OnelapSyncStrava main.go
```

## 配置指南

将项目根目录下的 `config.sample.json` 复制并重命名为 `config.json`，按需填入你的凭证：

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

> **为什么填写了 client_id 和 client_secret 后还需要授权？**
> 1. `client_id` 和 `client_secret` 是你的 **Strava API 应用凭证**，用于向 Strava 标识本程序。
> 2. **授权流程**（执行 `auth` 命令）是由于 Strava 的安全机制，需要你作为 **用户** 亲自同意授权给该应用上传数据的权限。
> 3. 授权成功后，程序会自动获取 `access_token` 和 `refresh_token` 并保存到 `config.json` 中，之后即可实现全自动同步，无需再次手动授权。


## 使用教程

基础命令格式：
```bash
./OnelapSyncStrava [command]
```
（如果不加 `[command]`，默认执行 `sync` 任务）

### 1. 检查配置 (`check`)

配置完毕后，首先运行 `check` 命令检查连通性是否正常：

```bash
./OnelapSyncStrava check
```
程序将分别测试 顽鹿 以及 Strava API 的连接状态及凭据有效性。

### 2. 授权 Strava (`auth`)

使用此命令启动内置验证服务进行 Strava 授权（由于访问限制，初次运行必须执行此步）：

```bash
./OnelapSyncStrava auth
```

程序会自动打开浏览器并前往 Strava 进行授权。
**重要**：授权请求中必须勾选 **「Upload your activities and posts to Strava」**，否则即使成功拉取也不能上传数据。
授权完成后，Token 会自动写入至 `config.json` 中保存，以后无需再重复操作。

*（提示：如果是运行在无头远程服务器等无法浏览器重定向的环境，你可以将终端输出的授权链接复制到本地浏览器中访问，授权同意之后再把最终报错跳转地址附带的 `?code=xxx` 重新手动回调请求给程序即可生效，或者请求智能助手帮忙代理处理。）*

### 3. 数据同步 (`sync`)

直接运行主程序或者附带 sync 参数，执行数据同步：

```bash
./OnelapSyncStrava sync
```

该模式会：
1. 登录顽鹿
2. 获取当天的骑行活动
3. 判断是否已在 `state.json` 中标记为同步完成，从而过滤重复任务
4. 检查并刷新 Strava 令牌（如失效）
5. 下载对应的 FIT 文件并将其上传到 Strava
6. 更新本地持久化记录 `state.json`

### 4. 查看状态 (`status`)

你可以随时通过此命令快速检查当前环境配置文件和历史同步情况：

```bash
./OnelapSyncStrava status
```
展示当前的账户设定、Strava验证状态，以及历史成功同步的骑行活动条目数。

## 定时执行（可选）

建议通过计划任务实现全自动后台同步。
- **Windows**: 使用 **任务计划程序** 定时执行 `OnelapSyncStrava.exe`。
- **Linux/macOS**: 配置 `crontab` 定时任务（例如每日晚间定期执行一次）。

## 项目结构

```
OnelapSyncStrava/
├── main.go                     # 项目入口与命令路透
├── config.json                 # 运行配置文件（需手动创建，勿提交）
├── config.sample.json          # 模板配置文件
├── state.json                  # 已同步记录状态（自动生成）
├── Makefile                    # 快捷构建脚本
├── .github/workflows/          # CI/CD 自动发布配置
├── .agents/skills/sync_wizard/ # 针对 Agent 辅助工具的使用指南
├── internal/
│   ├── config/                 # 负责读取配置和记录同步状态 
│   ├── onelap/                 # 顽鹿 API 客户端代码逻辑
│   └── strava/                 # Strava OAuth 与上传交互实现
├── go.mod
└── go.sum
```

## 许可证

MIT
