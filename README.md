A Lightspeed Retail SDK in Go!

To initialise the SDK:

```
func main() {
    // Read values from .env or another source
    baseURL := os.Getenv("BASE_URL")
    accessToken := os.Getenv("ACCESS_TOKEN")
    refreshToken := os.Getenv("REFRESH_TOKEN")

    // Create SDK instance with the retrieved values
    sdk := lightspeedsdk.NewSDK(baseURL, accessToken, refreshToken)
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
