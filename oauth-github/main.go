package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type GhResponse struct {
	Access_token string `json:"access_token"`
	Scope        string `json:"scope"`
	Token_type   string `json:"token_type"`
}

// type GhRequest struct {
// 	client_id     string
// 	client_secret string
// 	code          string
// }

func homePage(c *gin.Context) {
	c.HTML(200, "index.html", nil)
}

func oauth(c *gin.Context) {
	params := "?client_id=ff11d742cfea26f7953b&redirect_uri=http://localhost:8000/callback&"
	c.Redirect(http.StatusTemporaryRedirect, "https://github.com/login/oauth/authorize"+params)
}
func callback(c *gin.Context) {
	log.Println("Got here")
	code := c.Param("code")
	// reqBody := GhRequest{
	// 	client_id:     "ff11d742cfea26f7953b",
	// 	client_secret: "3d9eaed021bc91bd6427b6760a395e0319690c97",
	// 	code:          code,
	// }
	// post request to get access token
	reqBody := []byte(fmt.Sprintf(`{"client_id":"ff11d742cfea26f7953b","client_secret":"3d9eaed021bc91bd6427b6760a395e0319690c97","code":%q}`, code))
	req, err := http.NewRequest("POST", "https://github.com/login/oauth/access_token", bytes.NewBuffer(reqBody))
	req.Header.Set("Accept", "application:json")
	if err != nil {
		log.Println("request creation failed")
		return
	}
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("error in doing the request", err)
	} else {
		log.Println(resp)
	}
	var b []byte
	_, err = resp.Body.Read(b)
	if err != nil {
		log.Println(err)
	}
	jresponse := GhResponse{}
	log.Println(string(b))
	err = json.Unmarshal(b, &jresponse)
	if err != nil {
		log.Println("Marshalling error")
		return
	}
	log.Println("storing the token to the database", jresponse.Access_token)

}

func success(c *gin.Context) {
	log.Println("Oauth successful")
}

func main() {
	r := gin.Default()
	r.LoadHTMLGlob("templates/*")

	r.GET("/", homePage)
	r.GET("/oauth", oauth)
	r.GET("/callback", callback)
	r.GET("/success", success)

	err := r.Run(":8000")
	if err != nil {
		log.Fatal(err)
	}
}
