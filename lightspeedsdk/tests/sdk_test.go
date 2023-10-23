package lightspeedsdk_test

import (
	"fmt"
	"os"
	"testing"

	lightspeedsdk "github.com/darrylmorley/go-lightspeed-retail/lightspeedsdk"
	"github.com/joho/godotenv"
)

type Category struct {
	CategoryID   string `json:"categoryID"`
	Name         string `json:"name"`
	NodeDepth    string `json:"nodeDepth"`
	FullPathName string `json:"fullPathName"`
	LeftNode     string `json:"leftNode"`
	RightNode    string `json:"rightNode"`
	CreateTime   string `json:"createTime"`
	TimeStamp    string `json:"timeStamp"`
	ParentID     string `json:"parentID"`
}

type Response struct {
	Attributes struct {
		Next     string `json:"next"`
		Previous string `json:"previous"`
	} `json:"@attributes"`
	Category []Category `json:"Category"`
}

func TestMain(m *testing.M) {
	// Load environment variables from .env
	if err := godotenv.Load("../../.env"); err != nil {
		// Handle error if .env file is not found or cannot be loaded
		panic("Error loading .env file")
	}

	// Run the tests
	result := m.Run()

	// Clean up or perform other actions after tests if needed

	// Exit with the test result
	os.Exit(result)
}

func TestRefreshAccessToken(t *testing.T) {
	baseURL := os.Getenv("BASE_URL")
	clientID := os.Getenv("CLIENT_ID")
	clientSecret := os.Getenv("CLIENT_SECRET")
	refreshToken := os.Getenv("CLIENT_REFRESH_TOKEN")
	// Create a new SDK instance.
	sdk := lightspeedsdk.NewSDK(baseURL, clientID, clientSecret, refreshToken)

	// Use the SDK's RefreshAccessToken method.
	resp, err := sdk.RefreshAccessToken(refreshToken, clientID, clientSecret)
	if err != nil {
		t.Fatalf("Error refreshing access token: %s\n", err)
		return
	}

	fmt.Printf("New Access Token: %s\n", resp.AccessToken)
}

func TestDoGetCategories(t *testing.T) {
	baseURL := os.Getenv("BASE_URL")
	clientID := os.Getenv("CLIENT_ID")
	clientSecret := os.Getenv("CLIENT_SECRET")
	refreshToken := os.Getenv("CLIENT_REFRESH_TOKEN")
	// Create a new SDK instance with the mock server's URL as the BaseURL
	sdk := lightspeedsdk.NewSDK(baseURL, clientID, clientSecret, refreshToken)

	resp, err := sdk.RefreshAccessToken(refreshToken, clientID, clientSecret)
	if err != nil {
		t.Fatalf("Error refreshing access token: %s\n", err)
		return
	}
	fmt.Printf("New Access Token: %s\n", resp.AccessToken)

	// Define a struct to hold the response data
	var response Response

	// Call the doGet function to fetch categories
	err = sdk.DoGet("/Category", &response.Category)
	if err != nil {
		t.Fatalf("Error fetching categories: %v", err)
	}
	fmt.Printf("Categories: %+v\n", response)
	fmt.Println("Number of categories:", len(response.Category))
}
