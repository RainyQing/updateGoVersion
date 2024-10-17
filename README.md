# updateGoVersion
该 Golang 安装脚本用于自动更新/下载、安装并配置最新版本的 Golang（Go 编程语言），支持 Linux 和 Windows 系统。脚本提供了环境变量设置、Go 模块缓存配置、Go 代理配置等功能。它会检测并仅在有新版本时进行更新安装。

功能特点
自动检测现有的 Go 版本，并在必要时进行升级。
下载适用于操作系统的最新稳定版 Golang。
安装 Golang 到指定目录或默认目录。
自动设置 GOROOT、GOPATH 和 Go 代理设置。
配置 Go 模块缓存目录。
支持 Linux 和 Windows 平台。
前置要求
在运行脚本之前，请确保以下条件满足：

Go 命令行工具：脚本需要安装 go 命令行工具（至少用于检测现有 Go 环境）。
互联网连接：该脚本需要访问 Go 官方网站以下载最新版本。
安装步骤
克隆或下载该脚本到本地环境：

bash
复制代码
git clone <repository_url>
执行安装脚本：

Linux 系统：执行 install-go.sh 脚本。

bash
复制代码
sudo bash install-go.sh
Windows 系统：执行 install-go.ps1 脚本（需要管理员权限）。

powershell
复制代码
.\install-go.ps1
脚本将会自动检测现有的 Go 版本，并下载最新版本的 Go 安装包（如果需要更新）。

参数说明
INSTALL_DIR：自定义安装路径。如果未指定，脚本将默认安装到系统推荐目录（例如 /usr/local/go）。

bash
复制代码
INSTALL_DIR="/your/custom/path" bash install-go.sh
GO_PROXY：自定义 Go 代理地址，默认使用 https://proxy.golang.org。

bash
复制代码
GO_PROXY="https://goproxy.cn" bash install-go.sh
环境变量设置
脚本会自动将 Go 的环境变量配置到系统环境中，包括：

GOROOT：Go 安装路径。
GOPATH：Go 工作空间路径，默认位于用户主目录下。
GO111MODULE：启用 Go modules。
GOMODCACHE：Go 模块缓存路径。
GOPROXY：Go 代理设置，默认使用 https://proxy.golang.org，可以根据需要修改。
示例
在 Linux 系统上安装
bash
复制代码
sudo bash install-go.sh
在 Windows 系统上安装
powershell
复制代码
.\install-go.ps1
常见问题
如何验证安装是否成功？

在终端输入以下命令，查看 Go 版本：

bash
复制代码
go version
输出类似 go version go1.x.x linux/amd64 即表示安装成功。

如何更新到最新版本？

运行脚本时，脚本会自动检测是否有新版本可用，并自动更新。

贡献指南
欢迎通过 GitHub 提交 issues 和 pull requests 来改进此脚本。

通过此文档，您应该能够轻松使用该脚本安装并配置 Golang。如果有任何问题或改进建议，请随时反馈！
