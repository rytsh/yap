package config

import "github.com/rytsh/yap/internal/tui"

var Application = struct {
	LogLevel string     `cfg:"log-level"`
	Server   Server     `cfg:"server"`
	Screen   tui.Screen `cfg:"screen"`
}{
	LogLevel: "info",
	Server: Server{
		Host: "0.0.0.0",
		Port: 2222,
	},
}

type Server struct {
	Host string
	Port int
}
