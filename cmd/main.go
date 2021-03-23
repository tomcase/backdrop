package main

import (
	"backdropGo/reddit"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/joho/godotenv"
)

var (
	CLIENT_ID string = "VIz3jyOacujEuQ"
)

func main() {
	defaultContext := context.Background()

	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}

	version := os.Getenv("BG_VERSION")
	log.Println(fmt.Sprintf("Starting Backdrop v%s", version))

	checkDeviceIDExists(defaultContext)

	client := &http.Client{}

	authResp, err := authenticate(defaultContext, client)
	if err != nil {
		log.Fatal(err)
	}

	resp, err := getListing(client, authResp)
	if err != nil {
		log.Fatal(err)
	}

	for _, post := range resp.Data.Children {
		outputDir := "/output"
		outputFile := filepath.Join(outputDir, filepath.Base(post.Data.URL))
		err = downloadFile(post.Data.URL, outputFile)
		if err != nil {
			log.Fatal(err)
		}

		err = verifyImage(outputFile)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func checkDeviceIDExists(ctx context.Context) error {
	deviceID := os.Getenv("BG_DEVICE_ID")
	if deviceID == "" {
		return errors.New("You must set the DEVICE_ID environment variable to a unique string of 20-30 characters.")
	}
	return nil
}

func verifyImage(fileName string) error {
	ext := filepath.Ext(fileName)
	if ext != ".jpg" && ext != ".png" {
		os.Remove(fileName)
		return nil
	}

	file, err := os.Open(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	im, _, err := image.DecodeConfig(file)
	if err != nil {
		return err
	}

	if im.Width < 2048 || im.Height < 1280 || im.Width < im.Height {
		os.Remove(fileName)
	}

	return nil
}

func downloadFile(URL, fileName string) error {
	//Get the response bytes from the url
	response, err := http.Get(URL)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		msg := fmt.Sprintf("Failed to download from %s - StatusCode: %d -- SKIPPING", URL, response.StatusCode)
		log.Default().Println(msg)
		return nil
	}

	//Create a empty file
	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	//Write the bytes to the fiel
	written, err := io.Copy(file, response.Body)
	if err != nil {
		return err
	}

	log.Println(fmt.Sprintf("Finished writing %s with %d bytes written.", fileName, written))

	return nil
}

func getListing(c *http.Client, s *reddit.InstalledClientAuthentication) (*reddit.ListingResponse, error) {
	requestURL := "https://oauth.reddit.com/r/earthporn/top?limit=100"
	userAgent := os.Getenv("BG_USER_AGENT")

	r, _ := http.NewRequest(http.MethodGet, requestURL, nil) // URL-encoded payload
	r.Header.Add("Authorization", fmt.Sprintf("Bearer %s", s.AccessToken))
	r.Header.Add("User-Agent", userAgent)

	resp, err := c.Do(r)
	if err != nil {
		return nil, err
	}

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

func authenticate(ctx context.Context, c *http.Client) (*reddit.InstalledClientAuthentication, error) {
	authURL := "https://www.reddit.com/api/v1/access_token"
	deviceID := os.Getenv("BG_DEVICE_ID")
	grantType := "https://oauth.reddit.com/grants/installed_client"
	clientID := CLIENT_ID
	userAgent := os.Getenv("BG_USER_AGENT")

	data := url.Values{}
	data.Set("grant_type", grantType)
	data.Set("device_id", deviceID)

	r, _ := http.NewRequest(http.MethodPost, authURL, strings.NewReader(data.Encode())) // URL-encoded payload
	r.Header.Add("Authorization", fmt.Sprintf("Basic %s", basicAuth(clientID, "")))
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("User-Agent", userAgent)

	resp, err := c.Do(r)
	if err != nil {
		return nil, err
	}

	authResp := &reddit.InstalledClientAuthentication{}

	err = json.NewDecoder(resp.Body).Decode(authResp)
	if err != nil {
		return nil, err
	}
	return authResp, nil
}
