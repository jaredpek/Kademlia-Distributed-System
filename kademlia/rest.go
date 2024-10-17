package kademlia

/*
import (
	"encoding/hex"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Rest struct {
	Kademlia *Kademlia
	Router   *gin.Engine
}

// newRest returns a new instance of the REST interface
func newRest(kademlia *Kademlia) *Rest {
	gin := gin.Default()

	rest := &Rest{}
	rest.Kademlia = kademlia
	rest.Router = gin

	return rest
}

// Starts the Rest server and which listens for HTTP requests.
func (r *Rest) StartServer(ip string) {
	r.Router.GET("/objects/:hash", r.GetObject)
	r.Router.POST("/objects", r.CreateObject)

	r.Router.Run(ip)
}

// Looks for an object in the kademlia network. An REST response is sent back with the result.
func (r *Rest) GetObject(c *gin.Context) {
	hash := c.Param("hash")

	// TODO: change to actually using Kademlia.LookupData
	res, err := r.Kademlia.Network.FindData(hash)

	if err != nil {
		fmt.Printf("Error in getObject: %s", err)
		c.IndentedJSON(http.StatusNotFound, err.Error())
	}

	c.IndentedJSON(http.StatusFound, res)
}

// Creates a new object in the kademlia network. If no error is returned a 201 REST response is sent back.
func (r *Rest) CreateObject(c *gin.Context) {
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
*/
