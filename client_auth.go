package salt

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/pkg/errors"
)

// Define the struct to match the JSON response
type LoginApiResponse struct {
	Return []LoginTokenInfo `json:"return"`
}

type LoginTokenInfo struct {
	Token  string   `json:"token"`
	Expire float64  `json:"expire"`
	Start  float64  `json:"start"`
	User   string   `json:"user"`
	Eauth  string   `json:"eauth"`
	Perms  []string `json:"perms"`
}

func (c *Client) Login(ctx context.Context, username, password string) error {
	req := Request{
		"username": username,
		"password": password,
		"eauth":    "pam",
	}

	return c.do(ctx, "POST", "login", req, func(r *http.Response) error {
		c.Cookies = r.Cookies()

		c.Token = r.Header.Get("X-Auth-Token")

		body, err := io.ReadAll(r.Body)
		if err != nil {
			return errors.Wrap(err, "Error reading response body")
		}

		var response LoginApiResponse
		err = json.Unmarshal(body, &response)
		if err != nil {
			return errors.Wrap(err, "Error parsing Login JSON")
		}

		// Convert Expire and Start to time.Time
		c.TokenExpire = time.Unix(int64(response.Return[0].Expire), 0)

		return nil
	})
}

func (c *Client) Logout(ctx context.Context) error {
	c.TokenExpire = time.Time{}
	return c.do(ctx, "POST", "logout", nil, nil)
}
