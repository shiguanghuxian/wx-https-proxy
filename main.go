package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/google/logger"
	"github.com/naoina/toml"
)

var logPath = "./logs/wx-https-proxy.log"
var lf *os.File
var cfg *Config

func main() {
	// 系统日志显示文件和行号
	log.SetFlags(log.Lshortfile | log.LstdFlags)

	// 初始化日志
	initLog()
	// 初始化配置文件
	initConfig()

	proxy := &Proxy{
		Servers: cfg.Servers,
	}

	err := http.ListenAndServeTLS(cfg.ListenAddress, cfg.Cert.Crt, cfg.Cert.Key, proxy)
	if err != nil {
		logger.Fatalln("服务监听错误", err)
	}

	// 接收控制信号
	signalChan := make(chan os.Signal, 1)                      // 创建一个信号量的chan
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM) // 让进程收集信号量
	select {
	case <-signalChan:
		if lf != nil {
			lf.Close()
		}
	}
}

// 初始化日志
func initLog() {
	var err error
	lf, err = os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0660)
	if err != nil {
		logger.Fatalf("打开日志文件错误: %v", err)
	}
	logger.Init("wx-https-proxy", true, true, lf)
}

// 初始化配置文件
func initConfig() {
	f, err := os.Open("./config.toml")
	if err != nil {
		logger.Fatalf("配置文件读取错误: %v", err)
	}
	defer f.Close()
	cfg = new(Config)
	if err := toml.NewDecoder(f).Decode(cfg); err != nil {
		logger.Fatalf("日志文件解析错误: %v", err)
	}
}

// Proxy 代理处理对象
type Proxy struct {
	Servers []*ProxyServerConfig
}

// 代理实现
func (that *Proxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// 查找前缀对应的服务器
	currentKey := -1
	defaultKey := -1
	for k, s := range that.Servers {
		if strings.Index(r.URL.Path, "/"+s.Prefix) == 0 {
			currentKey = k
			break
		}
		if s.Default == true {
			defaultKey = k
		}
	}
	if currentKey == -1 && defaultKey == -1 {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("代理未找的对方服务器"))
		return
	}
	key := currentKey
	if key == -1 {
		key = defaultKey
	}
	serverURL := fmt.Sprintf("%s://%s:%d", that.Servers[key].Protocol, that.Servers[key].Address, that.Servers[key].Port)
	remote, err := url.Parse(serverURL)
	if err != nil {
		logger.Errorln(err)
	}
	if len(r.URL.Path) < len(that.Servers[key].Prefix)+1 {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("地址不能为空"))
		return
	}
	// 处理前缀
	r.URL.Path = r.URL.Path[len(that.Servers[key].Prefix)+1:]
	logger.Infoln(serverURL)
	logger.Infoln(r.URL.Path)
	proxy := httputil.NewSingleHostReverseProxy(remote)
	proxy.ServeHTTP(w, r)
}
