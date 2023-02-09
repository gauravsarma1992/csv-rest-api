package csvrestapi

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/gin-gonic/gin"
)

type (
	Server struct {
		apiEngine *gin.Engine
		Config    *Config
	}

	Config struct {
		Server struct {
			Host string `json:"host"`
			Port string `json:"port"`
		} `json:"server"`
	}
)

func New() (srv *Server, err error) {
	srv = &Server{
		apiEngine: gin.Default(),
		Config:    &Config{},
	}
	if err = srv.setConfig(); err != nil {
		return
	}
	if err = srv.setRoutes(); err != nil {
		return
	}
	return
}

func (srv *Server) setConfig() (err error) {
	var (
		contB []byte
	)
	if contB, err = ioutil.ReadFile("./config.json"); err != nil {
		return
	}
	if err = json.Unmarshal(contB, srv.Config); err != nil {
		return
	}
	return
}

func (srv *Server) setRoutes() (err error) {
	srv.apiEngine.GET("/ping", srv.PingHandler)
	srv.apiEngine.POST("/csv", srv.CsvHandler)
	return
}

func (srv *Server) PingHandler(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

func (srv *Server) CsvHandler(c *gin.Context) {
	type (
		CsvRequests struct {
			Files   []string          `json:"files"`
			Filters map[string]string `json:"filters"`
		}
	)
	var (
		csvFinder *CsvFinder
		csvReq    CsvRequests
		err       error
		resp      []interface{}
	)
	if err = c.ShouldBindJSON(&csvReq); err != nil {
		c.JSON(400, gin.H{
			"error": err,
		})
		return
	}

	if csvFinder, err = NewCsvReader(csvReq.Files, csvReq.Filters); err != nil {
		c.JSON(500, gin.H{
			"error": err,
		})
		return
	}
	if resp, err = csvFinder.Search(); err != nil {
		c.JSON(500, gin.H{
			"error": err,
		})
		return
	}

	c.JSON(200, gin.H{
		"length":        len(resp),
		"matched_elems": resp,
	})
}

func (srv *Server) Run() (err error) {
	log.Println("Running REST Server")
	srv.apiEngine.Run(fmt.Sprintf("%s:%s", srv.Config.Server.Host, srv.Config.Server.Port))
	return
}
