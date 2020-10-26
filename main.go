package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	// "github.com/hashicorp/go-multierror"
)

// StatusResponse is a type which used here as a general response for api calls.
type StatusResponse struct {
	Stat string      `json:"status"`
	Code int         `json:"code"`
	Resp interface{} `json:"resp"`
}

/*
responseBuilder will return a typical, generalized status response.

statuscode -> please use net/http constants for this, as we're also using http.StatusText.

response -> set as interface{}, anything will go. It will only crash when json.Marshal is failed.
*/
func responseBuilder(statuscode int, response interface{}) (sr StatusResponse) {
	res, err := json.Marshal(response)
	stat := &statuscode
	defer func() {
		if r := recover(); r != nil {
			*stat = http.StatusInternalServerError
			res = []byte(`{"message": "Recover from responseBuilder failure"}`)
			sr = StatusResponse{
				Stat: http.StatusText(statuscode),
				Code: statuscode,
				Resp: json.RawMessage(res),
			}
		}
	}()
	if err != nil {
		panic("responseBuilder failed at json.marshal")
	}
	if http.StatusText(statuscode) == "" {
		panic(fmt.Sprintf("Undefined statuscode '%v'", statuscode))
	}
	sr = StatusResponse{
		Stat: http.StatusText(statuscode),
		Code: statuscode,
		Resp: json.RawMessage(res),
	}
	return sr
}

func main() {
	r := gin.Default()
	r.GET("/env", func(c *gin.Context) {
		env := os.Environ()
		// a, _ := json.Marshal(env)
		// r := &StatusResponse{
		// 	Stat: http.StatusText(http.StatusOK),
		// 	Code: http.StatusOK,
		// 	Resp: json.RawMessage(a),
		// }
		r := responseBuilder(123, env)
		c.JSON(r.Code, r)
	})
	r.GET("/panic", func(c *gin.Context) {
		panic("this")
	})
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	r.Run(":1500")
}
