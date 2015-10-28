package main

import (
	"path"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/kelseyhightower/envconfig"
	"github.com/rrawrriw/gin-static"
)

const (
	SessionColl = "Session"
)

type (
	Specs struct {
		Host      string `envconfig:"host"`
		Port      int    `envconfig:"port"`
		PublicDir string `envconfig:"public_dir"`
	}
)

func main() {
	specs := Specs{}
	envconfig.Process("LIFE", &specs)
	host := specs.Host
	port := specs.Port
	srvRes := host + ":" + strconv.Itoa(port)

	publicDir := specs.PublicDir
	htmlDir := path.Join(publicDir, "html")

	srv := gin.Default()
	srv.Use(ginstatic.Serve("/", ginstatic.LocalFile(htmlDir, false)))
	srv.Static("/public", publicDir)

	srv.Run(srvRes)
}
