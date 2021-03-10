package main

import (
	"backgroundl/reddit"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	version := os.Getenv("VERSION")
	fmt.Println(fmt.Sprintf("Starting Downloadl v%s", version))

	configDir, err := os.UserConfigDir()
	if err != nil {
		panic(err)
	}

	fmt.Println(configDir) // C:\Users\YourUser

	client := &http.Client{}

	authResp, err := authenticate(client)
	if err != nil {
		panic(err)
	}

	fmt.Println(authResp.AccessToken)

	resp, err := getListing(client, authResp)
	if err != nil {
		panic(err)
	}

	for _, post := range resp.Data.Children {
		fmt.Println(post.Data.URLOverriddenByDest)
	}
}

func getListing(c *http.Client, s *reddit.InstalledClientAuthentication) (*reddit.ListingResponse, error) {
	requestURL := "https://oauth.reddit.com/r/earthporn/top"
	userAgent := os.Getenv("USER_AGENT")

	r, _ := http.NewRequest(http.MethodGet, requestURL, nil) // URL-encoded payload
	r.Header.Add("Authorization", fmt.Sprintf("Bearer %s", s.AccessToken))
	r.Header.Add("User-Agent", userAgent)

	resp, err := c.Do(r)
	if err != nil {
		return nil, err
	}

	fmt.Println(resp.StatusCode)

	result := &reddit.ListingResponse{}

	err = json.NewDecoder(resp.Body).Decode(result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func basicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}

func authenticate(c *http.Client) (*reddit.InstalledClientAuthentication, error) {
	authURL := "https://www.reddit.com/api/v1/access_token"
	deviceID := os.Getenv("DEVICE_ID")
	grantType := os.Getenv("GRANT_TYPE")
	clientID := os.Getenv("CLIENT_ID")
	userAgent := os.Getenv("USER_AGENT")

	data := url.Values{}
	data.Set("grant_type", grantType)
	data.Set("device_id", deviceID)

	r, _ := http.NewRequest(http.MethodPost, authURL, strings.NewReader(data.Encode())) // URL-encoded payload
	r.Header.Add("Authorization", fmt.Sprintf("Basic %s", basicAuth(clientID, "")))
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))
	r.Header.Add("User-Agent", userAgent)

	resp, err := c.Do(r)
	if err != nil {
		return nil, err
	}
	fmt.Printf("Authentication: %s\n", resp.Status)

	authResp := &reddit.InstalledClientAuthentication{}

	err = json.NewDecoder(resp.Body).Decode(authResp)
	if err != nil {
		return nil, err
	}
	return authResp, nil
}
