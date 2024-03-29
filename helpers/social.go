// THIS FILE WAS HATCHED WITH github.com/axiomzen/hatch
// THIS FILE IS SAFE TO EDIT BY HAND BUT YOU PROBABLY DON'T NEED TO

package helpers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

const (
	facebookTokenURL       = "https://graph.facebook.com/v2.10/debug_token?" //fields=id&access_token=@accesstoken
	facebookUserURL        = "https://graph.facebook.com/v2.10/"
	facebookUserPictureURL = "https://graph.facebook.com/%s/picture?type=large"
)

// FacebookAPIUser represents the parts of the user we're interested about
type FacebookAPIUser struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	Picture FBData `json:"picture"`
}

type FBData struct {
	Data map[string]interface{} `json:"data"`
}

func (fbu *FacebookAPIUser) ProfilePictureURL() string {
	return fbu.Picture.Data["url"].(string)
}

// ValidateFacebookLogin takes the id and token strings and sends them to the FACEBOOK_TOKEN_URL.
// If the inputs are valid, returns true, else it returns false and an error
func ValidateFacebookLogin(id, token, appID, appSecret string) (bool, error) {

	client := http.Client{}
	urlValues := url.Values{}
	urlValues.Set("input_token", token)
	urlValues.Set("access_token", appID+"|"+appSecret)
	req, _ := http.NewRequest("GET", facebookTokenURL+urlValues.Encode(), nil)
	req.Close = true
	// Accept type?
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	defer func(resp *http.Response) {
		if resp != nil && resp.Body != nil {
			resp.Body.Close()
		}
	}(resp)
	if err != nil {
		return false, err
	}

	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("Could not validate facebook token, response code: %d", resp.StatusCode)
	}
	// TODO: use a struct for this
	var respJSON map[string]interface{}
	contentType := req.Header.Get("Content-Type")
	if strings.Contains(contentType, "application/json") {
		decoder := json.NewDecoder(resp.Body)
		parseErr := decoder.Decode(&respJSON)
		if parseErr != nil {
			return false, parseErr
		}
		// Check ids match
		data := respJSON["data"].(map[string]interface{})
		if facebookError, ok := data["error"]; ok {
			return false, fmt.Errorf("Facebook error: %v", facebookError)
		} else if data["user_id"] == nil {
			return false, errors.New("Facebook User ID was nil")
		} else if data["user_id"].(string) != id {
			return false, fmt.Errorf("User ID %s did not match Facebook token: %s", id, data["user_id"].(string))
		}
	} else {
		return false, errors.New("unexpected content type")
	}

	return true, nil
}

// GetFacebookUserInfo takes the id and token strings and sends them to the facebookUserURL.
// Returns a FacebookAPIUser struct
func GetFacebookUserInfo(id, token, appID, appSecret string) (*FacebookAPIUser, error) {

	client := http.Client{}
	urlValues := url.Values{}
	urlValues.Set("input_token", token)
	urlValues.Set("access_token", appID+"|"+appSecret)
	urlValues.Set("fields", "name,first_name,last_name,email,picture.type(large){url}")
	req, _ := http.NewRequest("GET", facebookUserURL+id+"?"+urlValues.Encode(), nil)
	req.Close = true
	// Accept type?
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	defer func(resp *http.Response) {
		if resp != nil && resp.Body != nil {
			resp.Body.Close()
		}
	}(resp)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Could not get facebook user, response code: %d", resp.StatusCode)
	}
	apiUser := FacebookAPIUser{}
	contentType := req.Header.Get("Content-Type")
	if strings.Contains(contentType, "application/json") {
		decoder := json.NewDecoder(resp.Body)
		parseErr := decoder.Decode(&apiUser)
		if parseErr != nil {
			return nil, parseErr
		}
		return &apiUser, nil
	}
	return nil, errors.New("unexpected content type")
}
