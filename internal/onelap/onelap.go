package onelap

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/go-resty/resty/v2"
)

const (
	OnelapSecret    = "fe9f8382418fcdeb136461cac6acae7b"
	LoginBaseURL    = "https://www.onelap.cn/api"
	AnalysisBaseURL = "https://u.onelap.cn/analysis"
)

type Client struct {
	restyClient *resty.Client
	UID         string
	XSRFToken   string
	OToken      string
}

func NewClient() *Client {
	client := resty.New().
		SetTimeout(30 * time.Second).
		SetRetryCount(3).
		SetRetryWaitTime(2 * time.Second).
		SetRetryMaxWaitTime(10 * time.Second)

	return &Client{
		restyClient: client,
	}
}

func md5Hex(s string) string {
	h := md5.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}

func randomHex(n int) string {
	const hexChars = "0123456789abcdef"
	b := make([]byte, n)
	for i := range b {
		b[i] = hexChars[rand.Intn(len(hexChars))]
	}
	return string(b)
}

func (c *Client) Login(account, password string) error {
	if account == "" || password == "" {
		return fmt.Errorf("onelap account and password cannot be empty")
	}
	timestamp := fmt.Sprintf("%d", time.Now().Unix())
	nonce := randomHex(16)
	passwordMd5 := md5Hex(password)

	// Signature calculation matching Onelap's verification
	signStr := fmt.Sprintf("account=%s&nonce=%s&password=%s&timestamp=%s&key=%s", account, nonce, passwordMd5, timestamp, OnelapSecret)
	sign := md5Hex(signStr)

	body := fmt.Sprintf(`{"account":"%s","password":"%s"}`, account, passwordMd5)

	resp, err := c.restyClient.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("nonce", nonce).
		SetHeader("timestamp", timestamp).
		SetHeader("sign", sign).
		SetBody(body).
		Post(LoginBaseURL + "/login")

	if err != nil {
		return fmt.Errorf("login request failed: %w", err)
	}

	if resp.StatusCode() != http.StatusOK {
		return fmt.Errorf("login failed with status: %s, body: %s", resp.Status(), resp.String())
	}

	// Extract cookies and tokens from response
	// The response structure usually contains userinfo and token
	// data[0].userinfo.uid, data[0].token, data[0].refresh_token
	type LoginResponse struct {
		Data []struct {
			Token        string `json:"token"`
			RefreshToken string `json:"refresh_token"`
			UserInfo     struct {
				UID json.Number `json:"uid"` // Use json.Number to avoid scientific notation
			} `json:"userinfo"`
		} `json:"data"`
	}

	var loginData LoginResponse
	if err := json.Unmarshal(resp.Body(), &loginData); err != nil {
		return fmt.Errorf("failed to unmarshal login response: %w", err)
	}

	if len(loginData.Data) == 0 {
		return fmt.Errorf("invalid login response: no data")
	}

	c.UID = loginData.Data[0].UserInfo.UID.String()
	c.XSRFToken = loginData.Data[0].Token
	c.OToken = loginData.Data[0].RefreshToken

	return nil
}

func (c *Client) Check(account, password string) error {
	return c.Login(account, password)
}

type Activity struct {
	ExternalID string      `json:"_id"` // Unique activity ID from Onelap
	UserID     json.Number `json:"id"`  // User ID
	FileKey    string      `json:"fileKey"`
	StartTime  string      `json:"date"` // Use 'date' field from API
	DURL       string      `json:"durl"`
}

func (c *Client) GetActivities() ([]Activity, error) {
	resp, err := c.restyClient.R().
		SetCookie(&http.Cookie{Name: "ouid", Value: c.UID}).
		SetCookie(&http.Cookie{Name: "XSRF-TOKEN", Value: c.XSRFToken}).
		SetCookie(&http.Cookie{Name: "OTOKEN", Value: c.OToken}).
		Get(AnalysisBaseURL + "/list")

	if err != nil {
		return nil, fmt.Errorf("get activity list failed: %w", err)
	}

	type ListResponse struct {
		Data []Activity `json:"data"`
	}

	var dataResponse ListResponse
	if err := json.Unmarshal(resp.Body(), &dataResponse); err != nil {
		return nil, fmt.Errorf("failed to unmarshal activity list: %w", err)
	}

	return dataResponse.Data, nil
}

func (c *Client) GetTodayActivities() ([]Activity, error) {
	all, err := c.GetActivities()
	if err != nil {
		return nil, err
	}

	// We'll check for activities in the last 24 hours to be more robust 
	// against timezone differences, or just use the date field.
	today := time.Now().Format("2006-01-02")
	yesterday := time.Now().Add(-24 * time.Hour).Format("2006-01-02")
	
	var todayActivities []Activity
	for _, act := range all {
		// Onelap 'date' format is usually "YYYY-MM-DD HH:MM"
		if len(act.StartTime) >= 10 {
			dateStr := act.StartTime[:10]
			if dateStr == today || dateStr == yesterday { // Include yesterday to be safe with sync time
				todayActivities = append(todayActivities, act)
			}
		}
	}

	return todayActivities, nil
}

func (c *Client) DownloadFIT(durl, destPath string) error {
	resp, err := c.restyClient.R().
		SetOutput(destPath).
		Get(durl)

	if err != nil {
		return fmt.Errorf("failed to download FIT file: %w", err)
	}

	if resp.StatusCode() != http.StatusOK {
		return fmt.Errorf("download failed with status: %s", resp.Status())
	}

	return nil
}
