package kademlia

import (
	"encoding/hex"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Rest struct {
	Kademlia *Kademlia
	Router   *gin.Engine
}

func newRest(kademlia *Kademlia) *Rest {
	gin := gin.Default()

	rest := &Rest{}
	rest.Kademlia = kademlia
	rest.Router = gin

	return rest
}

func (r *Rest) startServer(ip string) {
	r.Router.GET("/objects/:hash", r.getObject)
	r.Router.POST("/objects", r.createObject)

	r.Router.Run(ip)
}

func (r *Rest) getObject(c *gin.Context) {
	hash := c.Param("hash")

	res := r.Kademlia.LookupData(hash)

	c.IndentedJSON(http.StatusNotFound, res)
}

func (r *Rest) createObject(c *gin.Context) {
	type inputData struct {
		Data string `json:"data"`
	}

	var data inputData

	err := c.BindJSON(&data)

	if err != nil {
		log.Fatalf("Error in createObject: %s", err)
	}

	d, _ := hex.DecodeString(data.Data)

	err, _ = r.Kademlia.Store(d)

	if err != nil {
		log.Fatalf("Error in Store: %s", err)
	}

	c.IndentedJSON(http.StatusCreated, data)
}
