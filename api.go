package api

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/gin-gonic/contrib/ginrus"
	"github.com/gin-gonic/gin"
	"github.com/sohlich/elogrus"
	"gopkg.in/olivere/elastic.v3"
)

var (
	mysupersecretpassword = "unicornsAreAwesome"
)

// func main() {
// 	r := gin.Default()
//
// 	public := r.Group("/api")
//
// 	public.GET("/", func(c *gin.Context) {
// 		// Create the token
// 		token := jwt_lib.New(jwt_lib.GetSigningMethod("HS256"))
// 		// Set some claims
// 		token.Claims["ID"] = "Christopher"
// 		token.Claims["exp"] = time.Now().Add(time.Hour * 1).Unix()
// 		// Sign and get the complete encoded token as a string
// 		tokenString, err := token.SignedString([]byte(mysupersecretpassword))
// 		if err != nil {
// 			c.JSON(500, gin.H{"message": "Could not generate token"})
// 		}
// 		c.JSON(200, gin.H{"token": tokenString})
// 	})
//
// 	private := r.Group("/api/private")
// 	private.Use(jwt.Auth(mysupersecretpassword))
//
// 	/*
// 		Set this header in your request to get here.
// 		Authorization: Bearer `token`
// 	*/
//
// 	private.GET("/", func(c *gin.Context) {
// 		c.JSON(200, gin.H{"message": "Hello from private"})
// 	})
//
// 	r.Run("localhost:8080")
// }

func Run() {
	r := gin.New()

	response, err := http.Get("http://192.168.99.100:9200")
	if err != nil {
		fmt.Printf("%s", err)
		os.Exit(1)
	} else {
		defer response.Body.Close()
		contents, err := ioutil.ReadAll(response.Body)
		if err != nil {
			fmt.Printf("%s", err)
			os.Exit(1)
		}
		fmt.Printf("%s\n", string(contents))
	}

	// Add a ginrus middleware, which:
	//   - Logs all requests, like a combined access and error log.
	//   - Logs to stdout.
	//   - RFC3339 with UTC time format.
	r.Use(ginrus.Ginrus(logrus.StandardLogger(), time.RFC3339, true))

	// Add similar middleware, but:
	//   - Only logs requests with errors, like an error log.
	//   - Logs to stderr instead of stdout.
	//   - Local time zone instead of UTC.
	logger := logrus.New()
	client, err := elastic.NewClient(
		elastic.SetURL("http://192.168.99.100:9200"),
		elastic.SetMaxRetries(10),
		elastic.SetHealthcheck(false),
		elastic.SetSniff(false),
	)
	if err != nil {
		panic(err)
	}
	// Getting the ES version number is quite common, so there's a shortcut
	esversion, err := client.ElasticsearchVersion("http://192.168.99.100:9200")
	if err != nil {
		panic(err)
	}
	fmt.Printf("Elasticsearch version %s\n", esversion)
	if err != nil {
		logger.Panic(err)
	}
	hook, _ := elogrus.NewElasticHook(client, "malicehost", logrus.DebugLevel, "mylog")
	logger.Hooks.Add(hook)
	// logger.WithFields(logrus.Fields{
	// 	"name": "joe",
	// 	"age":  42,
	// }).Error("Hello world!")
	logger.Level = logrus.InfoLevel
	logger.Out = os.Stderr
	r.Use(ginrus.Ginrus(logger, time.RFC3339, false))

	// Example ping request.
	r.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong "+fmt.Sprint(time.Now().Unix()))
	})

	// Listen and Server in 0.0.0.0:8080
	r.Run(":8080")
}
