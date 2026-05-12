package main

import "github.com/antongolenev23/voltake-services/services/auth/internal/config"

func main() {
	cfg := config.MustLoad()
	_ = cfg
}
