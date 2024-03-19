package canvas

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/ninja-software/terror/v2"
)

type Account struct {
	ID              int    `json:"id"`
	Name            string `json:"name"`
	ParentAccountID int    `json:"parent_account_id"`
	RootAccountID   int    `json:"root_account_id"`
}

func (c *APIClient) GetAccountByID(accountID int) (*Account, error) {
	account := &Account{}

	requestURL := fmt.Sprintf("%s/accounts/%d", c.BaseURL, accountID)
	req, err := http.NewRequest(http.MethodGet, requestURL, nil)
	if err != nil {
		fmt.Println(err)
		return nil, terror.Error(err, "cannot create a get request")
	}
	bearer := "Bearer " + c.AccessToken
	req.Header.Add("Authorization", bearer)

	res, err := c.do(req)
	if err != nil {
		fmt.Println(err)
		return nil, terror.Error(err, "error on get request call")
	}
	defer res.Body.Close()

	if res.Status != "200 OK" {
		fmt.Println(err)
		return nil, terror.Error(fmt.Errorf("status code: %d", res.StatusCode), "something went wrong and did not receive 200 OK status")
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return nil, terror.Error(err, "cannot read response body")
	}

	if err := json.Unmarshal(body, account); err != nil {
		fmt.Println(err)
		return nil, terror.Error(err, "cannot unmarshal response body")
	}
	return account, nil
}
