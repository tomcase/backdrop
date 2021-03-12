package main

import (
	"backdropGo/reddit"
	"backdropGo/sqlite"
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

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	version := os.Getenv("VERSION")
	log.Println(fmt.Sprintf("Starting Downloadl v%s", version))

	err = sqlite.CreateDatabase()
	if err != nil {
		log.Fatal(err)
	}

	client := &http.Client{}

	authResp, err := authenticate(client)
	if err != nil {
		log.Fatal(err)
	}

	resp, err := getListing(client, authResp)
	if err != nil {
		log.Fatal(err)
	}

	for _, post := range resp.Data.Children {
		outputDir := os.Getenv("OUTPUT_DIR")
		outputFile := filepath.Join(outputDir, filepath.Base(post.Data.URL))
		err = downloadFile(post.Data.URL, outputFile)
		if err != nil {
			log.Fatal(err)
		}

		err = checkImageWidth(outputFile)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func checkImageWidth(fileName string) error {
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
		return errors.New("Received non 200 response code")
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
	userAgent := os.Getenv("USER_AGENT")

	r, _ := http.NewRequest(http.MethodGet, requestURL, nil) // URL-encoded payload
	r.Header.Add("Authorization", fmt.Sprintf("Bearer %s", s.AccessToken))
	r.Header.Add("User-Agent", userAgent)

	resp, err := c.Do(r)
	if err != nil {
		return nil, err
	}

	log.Println(resp.StatusCode)

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
