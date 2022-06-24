package eureka_client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

const appURI = "apps/"

// 调用 eureka 服务端 rest API
// https://github.com/Netflix/eureka/wiki/Eureka-REST-operations

type Client struct {
	url                          string
	registryFetchIntervalSeconds time.Duration
	ticker                       *time.Ticker
	tickerStop                   chan struct{}
	stop                         chan struct{}
	httpClient                   http.Client
	instance                     *Instance
}

type DialOption func(clt *Client)

func Dial(ins *Instance, opts ...DialOption) *Client {
	clt := &Client{
		url:                          "http://admin:admin@localhost:8761/eureka/",
		registryFetchIntervalSeconds: 15 * time.Second,
		tickerStop:                   make(chan struct{}),
		stop:                         make(chan struct{}),
		httpClient:                   *http.DefaultClient,
		instance:                     ins,
	}
	clt.ticker = time.NewTicker(clt.registryFetchIntervalSeconds)

	for _, opt := range opts {
		opt(clt)
	}

	return clt
}

func WithURL(url string) DialOption {
	return func(clt *Client) {
		clt.url = url
	}
}

func (c *Client) GetApplications(ctx context.Context) (*GetApplicationsResponse, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.url+appURI, nil)
	if err != nil {
		return nil, fmt.Errorf("cannot create new request: %w", err)
	}
	req.Header.Set("Accept", "application/json")
	if err != nil {
		return nil, fmt.Errorf("cannot set header: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("cannot send request: %w", err)
	}
	if resp.StatusCode >= 400 {
		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("cannot read body from response: %w", err)
		}
		return nil, fmt.Errorf("response with error: %s", string(b))
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Printf("cannot close response body: %v", err) // TODO: log using zap
			return
		}
	}(resp.Body)

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("cannot read body from response: %w", err)
	}

	var apps GetApplicationsResponse
	err = json.Unmarshal(b, &apps)
	if err != nil {
		return nil, fmt.Errorf("cannot unmarshal to applications")
	}

	return &apps, nil
}

// Register 注册实例
// POST /eureka/v2/apps/appID
func (c *Client) Register(ctx context.Context) error {
	// instance 服务实例
	type InstanceInfo struct {
		Instance *Instance `json:"instance"`
	}
	var info = &InstanceInfo{
		Instance: c.instance,
	}

	b, err := json.Marshal(info)
	if err != nil {
		return fmt.Errorf("cannot marshal instance: %v", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.url+appURI+c.instance.App, bytes.NewReader(b))
	if err != nil {
		return fmt.Errorf("cannot create new request: %w", err)
	}
	req.Header.Add("Content-Type", "application/json")
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("cannot send request: %w", err)
	}
	if resp.StatusCode >= 400 {
		errRes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("cannot read body from response: %w", err)
		}
		return fmt.Errorf("cannot register app: %s", errRes)
	}

	return nil
}

// Heartbeat 发送心跳
// PUT /eureka/v2/apps/appID/instanceID
func (c *Client) Heartbeat(ctx context.Context) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, c.url+appURI+c.instance.App+"/"+c.instance.InstanceId, bytes.NewReader([]byte{'O', 'K'}))
	if err != nil {
		return fmt.Errorf("cannot create new request: %w", err)
	}
	err = req.ParseForm()
	if err != nil {
		return fmt.Errorf("cannot parse form: %w", err)
	}
	req.Form.Add("status", "UP")
	req.Header.Add("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("cannot send request: %w", err)
	}
	if resp.StatusCode == http.StatusNotFound {
		return fmt.Errorf(http.StatusText(http.StatusNotFound))
	}
	if resp.StatusCode >= 400 {
		errRes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("cannot read body from response: %w", err)
		}
		return fmt.Errorf("app failed to heartbeat: %v", errRes)
	}
	return nil
}

func (c *Client) Run(ctx context.Context) error {
	err := c.Register(ctx)
	if err != nil {
		return fmt.Errorf("failed to register: %v", err)
	}

	for {
		select {
		case <-c.ticker.C:
			err := c.Heartbeat(ctx)
			if err != nil {
				return fmt.Errorf("failed to heartbeat: %w", err)
			}
		case <-c.tickerStop:
			c.stop <- struct{}{}
			return nil
		}
	}
}

func (c *Client) Shutdown(ctx context.Context) {
	c.ticker.Stop()
	c.tickerStop <- struct{}{}
	close(c.tickerStop)

	select {
	case <-c.stop:
	case <-ctx.Done():
	}
}
