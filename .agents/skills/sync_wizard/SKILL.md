---
name: sync_wizard
description: 引导用户配置、授权并运行 MageneSync 同步工具。提供从源码编译或从 GitHub 获取二进制程序的详细说明。
---

# MageneSync 助手 (Sync Wizard)

当调用此 Skill 时，Agent 应作为同步专家，引导用户获取并使用 `MageneSync` 工具完成从 Onelap (Magene) 到 Strava 的活动同步。

## 1. 获取程序 (Get Program)

在使用工具之前，必须确保拥有对应操作系统的二进制程序。

### 方法 A：从 GitHub 下载（推荐）
1.  **稳定版**：前往 [GitHub Releases](https://github.com/kermit-r-wood/MageneSyncStrava/releases) 下载最新版本的压缩包，解压出对应平台的二进制文件（如 `MageneSync-windows-amd64.exe`）。
2.  **开发版**：从 GitHub Actions 的 [Build Binaries](https://github.com/kermit-r-wood/MageneSyncStrava/actions/workflows/binaries.yml) 工作流中下载最新的 Artifacts。

### 方法 B：从源码编译
如果本地已安装 Go 环境：
1.  **Windows**: 运行 `go build -o MageneSync.exe main.go` 或 `make build`。
2.  **Linux/macOS**: 运行 `go build -o MageneSync main.go` 并确保文件具有执行权限 (`chmod +x MageneSync`)。

## 2. 核心流程 (Core Workflow)

### 第一步：环境检查
1.  确认当前目录下存在 `MageneSync` 程序。
2.  运行 `MageneSync status` (Windows 使用 `.\MageneSync.exe status`) 查看当前同步统计和配置完整性。

### 第二步：基础配置 (获取 API 凭证)
1.  **获取 Strava API 凭证** (初次使用必做)：
    -   在浏览器中登录 Strava 并访问：[Strava API Settings](https://www.strava.com/settings/api)。
    -   创建一个新的 API 应用（若已有则跳过此步），填写必要信息。
    -   **关键点：`Authorization Callback Domain` 字段必须填写为 `localhost`**。
    -   创建成功后，在页面中找到 `Client ID` 及其对应的 `Client Secret`。
2.  设置并确保 `config.json` 包含正确的参数：
    -   `onelap_account`: Onelap 登录邮箱/手机。
    -   `onelap_password`: Onelap 密码。
    -   `strava_client_id`: 上一步获取的 Client ID。
    -   `strava_client_secret`: 上一步获取的 Client Secret。
3.  运行 `MageneSync check` 进行连通性测试。

### 第三步：Strava 授权 (OAuth)
1.  如果 `check` 提示 `missing refresh_token`，运行 `MageneSync auth`。
2.  **IM/无头环境下的授权引导 (OOB Fallback)**：
    -   Agent 提取终端输出的 Strava 授权链接，发送给用户，提示其在手机或个人计算机浏览器中打开并授权。
    -   向用户解释：由于在远程/IM环境下，用户的浏览器无法直接重定向到服务器的 `localhost`，所以授权后网页最终会显示访问报错（这是正常现象）。
    -   **关键交互指令**：通知用户：“请复制授权完成后浏览器地址栏那条报错的完整链接（必须包含 `?code=...` 等参数），并将其作为消息发送给我”。
    -   Agent 收到来自用户的带有授权码的链接后，在终端执行：`curl "[用户发回的完整链接]"`。
    -   此时本地 `MageneSync auth` 挂起的进程会捕获到该 `curl` 模拟的回调请求，结束堵塞、保存 Token 并完成授权。

### 第四步：执行同步
1.  运行 `MageneSync sync` 开始抓取并拉取活动。
2.  同步完成后，告知用户新增的活动数量。

## 3. 故障排除
-   **授权失败**：检查 Strava Client ID/Secret 是否正确，或尝试重新运行 `auth`。
-   **登录失败**：检查 Onelap 账号密码及网络连接。
-   **无新活动**：确认 Onelap 中是否有今日或近期尚未同步的记录。
