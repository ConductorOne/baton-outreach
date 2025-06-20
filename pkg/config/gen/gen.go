package main

import (
	cfg "github.com/conductorone/baton-outreach/pkg/config"
	"github.com/conductorone/baton-sdk/pkg/config"
)

func main() {
	config.Generate("outreach", cfg.Config)
}
