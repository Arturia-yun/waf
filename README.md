# WAF 管理系统 / WAF Management System

[English Version](#waf-management-system)

## 简介 

本项目是一个基于 Go 语言开发的 Web 应用防火墙（WAF）管理系统，配合开源 WAF 本体（如 x-waf），实现了规则的可视化管理、站点配置、自动同步等功能，适用于多台 WAF 主机的集中化运维和安全策略统一管理。

由两部分组成：

- Waf-control : 管理控制台，提供 Web 界面进行规则和站点管理
- x-waf : 基于 OpenResty 的 WAF 引擎，负责实际的请求过滤和攻击拦截



---

## 功能特性 

- **规则管理** ：支持多种类型规则的在线管理，包括 URL、参数、IP、User-Agent、Cookie 等
- **站点管理** ：支持添加、编辑、删除站点配置
- **多服务器同步** ：支持将规则和配置自动同步到多台 WAF 服务器
- **用户管**理 ：支持多用户管理，权限控制
- **实时监控** ：查看 WAF 拦截日志和统计信息
- **响应式界面**：现代化响应式前端界面
- **部署支持**：支持 Windows 和 Linux 环境部署



---

## 支持拦截的攻击行为 ##
**WAF 系统能够有效拦截以下类型的攻击：**
- SQL 注入攻击
- XSS 跨站脚本攻击
- 命令注入攻击
- 文件包含漏洞利用
- 敏感文件访问
- 恶意爬虫和扫描器
- CC 攻击（Challenge Collapsar）
- 恶意 IP 访问
- 异常 User-Agent
- 恶意 Cookie
- 异常 POST 数据
- 敏感 Header 信息

---

## 目录结构 

```bash
Waf-control/         # 管理端主程序
├─ conf/           # 配置文件目录（如 app.ini）
├─ models/         # ORM 数据模型
├─ modules/        # 功能模块
├─ routers/        # 路由与控制器
├─ setting/        # 配置加载
├─ templates/      # 前端模板与静态资源
├─ server.go       # 程序入口
 x-waf/               # WAF 本体（Lua 规则引擎等）
```
---
 
## 安装部署 ##
**依赖环境**

- Go 1.16+
- MySQL 5.7+ 或 SQLite3
- OpenResty 1.19.3.1+
- Nginx

## ubuntu平台安装 ##

### 编译安装openresty ###

```bash
apt-get install libreadline-dev libncurses5-dev libpcre3-dev libssl-dev perl make build-essential
sudo ln -s /sbin/ldconfig /usr/bin/ldconfig
wget https://openresty.org/download/openresty-1.9.15.1.tar.gz
tar -zxvf openresty-1.9.15.1.tar.gz
cd openresty-1.9.15.1
make && sudo make install
```
安装完成后，在/etc/profile中加入openresty的环境变量：
```bash
export PATH=/usr/local/openresty/bin:$PATH
```
### 部署WAF。 ###

将WAF代码目录放置在/usr/local/openresty/nginx/conf目录下，然后在Openresty的conf目录下新建Vhost目录，配置如下

```bash
cd /usr/local/openresty/nginx/conf
git clone it clone https://github.com/Arturia-yun/Waf/x-waf
mkdir -p /usr/local/openresty/nginx/conf/Vhost
```

用x-waf中的nginx.conf文件覆盖原有的nginx.conf文件

### 安装waf管理后台Waf-control ###

#### 二进制安装(推荐) #### 

直接从release获取对应操作系统的二进制版本

#### 源码安装 ####

- 安装依赖
```bash
go get gopkg.in/macaron.v1
go get gopkg.in/ini.v1
go get github.com/go-sql-driver/mysql
go get github.com/go-xorm/xorm
```
### 访问管理界面 ###
http://localhost:8081
默认用户名: admin
默认密码: admin


## WAF 拦截效果展示 ##


---

## 鸣谢 ##

本项目的 WAF 引擎基于 [unixhot](https://github.com/unixhot/waf) 开源项目，在此特别感谢 unixhot 的杰出贡献。同时感谢以下开源项目：

- Macaron
- XORM
- OpenResty

---

## 反馈与建议 ##

本项目持续不断地在改进中。如果您遇到任何问题或对改进有建议，请随时提供反馈。感谢您对本项目的支持和理解！

## 许可证
本项目采用 MIT 许可证



# WAF Management System

## Introduction

This project is a Web Application Firewall (WAF) management system developed in Go, designed to work with open-source WAF engines (such as x-waf). It provides visual rule management, site configuration, automatic synchronization, and is suitable for centralized operation and unified security policy management across multiple WAF hosts.

The system consists of two parts:

- **Waf-control**: Management console providing a web interface for rule and site management
- **x-waf**: OpenResty-based WAF engine responsible for actual request filtering and attack blocking

---

## Features

- **Rule Management**: Online management of multiple rule types, including URL, parameters, IP, User-Agent, Cookie, etc.
- **Site Management**: Add, edit, and delete site configurations
- **Multi-server Synchronization**: Automatically sync rules and configurations to multiple WAF servers
- **User Management**: Support for multiple users with permission control
- **Real-time Monitoring**: View WAF blocking logs and statistics
- **Responsive Interface**: Modern user interface supporting mobile devices
- **Deployment Support**: Supports both Windows and Linux environments

---

## Supported Attack Blocking

The WAF system effectively blocks the following types of attacks:

- SQL Injection Attacks
- XSS (Cross-Site Scripting) Attacks
- Command Injection Attacks
- File Inclusion Vulnerability Exploitation
- Sensitive File Access
- Malicious Crawlers and Scanners
- CC (Challenge Collapsar) Attacks
- Malicious IP Access
- Abnormal User-Agents
- Malicious Cookies
- Abnormal POST Data
- Sensitive Header Information

---

## Directory Structure

```bash
Waf-control/         # Main program for management
├─ conf/           # Configuration files (e.g., app.ini)
├─ models/         # ORM data models
├─ modules/        # Functional modules
├─ routers/        # Routes and controllers
├─ setting/        # Configuration loading
├─ templates/      # Frontend templates and static resources
├─ server.go       # Program entry point
x-waf/               # WAF engine (Lua rule engine, etc.)
```

---

## Installation and Deployment

**Dependencies**

- Go 1.16+
- MySQL 5.7+ or SQLite3
- OpenResty 1.19.3.1+
- Nginx

## Installation on Ubuntu

### Compile and Install OpenResty

```bash
apt-get install libreadline-dev libncurses5-dev libpcre3-dev libssl-dev perl make build-essential
sudo ln -s /sbin/ldconfig /usr/bin/ldconfig
wget https://openresty.org/download/openresty-1.9.15.1.tar.gz
tar -zxvf openresty-1.9.15.1.tar.gz
cd openresty-1.9.15.1
make && sudo make install
```
After installation, add OpenResty to the environment variables in /etc/profile:

```bash
export PATH=/usr/local/openresty/bin:$PATH
```

### Deploy WAF ###

Place the WAF code directory in /usr/local/openresty/nginx/conf, then create a Vhost directory in the OpenResty conf directory with the following configuration:
```bash
cd /usr/local/openresty/nginx/conf
git clone https://github.com/Arturia-yun/Waf/x-waf
mkdir -p /usr/local/openresty/nginx/conf/Vhost
```
Replace the original nginx.conf file with the one from x-waf.

### Install Waf-control Management Backend ###

#### Binary Installation (Recommended) ####

Directly obtain the binary version for the corresponding operating system from the release.

#### Source Installation ####

- Install dependencies
```bash
go get gopkg.in/macaron.v1
go get gopkg.in/ini.v1
go get github.com/go-sql-driver/mysql
go get github.com/go-xorm/xorm
```

### Access the Management Interface ###

http://localhost:8081
Default username: admin
Default password: admin

---

## WAF Blocking Effect Showcase 

---

## Acknowledgements

The WAF engine in this project is based on the [unixhot](https://github.com/unixhot/waf) open-source project. Special thanks to unixhot for their outstanding contribution. Also thanks to the following open-source projects:
- Macaron
- XORM
- OpenResty

## Feedback
This project is continuously being improved. If you encounter any issues or have suggestions for improvement, please feel free to provide feedback. Thank you for your support and understanding!

## License
This project is licensed under the MIT License

