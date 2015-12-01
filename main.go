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
	"github.com/rrawrriw/life-ctrl"
)

const (
	SessionColl = "Session"
)

type (
	Specs struct {
		Host      string `envconfig:"host"`
		Port      int    `envconfig:"port"`
		PublicDir string `envconfig:"public_dir"`
		StagesDir string `envconfig:"stages_dir"`
		PagesDir  string `envconfig:"pages_dir"`
		ConfigDir string `envconfig:"config_dir"`
	}
)

func main() {
	specs := Specs{}
	envconfig.Process("LIFE", &specs)
	host := specs.Host
	port := specs.Port
	srvRes := host + ":" + strconv.Itoa(port)

	publicDir := specs.PublicDir
	configDir := specs.ConfigDir
	htmlDir := path.Join(publicDir, "html")

	srv := gin.Default()
	srv.Use(ginstatic.Serve("/", ginstatic.LocalFile(htmlDir, false)))
	srv.Static("/public", publicDir)
	srv.Static("/config", configDir)
	srv.GET("/data", specWrap(specs, readStages))
	srv.GET("/page/:name", specWrap(specs, readMdFile))

	srv.Run(srvRes)
}

func readStages(c *gin.Context, specs Specs) {
	j, err := lifectrl.NewStageJSON(specs.StagesDir)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"Err": err.Error()})
		return
	}

	r := gin.H{}
	err = json.Unmarshal(j, &r)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"Err": err.Error()})
		return
	}

	c.JSON(http.StatusOK, r)
}

func readMdFile(c *gin.Context, specs Specs) {
	name := c.Params.ByName("name")
	str, err := ioutil.ReadFile(path.Join(specs.PagesDir, name+".md"))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"Err": err.Error()})
		return
	}

	m := map[string]string{}
	m[name] = string(str)

	c.JSON(http.StatusOK, m)
}

func specWrap(specs Specs, h func(*gin.Context, Specs)) gin.HandlerFunc {
	return func(c *gin.Context) {
		h(c, specs)
	}
}
