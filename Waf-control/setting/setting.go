package setting

import (
	"github.com/go-ini/ini"
	"log"
	"runtime"
	"strings"
	"path/filepath"
)

var (
    // App 设置
    AppName string
    RunMode string
    AppKey  string
    
    // 服务器设置
    HTTPPort int
    HTTPAddr string
    
    // 数据库设置
    DBType     string
    DBHost     string
    DBName     string
    DBUser     string
    DBPassword string
    DBSSLMode  string
    DBPath     string
    
    // API服务器
    APIServers []string
    
    // Nginx相关配置
    NginxBin    string
    NginxVhosts string
    RulePath    string
    
    // 操作系统类型
    IsWindows bool
)

// 加载配置
func Load() {
    // 检测操作系统类型
    IsWindows = runtime.GOOS == "windows"
    log.Printf("当前操作系统: %s, IsWindows: %v", runtime.GOOS, IsWindows)
    
    cfg, err := ini.Load("conf/app.ini")
    if err != nil {
        log.Fatalf("无法加载配置文件: %v", err)
    }
    
    // 加载App设置
    section, err := cfg.GetSection("app")
    if err != nil {
        log.Fatalf("无法获取 'app' 配置: %v", err)
    }
    AppName = section.Key("NAME").String()
    RunMode = section.Key("RUN_MODE").String()
    AppKey = section.Key("APP_KEY").String()
    
    // 加载服务器设置
    section, err = cfg.GetSection("server")
    if err != nil {
        log.Fatalf("无法获取 'server' 配置: %v", err)
    }
    HTTPPort = section.Key("HTTP_PORT").MustInt(5000)
    HTTPAddr = section.Key("HTTP_ADDR").String()
    
    // 尝试从 server 节加载 API 服务器设置
    apiServersStr := section.Key("API_SERVERS").String()
    if apiServersStr != "" {
        APIServers = strings.Split(apiServersStr, ",")
        for i, server := range APIServers {
            APIServers[i] = strings.TrimSpace(server)
        }
        log.Printf("从 [server] 节加载 API 服务器: %v", APIServers)
    } else {
        // 如果 server 节没有 API_SERVERS，尝试从 api 节加载
        apiSection, err := cfg.GetSection("api")
        if err == nil {
            servers := apiSection.Key("SERVERS").String()
            if servers != "" {
                APIServers = strings.Split(servers, ",")
                for i, server := range APIServers {
                    APIServers[i] = strings.TrimSpace(server)
                }
                log.Printf("从 [api] 节加载 API 服务器: %v", APIServers)
            }
        }
        
        // 如果两处都没有配置，使用默认值
        if len(APIServers) == 0 {
            log.Println("警告: 未找到 API 服务器配置，使用默认值 127.0.0.1")
            APIServers = []string{"127.0.0.1"}
        }
    }
    
    // 加载数据库设置
    section, err = cfg.GetSection("database")
    if err != nil {
        log.Fatalf("无法获取 'database' 配置: %v", err)
    }
    DBType = section.Key("TYPE").String()
    DBHost = section.Key("HOST").String()
    DBName = section.Key("NAME").String()
    DBUser = section.Key("USER").String()
    DBPassword = section.Key("PASSWORD").String()
    DBSSLMode = section.Key("SSL_MODE").String()
    DBPath = section.Key("PATH").String()
    
    // 如果是Windows系统，将数据库路径转换为Windows格式
    if IsWindows && DBType == "sqlite3" {
        DBPath = filepath.FromSlash(DBPath)
    }
    
    // 加载Nginx相关配置
    
    nginxSection, err := cfg.GetSection("nginx")
    if err == nil {
        // 默认使用 Linux 配置
        NginxBin = nginxSection.Key("bin").MustString("/usr/local/openresty/nginx/sbin/nginx")
        NginxVhosts = nginxSection.Key("vhosts").MustString("/usr/local/openresty/nginx/conf/vhosts")
        RulePath = nginxSection.Key("rules").MustString("/usr/local/openresty/nginx/conf/x-waf/rules")
        
        // 如果是 Windows 系统，则使用 Windows 配置
        if IsWindows {
            winSection, err := cfg.GetSection("nginx.windows")
            if err == nil {
                NginxBin = winSection.Key("bin").MustString("C:\\Program Files\\OpenResty\\nginx\\nginx.exe")
                NginxVhosts = winSection.Key("vhosts").MustString("C:\\Program Files\\OpenResty\\nginx\\conf\\vhosts")
                RulePath = winSection.Key("rules").MustString("C:\\Program Files\\OpenResty\\nginx\\conf\\x-waf\\rules")
                
                // 确保路径使用 Windows 格式
                NginxBin = filepath.FromSlash(NginxBin)
                NginxVhosts = filepath.FromSlash(NginxVhosts)
                RulePath = filepath.FromSlash(RulePath)
                
                log.Printf("使用 Windows Nginx 配置: bin=%s, vhosts=%s, rules=%s", NginxBin, NginxVhosts, RulePath)
            } else {
                log.Printf("警告: 无法获取 'nginx.windows' 配置节，将尝试使用 Linux 路径: %v", err)
            }
        }
    } else {
        // 如果没有 nginx 节，使用默认值
        log.Printf("警告: 无法获取 'nginx' 配置节，将使用默认值: %v", err)
    
        NginxBin = "/usr/local/openresty/nginx/sbin/nginx"
        NginxVhosts = "/usr/local/openresty/nginx/conf/vhosts"
        RulePath = "/usr/local/openresty/nginx/conf/x-waf/rules"
    }
    
    log.Printf("已加载Nginx配置: bin=%s, vhosts=%s, rules=%s", NginxBin, NginxVhosts, RulePath)
    log.Printf("注意: 这些配置将应用于远程 WAF 服务器 %v，而非本地系统", APIServers)
}

// 获取平台相关的路径
func GetPlatformPath(path string) string {
    if IsWindows {
        return filepath.FromSlash(path)
    }
    return path
}
