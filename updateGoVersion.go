package main

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"time"
)

type File struct {
	Filename string `json:"filename"`
	OS       string `json:"os"`
	Arch     string `json:"arch"`
	Version  string `json:"version"`
	Sha256   string `json:"sha256"`
	Size     int64  `json:"size"`
	Kind     string `json:"kind"`
}

type GoVersion struct {
	Version string `json:"version"`
	Stable  bool   `json:"stable"`
	Files   []File `json:"files"`
}

// 定义一个结构体来存储 go env 的输出信息
type GoEnv struct {
	GO111MODULE    string
	GOARCH         string
	GOBIN          string
	GOCACHE        string
	GOENV          string
	GOEXE          string
	GOEXPERIMENT   string
	GOFLAGS        string
	GOHOSTARCH     string
	GOHOSTOS       string
	GOINSECURE     string
	GOMODCACHE     string
	GONOPROXY      string
	GONOSUMDB      string
	GOOS           string
	GOPATH         string
	GOPRIVATE      string
	GOPROXY        string
	GOROOT         string
	GOSUMDB        string
	GOTMPDIR       string
	GOTOOLCHAIN    string
	GOTOOLDIR      string
	GOVCS          string
	GOVERSION      string
	GODEBUG        string
	GOTELEMETRY    string
	GOTELEMETRYDIR string
	GCCGO          string
	GOAMD64        string
	AR             string
	CC             string
	CXX            string
	CGO_ENABLED    string
	GOMOD          string
	GOWORK         string
	CGO_CFLAGS     string
	CGO_CPPFLAGS   string
	CGO_CXXFLAGS   string
	CGO_FFLAGS     string
	CGO_LDFLAGS    string
	PKG_CONFIG     string
	GOGCCFLAGS     string
}

// 下载网站链接
var downloadURL = "https://studygolang.com/dl/golang/"

// 获取golang版本信息接口
var linux = "linux"

// 获取golang版本信息接口
var windows = "windows"

// 获取golang版本信息接口
var latestVersionURL = "https://go.dev/dl/?mode=json"

// goproxy , 默认使用七牛云
var goProxy = "https://goproxy.cn,direct"

func main() {
	// 获取最新版本号, 下载链接 , 文件名
	latestVersion, url, fileName := getLatestGoVersion()
	// 尝试通过 `go env` 命令获取详细环境信息
	cmd := exec.Command("go", "env")
	output, err := cmd.CombinedOutput()
	filePath := ""
	env := GoEnv{}
	if err != nil {
		log.Println("未检测到 Go 环境 , 最新版本: ", latestVersion)
	} else {
		goEnv := strings.TrimSpace(string(output))

		// 解析 `go env` 输出
		env, err = parseGoEnv(goEnv)
		if err != nil {
			log.Printf("解析 `go env` 输出失败: %v\n", err)
			return
		}
		if latestVersion <= env.GOVERSION {
			log.Print("本地golang版本: ", env.GOVERSION, " 已经是最新版本, 不需要更新.")
			return
		}
		//如果成功解析 `go env` 输出, 则获取 `GOROOT` 路径
		filePath = env.GOROOT
		log.Print("本地golang版本: ", env.GOVERSION, " 最新版本: ", latestVersion)
	}
	//下载最新版本
	downloadFile(url, fileName)
	//如果未安装go环境让用户输入安装路径
	input := ""
	if filePath == "" {
		//请输入安装文件夹
		log.Print("请输入安装文件夹, 按回车键使用默认路径: ")
		fmt.Scanln(&input)
		//如果输入路径不为空
		if input != "" {
			filePath = input
			//判断当前环境
			if runtime.GOOS == windows {
				//判断输入路径是否符合windows格式
				isValidWindowsPath(filePath)
				//判断输入路径是否符合windows格式
			} else if runtime.GOOS == linux {
				isValidLinuxPath(filePath)
			} else {
				//其他系统暂不支持
				panic("不支持的操作系统")
			}
		} else {
			//如果输入路径为空, 则默认使用用户主目录
			homeDir, err := os.UserHomeDir()
			if err != nil {
				fmt.Println("获取用户主目录失败:", err)
				return
			}
			filePath = filepath.Join(homeDir, "go")
		}
	}
	log.Printf("安装路径: %s", filePath)
	// 安装最新版本
	installGo(filePath, fileName)
	//删除下载文件
	if err := os.Remove(fileName); err != nil {
		log.Printf("删除下载文件 %s 失败: %v", fileName, err)
	}
	log.Print("开始自动配置环境变量.")

	//设置环境变量
	setEnv(filePath)
	//设置gomod缓存目录
	cache := setGoModCache(filePath, err)
	//设置goproxy
	setGoProxy()
	//打印安装成功信息
	printInstallationInfo(latestVersion, filePath, cache)
	// 等待用户输入
	log.Printf("按任意键退出......")
	fmt.Scanln(&input)
	//完成安装退出
	log.Printf("程序已退出.")
}

func setGoModCache(filePath string, err any) string {
	//设置gomod缓存目录 , 在安装目录下创建 gomodcache 目录作为gomod缓存目录
	if runtime.GOOS == windows {
		filePath = filePath + "\\" + "gomodcache"
	} else if runtime.GOOS == linux {
		filePath = filePath + "/" + "gomodcache"
	}
	//校验gomodcache是否存在
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		// 创建 gomodcache 目录
		if err := os.MkdirAll(filePath, 0755); err != nil {
			log.Printf("创建 gomodcache 目录失败: %v", err)
			panic("创建gomod缓存目录失败.")
		}
	}
	gomodcache := "GOMODCACHE=" + filePath
	// 创建要执行的命令
	cmd := exec.Command("go", "env", "-w", gomodcache)

	// 运行命令并捕获输出
	_, err = cmd.CombinedOutput()
	if err != nil {
		panic("配置gomod缓存目录出错.")
	}
	return filePath
}

func setGoProxy() {
	// 创建要执行的命令
	cmd := exec.Command("go", "env", "-w", goProxy)

	// 运行命令并捕获输出
	_, err := cmd.CombinedOutput()
	if err != nil {
		panic("配置gomod缓存目录出错.")
	}
}

// 安装成功后的详细信息
func printInstallationInfo(goRoot string, goVersion string, gomodcache string) {
	log.Printf("golang环境安装成功！")
	log.Printf("安装详细信息：")
	log.Printf("-------------------------")
	log.Printf("版本: " + goVersion)         // 示例版本
	log.Printf("安装路径: " + goRoot)          // 示例安装路径
	log.Printf("gomod缓存目录: " + gomodcache) // gomod缓存目录
	log.Printf("环境变量: GOROOT 已设置")         // 环境变量设置情况
	log.Printf("使用 `go env` 查看详细信息")
	log.Printf("使用方法: 使用 'go run' 命令运行 Go 程序")
	log.Printf("-------------------------")
}

func setEnv(path string) {
	//如果是windows系统
	if runtime.GOOS == windows {

		cmd := exec.Command("setx", "GOROOT", path, "/M") // "/M" 是设置系统变量
		if err := cmd.Run(); err != nil {
			log.Fatalf("设置环境变量 %s 失败: %v", path, err)
		}

		path = "%GOROOT%" + "\\bin"

		// 获取当前 PATH 环境变量
		cmd = exec.Command("powershell", "-Command", "echo $env:Path")
		output, err := cmd.Output()
		if err != nil {
			log.Fatalf("获取 PATH 失败: %v", err)
		}

		// 检查新路径是否已经在 PATH 中
		currentPath := string(output)
		if !contains(currentPath, path) {
			// 添加新路径到 PATH，确保使用分号分隔
			currentPath = currentPath + ";" + path
			cmd = exec.Command("setx", "PATH", currentPath, "/M")
			if err := cmd.Run(); err != nil {
				log.Fatalf("添加到 PATH 失败: %v", err)
			}
			log.Printf("路径 %s 添加到 PATH 成功", path)
		} else {
			log.Printf("路径 %s 已在 PATH 中", path)
		}
		//如果是linux系统
	} else if runtime.GOOS == linux {
		// 获取当前用户信息
		usr, err := user.Current()
		if err != nil {
			log.Fatalf("获取用户信息失败: %v", err)
		}

		// 定义 bash 配置文件路径
		bashrcPath := filepath.Join(usr.HomeDir, ".bashrc")

		// 追加环境变量到 .bashrc
		file, err := os.OpenFile(bashrcPath, os.O_APPEND|os.O_WRONLY, 0600)
		if err != nil {
			log.Fatalf("打开文件 %s 失败: %v", bashrcPath, err)
		}
		defer file.Close()

		// 写入环境变量
		_, err = file.WriteString("\nexport " + "GOROOT" + "=" + path + "\n")
		if err != nil {
			log.Fatalf("写入文件失败: %v", err)
		}

		// 获取当前用户信息
		usr, err = user.Current()
		if err != nil {
			log.Fatalf("获取用户信息失败: %v", err)
		}

		// 定义 bash 配置文件路径
		bashrcPath = filepath.Join(usr.HomeDir, ".bashrc")

		// 追加 PATH 修改到 .bashrc
		file, err = os.OpenFile(bashrcPath, os.O_APPEND|os.O_WRONLY, 0600)
		if err != nil {
			log.Fatalf("打开文件 %s 失败: %v", bashrcPath, err)
		}
		defer file.Close()

		// 写入 PATH 修改
		_, err = file.WriteString("\nexport PATH=$PATH:" + path + "/bin" + "\n")
		if err != nil {
			log.Fatalf("写入文件失败: %v", err)
		}

	}
}

func contains(path string, target string) bool {
	// 检查路径是否存在
	for _, p := range strings.Split(path, ";") {
		if p == target {
			return true
		}
	}
	return false
}

// 检查路径是否符合 Windows 路径规范
func isValidWindowsPath(path string) {
	// 检查路径长度
	if len(path) > 260 {
		panic("非法路径: 长度超过 260 个字符")
	}

	// 检查非法字符，排除盘符中的冒号
	invalidChars := `<>:"|?*`
	for _, char := range invalidChars {
		// 允许盘符后面的冒号，如 C: 或 D:
		if char == ':' && len(path) > 1 && path[1] == ':' {
			continue
		}
		if strings.ContainsRune(path, char) {
			panic(fmt.Sprintf("非法路径: 包含非法字符 '%c'", char))
		}
	}

	// 检查绝对路径（以盘符开头）
	match, _ := regexp.MatchString(`^[A-Za-z]:[\\/]`, path)
	if !match {
		panic("非法路径: 必须以盘符开头")
	}

}

// 检查路径是否符合 Linux 路径规范
func isValidLinuxPath(path string) {
	// 检查路径长度（4096 是一个常见的限制，具体值视文件系统而定）
	if len(path) > 4096 {
		panic("非法路径")
	}

	// 检查路径是否包含 NULL 字符
	if strings.ContainsRune(path, '\x00') {
		panic("非法路径")
	}

	// 检查非法字符，Linux 路径中不能包含 /
	if strings.Contains(path, "/") && path != "/" {
		panic("非法路径")
	}

	// 检查路径的格式，不能以 / 开头并且不能只包含 /
	if path != "/" && strings.HasPrefix(path, "/") {
		// 使用正则表达式检查其他非法字符
		invalidChars := regexp.MustCompile(`[<>:"|?*]`)
		if invalidChars.MatchString(path) {
			panic("非法路径")
		}
	}

}

func installGo(goRoot string, fileName string) {
	//清空目标文件夹
	if err := clearDirectory(goRoot); err != nil {
		log.Printf("清空 %s 失败: %v", goRoot, err)
		return
	}
	//如果文件是zip压缩包
	if strings.HasSuffix(fileName, ".zip") {
		// 解压到指定目录
		if err := unzip(fileName, goRoot); err != nil {
			log.Printf("解压 %s 到 %s 失败: %v", fileName, goRoot, err)
			return
		}
	} else {
		// 解压到指定目录
		if err := untarGz(fileName, goRoot); err != nil {
			log.Printf("解压 %s 到 %s 失败: %v", fileName, goRoot, err)
			return
		}
	}

}

func clearDirectory(dir string) error {
	//如果目标文件夹存在
	if _, err := os.Stat(dir); err == nil {
		// 删除目标文件夹
		err := os.RemoveAll(dir)
		if err != nil {
			return err
		}
	}
	// 重新创建文件夹
	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}

func downloadFile(url string, fileName string) {
	log.Printf("开始下载最新版本")
	// 创建 HTTP 请求
	resp, err := http.Get(url)
	if err != nil {
		panic("创建下载请求失败: " + url + "错误信息: " + err.Error())
	}
	defer resp.Body.Close()

	// 检查 HTTP 响应状态
	if resp.StatusCode != http.StatusOK {
		panic("下载状态错误: " + url + "错误信息: " + err.Error())
	}

	// 创建文件
	out, err := os.Create(filepath.Join(".", fileName))
	if err != nil {
		panic("创建文件失败: " + fileName + "错误信息: " + err.Error())
	}
	defer out.Close()

	// 将响应的内容写入文件
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		panic("写入文件失败: " + fileName + "错误信息: " + err.Error())
	}
}

func getLatestGoVersion() (string, string, string) {
	log.Print("开始联网获取golang最新版本.")
	// 创建自定义的 TLS 配置
	tlsConfig := &tls.Config{
		MinVersion: tls.VersionTLS12, // 强制使用 TLS 1.2
	}

	// 创建 HTTP Transport
	transport := &http.Transport{
		TLSClientConfig: tlsConfig,
	}

	// 创建 HTTP 客户端
	client := &http.Client{
		Transport: transport,
		Timeout:   30 * time.Second, // 设置超时时间
	}
	// 发送请求
	resp, err := client.Get(latestVersionURL)
	if err != nil {
		log.Fatal("联网获取golang最新版本失败 , 错误信息: " + err.Error())
	}
	defer resp.Body.Close()
	// 将结果转成Response结构体
	data, err := io.ReadAll(resp.Body)

	var goVersions []GoVersion
	if err := json.Unmarshal(data, &goVersions); err != nil {
		log.Fatalf("Error parsing JSON: %v", err)
	}
	// 只获取最新版本
	version := goVersions[0]
	for _, file := range version.Files {
		//匹配操作系统和架构并且获取archive类型(压缩包类型 , 方便安装)
		if file.OS == runtime.GOOS && file.Arch == runtime.GOARCH && file.Kind == "archive" {
			downloadURL = downloadURL + file.Filename
			return version.Version, downloadURL, file.Filename
		}
	}
	panic("未找到合适的下载文件")

}

// 解析 go env 输出并填充结构体
func parseGoEnv(output string) (GoEnv, error) {
	lines := strings.Split(output, "\n")
	var env GoEnv

	fieldMap := map[string]*string{
		"GO111MODULE":    &env.GO111MODULE,
		"GOARCH":         &env.GOARCH,
		"GOBIN":          &env.GOBIN,
		"GOCACHE":        &env.GOCACHE,
		"GOENV":          &env.GOENV,
		"GOEXE":          &env.GOEXE,
		"GOEXPERIMENT":   &env.GOEXPERIMENT,
		"GOFLAGS":        &env.GOFLAGS,
		"GOHOSTARCH":     &env.GOHOSTARCH,
		"GOHOSTOS":       &env.GOHOSTOS,
		"GOINSECURE":     &env.GOINSECURE,
		"GOMODCACHE":     &env.GOMODCACHE,
		"GONOPROXY":      &env.GONOPROXY,
		"GONOSUMDB":      &env.GONOSUMDB,
		"GOOS":           &env.GOOS,
		"GOPATH":         &env.GOPATH,
		"GOPRIVATE":      &env.GOPRIVATE,
		"GOPROXY":        &env.GOPROXY,
		"GOROOT":         &env.GOROOT,
		"GOSUMDB":        &env.GOSUMDB,
		"GOTMPDIR":       &env.GOTMPDIR,
		"GOTOOLCHAIN":    &env.GOTOOLCHAIN,
		"GOTOOLDIR":      &env.GOTOOLDIR,
		"GOVCS":          &env.GOVCS,
		"GOVERSION":      &env.GOVERSION,
		"GODEBUG":        &env.GODEBUG,
		"GOTELEMETRY":    &env.GOTELEMETRY,
		"GOTELEMETRYDIR": &env.GOTELEMETRYDIR,
		"GCCGO":          &env.GCCGO,
		"GOAMD64":        &env.GOAMD64,
		"AR":             &env.AR,
		"CC":             &env.CC,
		"CXX":            &env.CXX,
		"CGO_ENABLED":    &env.CGO_ENABLED,
		"GOMOD":          &env.GOMOD,
		"GOWORK":         &env.GOWORK,
		"CGO_CFLAGS":     &env.CGO_CFLAGS,
		"CGO_CPPFLAGS":   &env.GOGCCFLAGS,
	}

	for _, line := range lines {
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		if strings.HasPrefix(key, "set ") {
			key = strings.TrimPrefix(key, "set ")
		}

		if ptr, ok := fieldMap[key]; ok {
			*ptr = value
		}
	}
	return env, nil
}

// unzip 解压 ZIP 文件到目标文件夹，并去掉第一层文件夹
// unzip 解压 ZIP 文件到目标文件夹，并去掉第一层文件夹
func unzip(src string, dest string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer r.Close()

	// 遍历 ZIP 中的每个文件/文件夹
	for _, f := range r.File {
		// 去除 ZIP 文件路径中的第一层文件夹
		fpath := filepath.Join(dest, trimFirstFolder(f.Name))

		// 检查是否为文件夹
		if f.FileInfo().IsDir() {
			// 创建文件夹
			if err := os.MkdirAll(fpath, os.ModePerm); err != nil {
				return err
			}
			continue
		}

		// 创建父文件夹
		if err = os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
			return err
		}

		// 打开 ZIP 文件内容
		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return err
		}

		rc, err := f.Open()
		if err != nil {
			return err
		}

		// 将文件内容复制到目标文件
		_, err = io.Copy(outFile, rc)

		// 关闭文件句柄
		outFile.Close()
		rc.Close()

		if err != nil {
			return err
		}
	}
	return nil
}

// trimFirstFolder 去掉路径中的第一层文件夹
func trimFirstFolder(path string) string {
	// 使用路径分隔符拆分路径
	parts := strings.SplitN(path, "/", 2) // Linux 使用 "/" 作为分隔符
	if len(parts) == 1 {
		parts = strings.SplitN(path, "\\", 2) // Windows 使用 "\\" 作为分隔符
	}
	if len(parts) > 1 {
		return parts[1] // 返回去掉第一层目录后的路径
	}
	return path // 如果没有找到分隔符，返回原路径
}

// 解压 .tar.gz 文件到指定文件夹
func untarGz(src string, dest string) error {
	// 打开 tar.gz 文件
	f, err := os.Open(src)
	if err != nil {
		return err
	}
	defer f.Close()

	// 使用 gzip 解压缩
	gzr, err := gzip.NewReader(f)
	if err != nil {
		return err
	}
	defer gzr.Close()

	// 创建 tar Reader
	tr := tar.NewReader(gzr)

	// 遍历 tar 文件中的每个文件/文件夹
	for {
		header, err := tr.Next()
		if err == io.EOF {
			// 读取到文件末尾，结束循环
			break
		}
		if err != nil {
			return err
		}

		// 处理文件路径，去掉第一层目录
		target := filepath.Join(dest, trimFirstFolder(header.Name))

		// 根据 header 类型进行不同的处理
		switch header.Typeflag {
		case tar.TypeDir:
			// 创建文件夹
			if err := os.MkdirAll(target, 0755); err != nil {
				return err
			}
		case tar.TypeReg:
			// 如果是文件，创建父文件夹
			if err := os.MkdirAll(filepath.Dir(target), 0755); err != nil {
				return err
			}

			// 创建文件
			outFile, err := os.OpenFile(target, os.O_CREATE|os.O_RDWR, os.FileMode(header.Mode))
			if err != nil {
				return err
			}
			defer outFile.Close()

			// 将内容写入文件
			if _, err := io.Copy(outFile, tr); err != nil {
				return err
			}
		default:
			fmt.Printf("忽略未知类型: %x 在 %s\n", header.Typeflag, header.Name)
		}
	}

	return nil
}
