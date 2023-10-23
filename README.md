A Lightspeed Retail SDK in Go!

To initialise the SDK:

```
func main() {
    // Read values from .env or another source
    baseURL := os.Getenv("BASE_URL")
	clientID := os.Getenv("CLIENT_ID")
	clientSecret := os.Getenv("CLIENT_SECRET")
	refreshToken := os.Getenv("CLIENT_REFRESH_TOKEN")

    // Create a new SDK instance with the mock server's URL as the BaseURL
	sdk := lightspeedsdk.NewSDK(baseURL, clientID, clientSecret, refreshToken)
    // ...
}
```

Example, call Categories:

```
	err = sdk.DoGet("/Category", &response)
	if err != nil {
		t.Fatalf("Error fetching categories: %v", err)
	}
	fmt.Printf("Categories: %+v\n", response)
	fmt.Println("Number of categories:", len(response))
```

Currently only DoGet is implemented.
