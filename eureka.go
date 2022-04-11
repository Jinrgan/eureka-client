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
	"strings"
	"time"
)

const appURI = "/apps/"

// 调用 eureka 服务端 rest API
// https://github.com/Netflix/eureka/wiki/Eureka-REST-operations

type InstanceOption func(ins *Instance)

func NewInstance(app, ip string, port int, opts ...InstanceOption) *Instance {
	ins := &Instance{
		InstanceId:       fmt.Sprintf("%s:%s:%d", ip, app, port),
		HostName:         ip,
		App:              strings.ToLower(app),
		Status:           "UP",                   // TODO: enum
		OverriddenStatus: "UNKNOWN",              // TODO: enum
		Port:             &Port{Enabled: "true"}, // TODO: bool
		SecurePort:       nil,
		CountryId:        0,
		DataCenterInfo: &DataCenterInfo{
			Class: "MyOwn",
			Name:  "com.netflix.appinfo.InstanceInfo$DefaultDataCenterInfo",
		},
		LeaseInfo: &LeaseInfo{
			RenewalIntervalInSecs: 30,
			DurationInSecs:        15,
		},
	}

	for _, opt := range opts {
		opt(ins)
	}

	return ins
}

type Client struct {
	url                          string
	RegistryFetchIntervalSeconds time.Duration
	http.Client
}

type DialOption func(clt *Client)

func Dial(opts ...DialOption) *Client {
	clt := &Client{
		url:                          "http://admin:admin@localhost:8761/eureka/",
		RegistryFetchIntervalSeconds: 15 * time.Second,
	}

	for _, opt := range opts {
		opt(clt)
	}

	return clt
}

func WithURL(url string) DialOption {
	return func(s *Client) {
		s.url = url
	}
}

func (c *Client) GetApplications(ctx context.Context) (*GetApplicationsResponse, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.url+appURI, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	if err != nil {
		return nil, fmt.Errorf("cannot set header: %w", err)
	}

	resp, err := c.Do(req)
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
func (c *Client) Register(ctx context.Context, ins *Instance) error {
	b, err := json.Marshal(ins)
	if err != nil {
		return fmt.Errorf("cannot marshal instance: %v", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.url+appURI, bytes.NewReader(b))
	if err != nil {
		return err
	}
	if req.Response.StatusCode >= http.StatusBadRequest {
		errRes, err := ioutil.ReadAll(req.Body)
		if err != nil {
			return fmt.Errorf("cannot read body from response: %w", err)
		}
		return fmt.Errorf("cannot register app: %v", errRes)
	}

	return nil
}
