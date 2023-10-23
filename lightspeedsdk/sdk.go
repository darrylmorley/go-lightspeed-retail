package lightspeedsdk

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"
)

type SDK struct {
	BaseURL      string
	AccessToken  string
	RefreshToken string
	ClientID     string
	ClientSecret string
	Client       *http.Client
}

type RefreshResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	TokenType   string `json:"token_type"`
	Scope       string `json:"scope"`
}

func NewSDK(baseURL, refreshToken, clientID, clientSecret string) *SDK {
	return &SDK{
		BaseURL:      baseURL,
		AccessToken:  "",
		RefreshToken: refreshToken,
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Client:       &http.Client{},
	}
}

func (sdk *SDK) RefreshAccessToken(refreshToken, clientID, clientSecret string) (*RefreshResponse, error) {
	data := map[string]string{
		"refresh_token": refreshToken,
		"client_id":     clientID,
		"client_secret": clientSecret,
		"grant_type":    "refresh_token",
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", "https://cloud.lightspeedapp.com/oauth/access_token.php", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	resp, err := sdk.doWithRetry(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusUnauthorized {
			return nil, &UnauthorizedError{Message: "Token has expired or is invalid"}
		}
		return nil, parseAPIError(resp)
	}

	var refreshResp RefreshResponse
	err = json.NewDecoder(resp.Body).Decode(&refreshResp)
	if err != nil {
		return nil, err
	}

	sdk.AccessToken = refreshResp.AccessToken
	return &refreshResp, nil
}

func (sdk *SDK) DoGet(endpoint string, result interface{}) error {
	var allData []interface{}

	if !strings.HasSuffix(endpoint, ".json") {
		endpoint += ".json"
	}

	for endpoint != "" {
		fmt.Println("Endpoint:", endpoint)

		var url string
		if strings.HasPrefix(endpoint, "http") {
			url = endpoint
		} else {
			url = sdk.BaseURL + endpoint
		}

		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return err
		}

		req.Header.Set("Authorization", "Bearer "+sdk.AccessToken)

		resp, err := sdk.doWithRetry(req)
		if err != nil {
			return err
		}

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		fmt.Println("Raw Response:", string(body))

		// Reset the response body so it can be read again
		resp.Body = ioutil.NopCloser(bytes.NewBuffer(body))

		if resp.StatusCode != http.StatusOK {
			resp.Body.Close()
			return parseAPIError(resp)
		}

		defer resp.Body.Close() // Close the body here

		var tempResult map[string]interface{}
		err = json.Unmarshal(body, &tempResult)
		if err != nil {
			return err
		}

		for _, value := range tempResult {
			switch v := value.(type) {
			case []interface{}:
				for _, singleItem := range v {
					dataAsBytes, err := json.Marshal(singleItem)
					if err != nil {
						return err
					}
					var singleResult interface{}
					resultVal := reflect.ValueOf(result)
					if resultVal.Kind() == reflect.Ptr && resultVal.Elem().Kind() == reflect.Slice {
						singleResultType := resultVal.Elem().Type().Elem()
						singleResultPtr := reflect.New(singleResultType)
						singleResult = singleResultPtr.Interface()
						err = json.Unmarshal(dataAsBytes, singleResult)
						if err != nil {
							return err
						}
					} else {
						return errors.New("expected a pointer to a slice as the result parameter")
					}
					allData = append(allData, singleResult)
				}
			default:
				// handle or ignore unexpected types
			}
		}

		endpoint = getNextPageURL(resp)
	}

	// Transfer aggregated data to the result variable
	resultVal := reflect.ValueOf(result)
	if resultVal.Kind() == reflect.Ptr && resultVal.Elem().Kind() == reflect.Slice {
		resultSlice := resultVal.Elem()
		for _, v := range allData {
			resultSlice.Set(reflect.Append(resultSlice, reflect.Indirect(reflect.ValueOf(v))))
		}
	} else {
		return errors.New("expected a pointer to a slice as the result parameter")
	}

	return nil
}

func getNextPageURL(resp *http.Response) string {
	// Parse the response body as JSON
	var jsonResponse map[string]interface{}
	err := json.NewDecoder(resp.Body).Decode(&jsonResponse)
	if err != nil {
		return ""
	}

	// Extract the "next" attribute from the JSON response
	nextURL, _ := jsonResponse["@attributes"].(map[string]interface{})["next"].(string)
	fmt.Println(nextURL)
	return nextURL
}

// func (sdk *SDK) DoPost(endpoint string, data interface{}, result interface{}) error {
// 	// Preparation
// 	jsonData, err := json.Marshal(data)
// 	if err != nil {
// 		return err
// 	}

// 	// Request Creation
// 	url := sdk.BaseURL + endpoint
// 	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
// 	if err != nil {
// 		return err
// 	}

// 	// Headers
// 	req.Header.Set("Content-Type", "application/json")
// 	req.Header.Set("Authorization", "Bearer "+sdk.AccessToken)

// 	// Sending the Request
// 	resp, err := sdk.doWithRetry(req)
// 	if err != nil {
// 		return err
// 	}
// 	defer resp.Body.Close()

// 	// Error Handling
// 	if resp.StatusCode != http.StatusOK {
// 		if resp.StatusCode == http.StatusUnauthorized {
// 			return &UnauthorizedError{Message: "Token has expired or is invalid"}
// 		}
// 		return parseAPIError(resp)
// 	}

// 	// Response Parsing
// 	if result != nil {
// 		err = json.NewDecoder(resp.Body).Decode(result)
// 		if err != nil {
// 			return err
// 		}
// 	}

// 	return nil
// }
