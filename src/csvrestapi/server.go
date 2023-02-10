package csvrestapi

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"io/ioutil"
	"log"
	"strings"

	"github.com/gin-contrib/cors"
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
	srv.apiEngine.Use(cors.Default())
	srv.apiEngine.GET("/ping", srv.PingHandler)
	srv.apiEngine.POST("/csv", srv.CsvHandler)
	srv.apiEngine.GET("/csv_folder", srv.CsvFolderReadHandler)
	return
}

func (srv *Server) PingHandler(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

func (srv *Server) CsvFolderReadHandler(c *gin.Context) {
	var (
		folderName      string
		err             error
		folderFileNames []string
		folderFiles     []fs.FileInfo
	)
	folderName = c.Query("folder_name")
	log.Println("here", folderName)
	if folderName == "" {
		c.JSON(400, gin.H{
			"error": "Folder name cannot be empty",
		})
		return
	}
	if folderFiles, err = ioutil.ReadDir(folderName); err != nil {
		c.JSON(400, gin.H{
			"error": "No files present in folder",
		})
		return
	}
	for _, file := range folderFiles {
		fileName := file.Name()
		if !strings.HasSuffix(fileName, ".csv") {
			continue
		}
		folderFileNames = append(folderFileNames, fileName)
	}

	c.JSON(200, gin.H{
		"folder_files": folderFileNames,
	})
	return
}

func (srv *Server) CsvHandler(c *gin.Context) {
	type (
		CsvRequests struct {
			Files   []string          `json:"files"`
			Folders []string          `json:"folders"`
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

	if csvFinder, err = NewCsvReader(csvReq.Files, csvReq.Folders, csvReq.Filters); err != nil {
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
	return
}

func (srv *Server) Run() (err error) {
	log.Println("Running REST Server")
	srv.apiEngine.Run(fmt.Sprintf("%s:%s", srv.Config.Server.Host, srv.Config.Server.Port))
	return
}
