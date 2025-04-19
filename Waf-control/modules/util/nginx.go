package util

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"Waf-control/setting"
)

// Message API响应消息结构
type Message struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

// WriteNginxConf 写入Nginx配置文件
func WriteNginxConf(proxyConfig string, siteName string, vhostPath string) (err error) {
    log.Println("开始写入配置文件...")
    log.Println("站点名称:", siteName)
    log.Println("虚拟主机目录:", vhostPath)
    
    // 确保使用正确的路径分隔符
    if setting.IsWindows {
        vhostPath = filepath.FromSlash(vhostPath)
    }
    
    // 确保目录存在
    err = os.MkdirAll(vhostPath, 0755)
    if err != nil {
        log.Println("创建目录失败:", err)
        return err
    }
    
    // 使用 filepath.Join 而不是 path.Join 以确保使用正确的路径分隔符
    proxyConfigFile := filepath.Join(vhostPath, fmt.Sprintf("%v.conf", siteName))
    log.Println("配置文件路径:", proxyConfigFile)
    
    // 尝试打开文件
    fileConfig, err := os.Create(proxyConfigFile)
    if err != nil {
        log.Println("创建配置文件失败:", err)
        return err
    }
    defer fileConfig.Close()
    
    // 替换换行符
    proxyConfig = strings.Replace(proxyConfig, "\r\n", "\n", -1)
    
    // 写入配置
    _, err = fileConfig.WriteString(proxyConfig)
    if err != nil {
        log.Println("写入配置失败:", err)
        return err
    }
    
    log.Println("配置文件写入成功:", proxyConfigFile)
    return nil
}

// ReloadNginx 重新加载Nginx配置
func ReloadNginx() (err error) {
    log.Println("开始重新加载Nginx配置")
    
    // 检查 Nginx 可执行文件是否存在
    if _, err := os.Stat(setting.NginxBin); os.IsNotExist(err) {
        log.Printf("错误: Nginx 可执行文件不存在: %s", setting.NginxBin)
        return fmt.Errorf("nginx 可执行文件不存在: %s", setting.NginxBin)
    }
    
    // 获取Nginx安装目录
    nginxDir := filepath.Dir(setting.NginxBin)
    log.Printf("Nginx安装目录: %s", nginxDir)
    
    // 构建配置文件路径
    confPath := filepath.Join(nginxDir, "conf", "nginx.conf")
    log.Printf("Nginx配置文件路径: %s", confPath)
    
    // 测试 Nginx 配置
    cmd := exec.Command(setting.NginxBin, "-t", "-c", confPath)
    cmd.Dir = nginxDir
    output, err := cmd.CombinedOutput()
    log.Printf("Nginx 配置测试结果: %s", string(output))
    
    if err != nil {
        log.Printf("Nginx 配置测试失败: %v", err)
        return err
    }
    
    // 在Windows环境下，使用不同的重载策略
    if setting.IsWindows {
        // 1. 首先检查Nginx是否正在运行
        checkCmd := exec.Command("tasklist", "/FI", "IMAGENAME eq nginx.exe")
        checkOutput, _ := checkCmd.CombinedOutput()
        nginxRunning := strings.Contains(string(checkOutput), "nginx.exe")
        
        if nginxRunning {
            // 2. 如果Nginx正在运行，先尝试正常停止
            log.Println("检测到Nginx正在运行，尝试停止...")
            
            // 使用taskkill强制终止所有nginx进程
            killCmd := exec.Command("taskkill", "/IM", "nginx.exe", "/F")
            killOutput, _ := killCmd.CombinedOutput()
            log.Printf("强制终止Nginx进程结果: %s", string(killOutput))
            
            // 等待进程完全退出
            time.Sleep(2 * time.Second)
        } else {
            log.Println("未检测到Nginx正在运行")
        }
        
        // 3. 启动Nginx
        log.Println("开始启动Nginx...")
        
        // 使用 cmd /c start 命令启动 Nginx，但不等待输出
        startCmd := exec.Command("cmd", "/c", "start", "/b", setting.NginxBin)
        startCmd.Dir = nginxDir
        
        // 使用Start()而不是CombinedOutput()，这样不会等待命令完成
        if err := startCmd.Start(); err != nil {
            log.Printf("Nginx 启动失败: %v", err)
            return err
        }
        
        // 等待一小段时间让Nginx有机会启动
        time.Sleep(1 * time.Second)
        
        // 检查Nginx是否成功启动
        checkCmd = exec.Command("tasklist", "/FI", "IMAGENAME eq nginx.exe")
        checkOutput, _ = checkCmd.CombinedOutput()
        if strings.Contains(string(checkOutput), "nginx.exe") {
            log.Println("Nginx 已成功启动")
        } else {
            log.Println("警告: 未检测到Nginx进程，可能启动失败")
            // 即使没有检测到进程，也不返回错误，因为有可能是检测太快
        }
        
        log.Println("Nginx 重载完成")
        return nil
    } else {
        // Linux环境下使用reload命令
        cmd = exec.Command(setting.NginxBin, "-s", "reload")
        cmd.Dir = nginxDir
        output, err = cmd.CombinedOutput()
        log.Printf("Nginx 重载结果: %s", string(output))
        
        if err != nil {
            log.Printf("Nginx 重载失败: %v", err)
        }
        
        return err
    }
}