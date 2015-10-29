package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
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
		Data      string `envconfig:"data"`
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
	srv.GET("/data", specWrap(specs))

	srv.Run(srvRes)
}

func specWrap(specs Specs) gin.HandlerFunc {
	return func(c *gin.Context) {
		str, err := ioutil.ReadFile(specs.Data)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"Err": err.Error()})
		}

		r := gin.H{}
		err = json.Unmarshal([]byte(str), &r)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"Err": err.Error()})
		}

		c.JSON(http.StatusOK, r)
	}
}
