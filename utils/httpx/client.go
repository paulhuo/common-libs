package httpclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

/*
 * 该键名用于标识每个请求的唯一性，便于日志追踪和错误排查
 * 默认值为 "X-Request-ID"
 * 可以通过 SetRequestIdKey 方法自定义该键名
 */
const DefaultRequestIdKey = "X-Request-ID"

/*
 * DefaultRetryMaxRetries 定义默认的最大重试次数
 * 默认为 0，表示不进行重试
 */
const DefaultRetryMaxRetries = 0

/*
 * DefaultRetryInterval 定义默认的重试间隔时间
 * 默认为 1 秒
 */
const DefaultRetryInterval = 1 * time.Second

/*
 * DefaultTimeoutSeconds 定义默认的请求超时时间
 * 默认为 5 秒
 */
const DefaultTimeoutSeconds = 5

type Client struct {
	HTTPClient   *http.Client
	Headers      map[string]string
	requestIdKey string

	retryCfg RetryConfig
}

func NewClient() *Client {
	timeout := time.Duration(DefaultTimeoutSeconds) * time.Second
	return &Client{
		HTTPClient: &http.Client{
			Timeout: timeout,
		},
		Headers:      make(map[string]string),
		requestIdKey: DefaultRequestIdKey,
		retryCfg:     RetryConfig{MaxRetries: DefaultRetryMaxRetries, Interval: DefaultRetryInterval},
	}
}

func (c *Client) WithRetry(maxRetries int, interval time.Duration) *Client {
	c.retryCfg = RetryConfig{MaxRetries: maxRetries, Interval: interval}
	return c
}

func (c *Client) SetHeader(key, value string) {
	c.Headers[key] = value
}

func (c *Client) SetRequestIdKey(key string) {
	c.requestIdKey = key
}

/*
 * SetTimeout 设置 HTTP 客户端的超时时间
 * 参数 timeout 单位为秒
 * 如果设置为 0，则不设置超时
 */
func (c *Client) SetTimeout(timeout uint) {
	c.HTTPClient.Timeout = time.Duration(timeout) * time.Second
}

func (c *Client) addHeaders(req *http.Request) {
	for k, v := range c.Headers {
		req.Header.Set(k, v)
	}
}

func (c *Client) doWithContext(ctx context.Context, req *http.Request) (body []byte, err error) {
	traceID := GetTraceID(ctx)
	req = req.WithContext(ctx)

	if req.Header.Get(c.requestIdKey) == "" {
		req.Header.Set(c.requestIdKey, traceID)
	}

	body, err = c.withRetry(func() ([]byte, error) {
		resp, err := c.HTTPClient.Do(req)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		if resp.StatusCode >= 300 {
			return nil, fmt.Errorf("status %d: %s", resp.StatusCode, string(body))
		}
		return body, nil
	})

	return body, err
}

func (c *Client) PostJSON(ctx context.Context, fullURL string, body interface{}) (content []byte, err error) {
	start := time.Now()
	defer func() {
		if err != nil {
			logger(ctx).WithDuration(time.Since(start)).Errorf("Request: PostJSON %s %v, Error: %v", fullURL, body, err)
		} else {
			logger(ctx).WithDuration(time.Since(start)).Infof("Request: PostJSON %s %v, Response: %d, Body: %s", fullURL, body, http.StatusOK, removeNewline(content))
		}
	}()

	data, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", fullURL, bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	c.addHeaders(req)
	content, err = c.doWithContext(ctx, req)

	return content, err
}

func (c *Client) Get(ctx context.Context, fullURL string, query map[string]string) (content []byte, err error) {
	if len(query) > 0 {
		params := make([]string, 0)
		for k, v := range query {
			params = append(params, fmt.Sprintf("%s=%s", k, v))
		}
		fullURL += "?" + strings.Join(params, "&")
	}

	start := time.Now()
	defer func() {
		if err != nil {
			logger(ctx).WithDuration(time.Since(start)).Errorf("Request: Get %s, Error: %v", fullURL, err)
		} else {
			logger(ctx).WithDuration(time.Since(start)).Infof("Request: Get %s, Response: %d, Body: %s", fullURL, http.StatusOK, removeNewline(content))
		}
	}()

	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		return nil, err
	}
	c.addHeaders(req)
	content, err = c.doWithContext(ctx, req)

	return content, err
}

func (c *Client) PostForm(ctx context.Context, fullURL string, form map[string]interface{}) (content []byte, err error) {
	hasFile := false
	for _, v := range form {
		if _, ok := v.(*FormFile); ok {
			hasFile = true
			break
		}
	}

	start := time.Now()
	defer func() {
		if err != nil {
			logger(ctx).WithDuration(time.Since(start)).Errorf("Request: PostForm %s %v, Error: %v", fullURL, form, err)
		} else {
			logger(ctx).WithDuration(time.Since(start)).Infof("Request: PostForm %s %v, Response: %d, Body: %s", fullURL, form, http.StatusOK, removeNewline(content))
		}
	}()

	var req *http.Request

	if hasFile {
		body, contentType, err := buildMultipartBodyInterface(form)
		if err != nil {
			return nil, err
		}
		req, err = http.NewRequest("POST", fullURL, body)
		if err != nil {
			return nil, err
		}
		req.Header.Set("Content-Type", contentType)
	} else {
		encoded, err := buildURLEncodedBodyInterface(form)
		if err != nil {
			return nil, err
		}
		req, err = http.NewRequest("POST", fullURL, strings.NewReader(encoded))
		if err != nil {
			return nil, err
		}
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}

	c.addHeaders(req)
	content, err = c.doWithContext(ctx, req)

	return content, err
}
