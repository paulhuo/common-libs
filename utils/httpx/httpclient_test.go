package httpclient

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gitlab.yx/base-service/common-go/utils/logx"
)

func TestPostJSON(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	}
	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	// 设置日志logger
	logx.SetWriter(logx.NewWriter(os.Stdout))

	client := NewClient()
	resp, err := client.PostJSON(context.Background(), server.URL+"/json", map[string]string{"foo": "bar"})
	assert.NoError(t, err)
	assert.Contains(t, string(resp), "ok")
}

func TestPostJSON_500(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`"error":"internal server error"}`))
	}
	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	client := NewClient()
	_, err := client.PostJSON(context.Background(), server.URL+"/json", map[string]string{"foo": "bar"})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "status 500")
}

func TestPostJSON_Timeout(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(2 * time.Second)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	}
	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	client := NewClient()
	client.SetTimeout(1) // 设置超时时间为1秒
	_, err := client.PostJSON(context.Background(), server.URL+"/json", map[string]string{"foo": "bar"})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "context deadline exceeded")
}

func TestGet(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "v", r.URL.Query().Get("k"))
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"result":"ok"}`))
	}
	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	client := NewClient()
	resp, err := client.Get(context.Background(), server.URL+"/", map[string]string{"k": "v"})
	assert.NoError(t, err)
	assert.Contains(t, string(resp), "ok")
}

func TestPostFormURLEncoded(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		assert.Contains(t, string(body), "name=go")
		assert.Contains(t, string(body), "ids%5B%5D=1") // ids[]=1
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`ok`))
	}
	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	client := NewClient()
	form := map[string]interface{}{
		"name": "go",
		"ids":  []int{1, 2},
	}
	resp, err := client.PostForm(context.Background(), server.URL+"/form", form)
	assert.NoError(t, err)
	assert.Contains(t, string(resp), "ok")
}

func TestPostFormMultipart(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		assert.Contains(t, r.Header.Get("Content-Type"), "multipart/form-data")
		err := r.ParseMultipartForm(10 << 20)
		assert.NoError(t, err)
		assert.Equal(t, "file.txt", r.MultipartForm.File["file"][0].Filename)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("uploaded"))
	}
	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	// 创建临时文件
	f, _ := os.CreateTemp("", "file.txt")
	f.WriteString("hello")
	f.Close()
	defer os.Remove(f.Name())

	client := NewClient()
	form := map[string]interface{}{
		"file": &FormFile{
			Path:     f.Name(),
			FileName: "file.txt",
		},
	}
	resp, err := client.PostForm(context.Background(), server.URL+"/upload", form)
	assert.NoError(t, err)
	assert.Contains(t, string(resp), "uploaded")
}

// test set request ID header
func TestSetRequestID(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		assert.NotEmpty(t, r.Header.Get("X-Request-ID"))
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	}
	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	client := NewClient()
	ctx := context.Background()
	resp, err := client.PostJSON(ctx, server.URL+"/json", map[string]string{"foo": "bar"})
	assert.NoError(t, err)
	assert.Contains(t, string(resp), "ok")
}
