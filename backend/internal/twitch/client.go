package twitch

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

const (
	twitchAPIBaseURL = "https://api.twitch.tv/helix"
	twitchAuthURL    = "https://id.twitch.tv/oauth2/token"
)

type Client struct {
	clientID     string
	clientSecret string
	httpClient   *http.Client
}

func NewClient(clientID, clientSecret string) (*Client, error) {
	return &Client{
		clientID:     clientID,
		clientSecret: clientSecret,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}, nil
}

func (c *Client) makeRequest(method, endpoint string, headers map[string]string, params url.Values) (*http.Response, error) {
	reqURL := twitchAPIBaseURL + endpoint
	if len(params) > 0 {
		reqURL += "?" + params.Encode()
	}

	req, err := http.NewRequest(method, reqURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Client-ID", c.clientID)
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	return c.httpClient.Do(req)
}

func (c *Client) GetChannelInfoWithToken(accessToken string) (*ChannelInfo, error) {
	userID, err := c.getUserID(accessToken)
	if err != nil {
		return nil, err
	}

	headers := map[string]string{
		"Authorization": "Bearer " + accessToken,
	}

	params := url.Values{}
	params.Set("broadcaster_id", userID)

	resp, err := c.makeRequest("GET", "/channels", headers, params)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("twitch API error: %d", resp.StatusCode)
	}

	var channelResp ChannelResponse
	if err := json.NewDecoder(resp.Body).Decode(&channelResp); err != nil {
		return nil, err
	}

	if len(channelResp.Data) == 0 {
		return nil, fmt.Errorf("no channel data returned")
	}

	return &channelResp.Data[0], nil
}

func (c *Client) GetFollowerCount(accessToken string) (int, error) {
	userID, err := c.getUserID(accessToken)
	if err != nil {
		return 0, err
	}

	headers := map[string]string{
		"Authorization": "Bearer " + accessToken,
	}

	params := url.Values{}
	params.Set("broadcaster_id", userID)
	params.Set("first", "1")

	resp, err := c.makeRequest("GET", "/channels/followers", headers, params)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("twitch API error: %d", resp.StatusCode)
	}

	var followersResp FollowersResponse
	if err := json.NewDecoder(resp.Body).Decode(&followersResp); err != nil {
		return 0, err
	}

	return followersResp.Total, nil
}

func (c *Client) GetSubscriberCount(accessToken string) (int, error) {
	userID, err := c.getUserID(accessToken)
	if err != nil {
		return 0, err
	}

	headers := map[string]string{
		"Authorization": "Bearer " + accessToken,
	}

	params := url.Values{}
	params.Set("broadcaster_id", userID)

	resp, err := c.makeRequest("GET", "/subscriptions", headers, params)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, nil
	}

	var subsResp SubscribersResponse
	if err := json.NewDecoder(resp.Body).Decode(&subsResp); err != nil {
		return 0, err
	}

	return len(subsResp.Data), nil
}

func (c *Client) GetVideos(accessToken, videoType string, limit int) ([]VideoInfo, error) {
	userID, err := c.getUserID(accessToken)
	if err != nil {
		return nil, err
	}

	headers := map[string]string{
		"Authorization": "Bearer " + accessToken,
	}

	params := url.Values{}
	params.Set("user_id", userID)
	params.Set("type", videoType)
	params.Set("first", fmt.Sprintf("%d", limit))

	resp, err := c.makeRequest("GET", "/videos", headers, params)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("twitch API error: %d", resp.StatusCode)
	}

	var videosResp VideosResponse
	if err := json.NewDecoder(resp.Body).Decode(&videosResp); err != nil {
		return nil, err
	}

	return videosResp.Data, nil
}

func (c *Client) GetStreamInfo(accessToken string) (*StreamInfo, error) {
	userID, err := c.getUserID(accessToken)
	if err != nil {
		return nil, err
	}

	headers := map[string]string{
		"Authorization": "Bearer " + accessToken,
	}

	params := url.Values{}
	params.Set("user_id", userID)

	resp, err := c.makeRequest("GET", "/streams", headers, params)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("twitch API error: %d", resp.StatusCode)
	}

	var streamResp StreamResponse
	if err := json.NewDecoder(resp.Body).Decode(&streamResp); err != nil {
		return nil, err
	}

	if len(streamResp.Data) == 0 {
		return nil, nil
	}

	return &streamResp.Data[0], nil
}

func (c *Client) GetUserInfo(accessToken string) (*User, error) {
	headers := map[string]string{
		"Authorization": "Bearer " + accessToken,
	}

	resp, err := c.makeRequest("GET", "/users", headers, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("twitch API error: %d", resp.StatusCode)
	}

	var userResp UsersResponse
	if err := json.NewDecoder(resp.Body).Decode(&userResp); err != nil {
		return nil, err
	}

	if len(userResp.Data) == 0 {
		return nil, fmt.Errorf("no user data returned")
	}

	return &userResp.Data[0], nil
}

func (c *Client) getUserID(accessToken string) (string, error) {
	user, err := c.GetUserInfo(accessToken)
	if err != nil {
		return "", err
	}
	return user.ID, nil
}
