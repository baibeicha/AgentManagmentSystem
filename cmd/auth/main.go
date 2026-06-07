package main

import (
	"AgentManagmentSystem/pkg/config"
)

func main() {
	cfg := config.MustLoad("test", "cfg-path", "./cmd/auth")
}
