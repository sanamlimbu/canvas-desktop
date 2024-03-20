package canvas

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/ninja-software/terror/v2"
)

type User struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	SISUserID string `json:"sis_user_id"`
}

func (c *APIClient) GetUserBySisID(sisID string) (*User, error) {
	user := &User{}

	requestURL := fmt.Sprintf("%s/users/sis_user_id:%s", c.BaseURL, sisID)
	req, err := http.NewRequest(http.MethodGet, requestURL, nil)
	if err != nil {
		return nil, terror.Error(err, "cannot create a get request")
	}
	bearer := "Bearer " + c.AccessToken
	req.Header.Add("Authorization", bearer)

	res, err := c.do(req)
	if err != nil {
		return nil, terror.Error(err, "error on get request call")
	}
	defer res.Body.Close()

	if res.Status != "200 OK" {
		return nil, terror.Error(fmt.Errorf("user not found"), "something went wrong and did not receive 200 OK status")
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, terror.Error(err, "cannot read response body")
	}

	if err := json.Unmarshal(body, user); err != nil {
		return nil, terror.Error(err, "cannot unmarshal response body")
	}
	return user, nil

}
