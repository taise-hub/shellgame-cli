package shellgame

import (
	"bytes"
	"io/ioutil"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/taise-hub/shellgame-cli/server/domain/model"
	"net/http"
	"net/url"
	"time"
)

const (
	HOST = "localhost"
)

var (
	baseEndpoint    = &url.URL{Scheme: "http", Host: HOST, Path: "/"}
	profileEndpoint = &url.URL{Scheme: "http", Host: HOST, Path: "/profile"}
	shellEndpoint   = &url.URL{Scheme: "ws", Host: HOST, Path: "/shell"}
)

// Clientはシェルゲーサーバとネットワーク的なやりとりを行う役割を持つ
type Client struct {
	*http.Client
	baseEndpoint    *url.URL
	profileEndpoint *url.URL
	shellEndpoint 	*url.URL
}

func newClient() (*Client, error) {
	c := &Client{}
	jar, err := getJar()
	if err != nil {
		return nil, err
	}
	c.Client = &http.Client{
		Jar:     jar,
		Timeout: 20 * time.Second,
	}
	c.baseEndpoint =  &url.URL{Scheme: "http", Host: HOST, Path: "/"}
	c.profileEndpoint = &url.URL{Scheme: "http", Host: HOST, Path: "/profile"}
	c.shellEndpoint   = &url.URL{Scheme: "ws", Host: HOST, Path: "/shell"}
	return c, nil
}

// シェルゲーサーバで稼働するコンテナにWebSocketを利用して接続する。
func ConnectShell() (*websocket.Conn, error) {
	var header http.Header
	jar, err := getJar()
	if err != nil {
		return nil, err
	}
	for _, cookie := range jar.Cookies(baseEndpoint) {
		header.Add("Cookie", fmt.Sprintf("%s=%s", cookie.Name, cookie.Value))
	}

	wsconn, _, err := websocket.DefaultDialer.Dial(shellEndpoint.String(), header)
	if err != nil {
		return nil, err
	}

	return wsconn, nil
}

// シェルゲーサーバにプレイヤー名を登録する。
func (c *Client) PostProfile(name string) error {
	id := uuid.New()
	profile := &model.Profile{ID: id.String(), Name: name}
	p, err := json.Marshal(profile)
	if err != nil {
		return err
	}

	request, err := http.NewRequest("POST", c.profileEndpoint.String(), bytes.NewBuffer(p))
	if err != nil {
		return err
	}
	request.Header.Set("Content-Type", "application/json; charset=UTF-8")

	resp, err := c.Do(request)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body")
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("%s", body)
	}

	return nil
}
