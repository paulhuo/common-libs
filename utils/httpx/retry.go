
package httpclient

import (
    "time"
)

type RetryConfig struct {
    MaxRetries int
    Interval   time.Duration
}

func (c *Client) withRetry(doFunc func() ([]byte, error)) ([]byte, error) {
    var resp []byte
    var err error

    for i := 0; i <= c.retryCfg.MaxRetries; i++ {
        resp, err = doFunc()
        if err == nil {
            return resp, nil
        }
        time.Sleep(c.retryCfg.Interval)
    }
    return nil, err
}
