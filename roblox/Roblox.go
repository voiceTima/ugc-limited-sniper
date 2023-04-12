package roblox

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
	"time"
)

type Roblox struct {
	client    *http.Client
	cookie    string
	xsrfToken string
	userID    int
	username  string
	mu        sync.Mutex
}

type Config struct {
	Cookie string
}

func New(config Config) *Roblox {
	client := &http.Client{
		Timeout: time.Second * 10,
	}
	return &Roblox{client: client, cookie: config.Cookie}
}

func (r *Roblox) request(method, url string, headers map[string]string, body string) (*http.Response, error) {
	req, err := http.NewRequest(method, url, strings.NewReader(body))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Cookie", ".ROBLOSECURITY="+r.cookie)
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	return r.client.Do(req)
}

func (r *Roblox) RefreshXsrfToken() (string, error) {
	// This function does not actually invalidate your cookie. You would need to post the XSRF token in the headers, otherwise anyone could log you out of your Roblox account by making you click a link.
	// https://owasp.org/www-community/attacks/csrf
	// https://portswigger.net/web-security/csrf
	resp, err := r.request("POST", "https://auth.roblox.com/v2/logout", nil, "")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	token := resp.Header.Get("x-csrf-token")
	if token == "" {
		return "", errors.New("no token found in response headers")
	}

	return token, nil
}

func (r *Roblox) SetXsrfToken(token string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.xsrfToken = token
}

func (r *Roblox) GetXsrfToken() string {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.xsrfToken
}

func (r *Roblox) GetUserInfo() (*CurrentUserResponse, error) {
	resp, err := r.request("GET", "https://users.roblox.com/v1/users/authenticated", nil, "")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var currentUser CurrentUserResponse
	err = json.Unmarshal(body, &currentUser)
	if err != nil {
		return nil, err
	}

	return &currentUser, nil
}

func (r *Roblox) SetCurrentUser(userID int, username string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.userID = userID
	r.username = username
}

func (r *Roblox) GetCurrentUser() (int, string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.userID, r.username
}

func (r *Roblox) GetAssetsInfo(assetIds []int) ([]POSTItemDetails, error) {
	items := make([]map[string]interface{}, len(assetIds))
	for i, id := range assetIds {
		items[i] = map[string]interface{}{
			"itemType": "Asset",
			"id":       id,
		}
	}

	jsonBody, err := json.Marshal(map[string]interface{}{
		"items": items,
	})
	if err != nil {
		return nil, err
	}

	resp, err := r.request("POST", "https://catalog.roblox.com/v1/catalog/items/details", map[string]string{
		"x-csrf-token": r.xsrfToken,
		"Content-Type": "application/json",
	}, string(jsonBody))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var catalogDetailsResponse CatalogDetailsResponse
	err = json.Unmarshal(body, &catalogDetailsResponse)
	if err != nil {
		return nil, err
	}

	return catalogDetailsResponse.Data, nil
}

func (r *Roblox) BuyItem(itemId int, productId string, sellerId int) (*PurchaseItemResponse, error) {
	// This function will only work on a free item, but you can edit the code to make it support paid items as well.
	// If you decide to do this, I recommend limiting the price so you don't accidentally buy an item worth some insane amount of Robux.
	headers := map[string]string{
		"x-csrf-token": r.xsrfToken,
		"Content-Type": "application/json",
	}

	jsonBody, err := json.Marshal(map[string]interface{}{
		"collectibleItemId":     itemId,
		"expectedCurrency":      1,
		"expectedPrice":         0,
		"expectedPurchaserId":   r.userID,
		"expectedPurchaserType": "User",
		"expectedSellerId":      sellerId,
		"expectedSellerType":    "User",
		"idempotencyKey":        uuid.New().String(),
		"collectibleProductId":  productId,
	})
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("https://apis.roblox.com/marketplace-sales/v1/item/%d/purchase-item", itemId)
	resp, err := r.request("POST", url, headers, string(jsonBody))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var purchaseItemResponse PurchaseItemResponse
	err = json.Unmarshal(body, &purchaseItemResponse)
	if err != nil {
		return nil, err
	}

	return &purchaseItemResponse, nil
}
