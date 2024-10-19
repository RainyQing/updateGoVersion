# updateGoVersion
## Golang 安装脚本文档

该 Golang 安装脚本用于自动更新/下载、安装并配置最新版本的 Golang（Go 编程语言），支持 Linux 和 Windows 系统。脚本提供了环境变量设置、Go 模块缓存配置、Go 代理配置等功能。它会检测并仅在有新版本时进行更新安装。

### 功能特点

- 自动检测现有的 Go 版本，并在必要时进行升级。
- 下载适用于操作系统的最新稳定版 Golang。
- 安装 Golang 到指定目录或默认目录。
- 自动设置 `GOROOT` 和 Go 代理设置。
- 配置 Go 模块缓存目录。
- 支持 Linux 和 Windows 平台。
- 由于windows 系统配置环境变量是配置到系统变量，所以需要管理员权限。

### 前置要求

在运行已编译程序/原go脚本之前，请确保以下条件满足：
- **网络要求**: 该脚本需要互联网连接访问下载最新版本 , 但是因为检测最新版本接口是调用go官网 , 或许会有链接失败情况。
- **已编译程序**：已编译过后得程序不需要环境 , 仅需要有网
- **原go脚本**：脚本需要安装go环境。
- **注意**: linux版本未经过严格测试 , 或许会有bug , 如有bug按需修改过后在使用

### 安装步骤

1. 克隆或下载该脚本到本地环境：

   ```bash
   git clone <repository_url>
   ```

2. 执行安装脚本：
    - **Linux 系统**：执行 已编译好的`updateGoVersion` 程序或原`updateGoVersion.go` 脚本。

      ```bash
      updateGoVersion
      或者
      go run updateGoVersion.go
      ```

    - **Windows 系统**：执行 已编译好的`updateGoVersion.exe` 程序(需要管理员权限)或原`updateGoVersion.go` 脚本。

      ```powershell
      .\updateGoVersion.exe
      或者
      go run updateGoVersion.go // 注意：需要安装go环境且命令行需要管理员权限
      ```

3. 脚本将会自动检测现有的 Go 版本，并下载最新版本的 Go 安装包（如果需要更新）。


### 环境变量设置

脚本会自动将 Go 的环境变量配置到系统环境中，包括：

- `GOROOT`：Go 安装路径。
- `GO111MODULE`：启用 Go modules。
- `GOMODCACHE`：Go 模块缓存路径。
- `GOPROXY`：Go 代理设置，默认使用 `https://proxy.golang.org` (七牛云)，可以根据需要修改。


### 常见问题

1. **如何验证安装是否成功？**

   在终端输入以下命令，查看 Go 版本：

   ```bash
   go version
   ```

   输出类似 `go version go1.x.x linux/amd64` 即表示安装成功。

2. **如何更新到最新版本？**

   运行脚本时，脚本会自动检测是否有新版本可用，并自动更新。

### 更新日志
v1.1.0 优化按任意键退出

v1.2.0    1. 优化代码 2. 配置如果本地已经安装go环境则不会设置环境变量和go mod缓存目录以及goproxy
### 贡献指南

欢迎通过 GitHub 提交 issues 和 pull requests 来改进此脚本。

---

通过此文档，您应该能够轻松使用该脚本安装并配置 Golang。如果有任何问题或改进建议，请随时反馈！
