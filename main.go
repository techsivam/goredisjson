package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/nitishm/go-rejson/v4"
)

/* type Data struct {
	Key   string      `json:"key"`
	Value interface{} `json:"value"`
} */

var rh *rejson.Handler

func init() {
	var rdb = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	rh = rejson.NewReJSONHandler()
	rh.SetGoRedisClient(rdb)
}

func GetRedis(c *gin.Context) {
	tenant := c.Param("tenant")

	// Retrieve the JSON value from Redis as bytes
	jsonBytes, err := rh.JSONGet(tenant, ".")
	if err != nil {
		c.JSON(404, gin.H{"error": "Key not found"})
		return
	}

	// Perform a type assertion to convert the bytes to []byte
	bytes, ok := jsonBytes.([]byte)
	if !ok {
		c.JSON(500, gin.H{"error": "Invalid data type"})
		return
	}

	// Unmarshal the JSON bytes into a map[string]interface{}
	var jsonData map[string]interface{}
	err = json.Unmarshal(bytes, &jsonData)
	if err != nil {
		c.JSON(500, gin.H{"error": "Error decoding JSON"})
		return
	}
	fmt.Println("JSON READ: ", jsonData)
	// Return the JSON data as a response
	c.JSON(200, jsonData)
}

func PutRedis(c *gin.Context) {
	tenant := c.Param("tenant")
	var jsonData map[string]interface{}

	if err := c.ShouldBindJSON(&jsonData); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	redisKey := tenant
	_, err := rh.JSONSet(redisKey, ".", jsonData)

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	fmt.Println("JSON WRITE: ", jsonData)
	c.JSON(200, gin.H{"status": "success"})
	/*file, _, err := c.Request.FormFile("file")
		if err != nil {
			c.JSON(400, gin.H{"error": "file not found"})
			return
		}

		content, err := ioutil.ReadAll(file)
		if err != nil {
			c.JSON(400, gin.H{"error": "error reading file"})
			return
		}

		// Unmarshal the JSON bytes into a map[string]interface{}
	 var jsonData map[string]interface{}
		err = json.Unmarshal(content, &jsonData)
		if err != nil {
			c.JSON(400, gin.H{"error": "error decoding JSON"})
			return
		}

		_, err = rh.JSONSet(tenant, ".", jsonData)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		fmt.Println("JSON WRITE: ", jsonData)
		c.JSON(200, gin.H{"status": "success"}) */
}
func PutRedis1(c *gin.Context) {
	tenant := c.Param("tenant")
	file, _, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(400, gin.H{"error": "file not found"})
		return
	}

	content, err := ioutil.ReadAll(file)
	if err != nil {
		c.JSON(400, gin.H{"error": "error reading file"})
		return
	}

	// Unmarshal the JSON bytes into a map[string]interface{}
	var jsonData map[string]interface{}
	err = json.Unmarshal(content, &jsonData)
	if err != nil {
		c.JSON(400, gin.H{"error": "error decoding JSON"})
		return
	}

	_, err = rh.JSONSet(tenant, ".", jsonData)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	fmt.Println("JSON WRITE: ", jsonData)
	c.JSON(200, gin.H{"status": "success"})
}

func main() {
	router := gin.Default()
	router.GET("/:tenant", GetRedis)
	router.POST("/:tenant", PutRedis)

	router.Run(":8082")
}
