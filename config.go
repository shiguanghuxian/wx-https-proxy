package main

// Config 配置文件
type Config struct {
	ListenAddress string `toml:"listen_address"`   // 监听地址
	Cert    *CertConfig          `toml:"cert"`   // 证书路径配置
	Servers []*ProxyServerConfig `toml:"server"` // 要代理的服务端完整地址(到端口)
}

// CertConfig 证书路径配置
type CertConfig struct {
	Crt string `toml:"crt"`
	Key string `toml:"key"`
}

// ProxyServerConfig 被代理服务地址配置
type ProxyServerConfig struct {
	Prefix   string `toml:"prefix"`   // 请求地址前缀
	Protocol string `toml:"protocol"` // http | https
	Default  bool   `toml:"default"`  // 是否是默认代理地址
	Address  string `toml:"address"`  // 被代理服务器地址
	Port     int    `toml:"port"`     // 被代理服务器端口
}
