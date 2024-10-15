package cmd

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

const (
	deviceCodeURL = "https://github.com/login/device/code"
	tokenURL      = "https://github.com/login/oauth/access_token"
)

var (
	clientID     string
	clientSecret string
)

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Log in to your GitHub account using Device Flow",
	Run: func(cmd *cobra.Command, args []string) {
		startDeviceFlow()
	},
}

type DeviceCodeResponse struct {
	DeviceCode      string `json:"device_code"`
	UserCode        string `json:"user_code"`
	VerificationURI string `json:"verification_uri"`
	ExpiresIn       int    `json:"expires_in"`
	Interval        int    `json:"interval"`
}

type TokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	Scope       string `json:"scope"`
}

func init() {
	err := godotenv.Load(GetConfigPath(".env"))
	if err != nil {
		logrus.Fatal("Error loading .env file")
	}

	logrus.Println("Loaded .env")
	logrus.Println(os.Getenv("CLIENT_ID"))
}

func startDeviceFlow() {
	clientID = os.Getenv("CLIENT_ID")
	clientSecret = os.Getenv("CLIENT_SECRET")

	if len(clientID) == 0 || len(clientSecret) == 0 {
		logrus.Fatalln("Invalid clientID or client secret")
	}

	codeResponse, err := requestDeviceCode()
	if err != nil {
		logrus.WithError(err).Warnln("Error requesting device code")
		return
	}

	fmt.Printf("Please visit %s and enter the code: %s\n", codeResponse.VerificationURI, codeResponse.UserCode)

	tokenResponse, err := pollForAccessToken(codeResponse.DeviceCode, codeResponse.Interval)
	if err != nil {
		logrus.WithError(err).Warnln("Error during token polling")
		return
	}

	saveToken(tokenResponse.AccessToken)
	fmt.Println("Login successful! Token saved.")
}

func requestDeviceCode() (*DeviceCodeResponse, error) {
	data := fmt.Sprintf("client_id=%s&scope=repo", clientID)
	resp, err := http.Post(deviceCodeURL, "application/x-www-form-urlencoded", strings.NewReader(data))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	// fmt.Println("Response body:", string(body))

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received non-200 response: %s", string(body))
	}

	values, err := url.ParseQuery(string(body))
	if err != nil {
		return nil, err
	}

	codeResponse := &DeviceCodeResponse{
		DeviceCode:      values.Get("device_code"),
		UserCode:        values.Get("user_code"),
		VerificationURI: values.Get("verification_uri"),
		ExpiresIn:       parseInt(values.Get("expires_in")),
		Interval:        parseInt(values.Get("interval")),
	}

	return codeResponse, nil
}

func parseInt(s string) int {
	if val, err := strconv.Atoi(s); err == nil {
		return val
	}
	return 0
}

func pollForAccessToken(deviceCode string, interval int) (*TokenResponse, error) {
	for {
		time.Sleep(time.Duration(interval) * time.Second)

		data := fmt.Sprintf("client_id=%s&client_secret=%s&device_code=%s&grant_type=urn:ietf:params:oauth:grant-type:device_code", clientID, clientSecret, deviceCode)
		resp, err := http.Post(tokenURL, "application/x-www-form-urlencoded", strings.NewReader(data))
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		if strings.Contains(string(body), "authorization_pending") {
			continue
		}

		values, err := url.ParseQuery(string(body))
		if err != nil {
			logrus.WithError(err).Warnln("Error parsing token response")
			return nil, fmt.Errorf("error parsing token response: %v", err)
		}

		if accessToken := values.Get("access_token"); accessToken != "" {
			return &TokenResponse{
				AccessToken: accessToken,
				Scope:       values.Get("scope"),
				TokenType:   values.Get("token_type"),
			}, nil
		}

		logrus.WithError(err).Warnln("Error parsing token response")
		return nil, fmt.Errorf("error parsing token response: %v", err)
	}
}

func saveToken(token string) {
	err := os.WriteFile(GetConfigPath(".token"), []byte(token), 0600)
	if err != nil {
		logrus.WithError(err).Warnln("Error saving token")
		return
	}
}
