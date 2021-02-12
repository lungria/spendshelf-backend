package mono

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// GetUserInfo loads user accounts list.
func (c *Client) GetUserInfo(ctx context.Context) ([]Account, error) {
	uri := fmt.Sprintf("%s/personal/client-info", c.baseURL)

	response, err := c.performRequest(ctx, uri, http.MethodGet, nil)
	if err != nil {
		return nil, fmt.Errorf("unable to get accounts: %w", err)
	}

	if len(response) == 0 {
		return nil, fmt.Errorf("unable to get accounts: empty response without error received")
	}

	parsedResponse := userInfoResponse{}
	if err = json.Unmarshal(response, &parsedResponse); err != nil {
		return nil, fmt.Errorf("failed unmarshal accounts form json body: %w", err)
	}

	return parsedResponse.Accounts, nil
}

type userInfoResponse struct {
	Accounts []Account `json:"accounts"`
}

// Account describes monobank account.
type Account struct {
	ID           string `json:"id"`
	Balance      int64  `json:"balance"`
	CurrencyCode int16  `json:"currencyCode"`
}
