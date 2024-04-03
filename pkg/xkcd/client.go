package xkcd

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

var NotFound = errors.New("client: comic not found")
var UnexpectedStatus = errors.New("client: unexpected status")

type Comic struct {
	Num        int    `json:"num"`
	Title      string `json:"title"`
	Transcript string `json:"transcript"`
	Alt        string `json:"alt"`
	Img        string `json:"img"`
}

type HttpClient struct {
	c   http.Client
	url string
}

type IClient interface {
	GetById(id int) (Comic, error)
}

func NewHttpClient(url string, timeout time.Duration) IClient {
	c := http.Client{Timeout: timeout}
	return &HttpClient{c, url}
}

func (xc *HttpClient) GetById(id int) (Comic, error) {
	resp, err := xc.doGet(xc.makeComicUrl(id))
	if err != nil {
		return Comic{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusNotFound {
			return Comic{}, NotFound
		}

		return Comic{}, UnexpectedStatus
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return Comic{}, err
	}

	comic, err := parseBody(body)
	if err != nil {
		return Comic{}, err
	}

	return comic, nil
}

func (xc *HttpClient) doGet(url string) (*http.Response, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")

	resp, err := xc.c.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (xc *HttpClient) makeComicUrl(id int) string {
	return fmt.Sprintf("%s/%d/info.0.json", xc.url, id)
}

func parseBody(b []byte) (Comic, error) {
	comic := &Comic{}
	return *comic, json.Unmarshal(b, comic)
}
