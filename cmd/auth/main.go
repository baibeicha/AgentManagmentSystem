package main

import (
	"AgentManagmentSystem/pkg/config"
)

func main() {
	_ = config.MustLoad("auth-config")
}
