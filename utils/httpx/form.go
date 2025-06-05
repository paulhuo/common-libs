package httpclient

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/url"
	"os"
	"path/filepath"
)

func buildURLEncodedBody(form map[string]FormValue) (string, error) {
	values := url.Values{}
	for key, v := range form {
		switch val := v.Value.(type) {
		case string, int, float64, int64:
			values.Add(key, fmt.Sprintf("%v", val))
		case []string:
			for _, item := range val {
				values.Add(key+"[]", item)
			}
		case []int:
			for _, item := range val {
				values.Add(key+"[]", fmt.Sprintf("%d", item))
			}
		case map[string]interface{}:
			for k2, v2 := range val {
				values.Add(fmt.Sprintf("%s[%s]", key, k2), fmt.Sprintf("%v", v2))
			}
		default:
			return "", fmt.Errorf("unsupported value type for key %s", key)
		}
	}
	return values.Encode(), nil
}

func buildMultipartBody(form map[string]FormValue) (*bytes.Buffer, string, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	for key, val := range form {
		if val.IsFile {
			file, err := os.Open(val.Value.(string))
			if err != nil {
				return nil, "", err
			}
			defer file.Close()

			part, err := writer.CreateFormFile(val.FieldName, val.FileName)
			if err != nil {
				return nil, "", err
			}
			if _, err := io.Copy(part, file); err != nil {
				return nil, "", err
			}
		} else {
			switch v := val.Value.(type) {
			case string, int, float64, int64:
				writer.WriteField(key, fmt.Sprintf("%v", v))
			case []string:
				for _, item := range v {
					writer.WriteField(key+"[]", item)
				}
			case map[string]interface{}:
				for k2, v2 := range v {
					writer.WriteField(fmt.Sprintf("%s[%s]", key, k2), fmt.Sprintf("%v", v2))
				}
			default:
				return nil, "", fmt.Errorf("unsupported multipart value type for key %s", key)
			}
		}
	}

	err := writer.Close()
	if err != nil {
		return nil, "", err
	}

	return body, writer.FormDataContentType(), nil
}

func buildURLEncodedBodyInterface(form map[string]interface{}) (string, error) {
	values := url.Values{}
	for key, val := range form {
		switch v := val.(type) {
		case string, int, float64, int64, int32, uint32, uint64:
			values.Add(key, fmt.Sprintf("%v", v))
		case []string:
			for _, item := range v {
				values.Add(key+"[]", item)
			}
		case []int:
			for _, item := range v {
				values.Add(key+"[]", fmt.Sprintf("%d", item))
			}
		case []int32:
			for _, item := range v {
				values.Add(key+"[]", fmt.Sprintf("%d", item))
			}
		case []int64:
			for _, item := range v {
				values.Add(key+"[]", fmt.Sprintf("%d", item))
			}
		case []uint32:
			for _, item := range v {
				values.Add(key+"[]", fmt.Sprintf("%d", item))
			}
		case []uint64:
			for _, item := range v {
				values.Add(key+"[]", fmt.Sprintf("%d", item))
			}
		case map[string]interface{}:
			for k2, v2 := range v {
				values.Add(fmt.Sprintf("%s[%s]", key, k2), fmt.Sprintf("%v", v2))
			}
		case []map[string]interface{}:
			for i, item := range v {
				for k2, v2 := range item {
					values.Add(fmt.Sprintf("%s[%d][%s]", key, i, k2), fmt.Sprintf("%v", v2))
				}
			}
		case *FormFile:
			continue // skip file in URL encoding
		default:
			return "", fmt.Errorf("unsupported type for key %s", key)
		}
	}
	return values.Encode(), nil
}

func buildMultipartBodyInterface(form map[string]interface{}) (*bytes.Buffer, string, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	for key, val := range form {
		switch v := val.(type) {
		case *FormFile:
			file, err := os.Open(v.Path)
			if err != nil {
				return nil, "", err
			}
			defer file.Close()
			field := v.Field
			if field == "" {
				field = key
			}
			filename := v.FileName
			if filename == "" {
				filename = filepath.Base(v.Path)
			}
			part, err := writer.CreateFormFile(field, filename)
			if err != nil {
				return nil, "", err
			}
			if _, err := io.Copy(part, file); err != nil {
				return nil, "", err
			}
		case string, int, float64:
			writer.WriteField(key, fmt.Sprintf("%v", v))
		case []string:
			for _, item := range v {
				writer.WriteField(key+"[]", item)
			}
		case map[string]interface{}:
			for k2, v2 := range v {
				writer.WriteField(fmt.Sprintf("%s[%s]", key, k2), fmt.Sprintf("%v", v2))
			}
		default:
			return nil, "", fmt.Errorf("unsupported multipart type for key %s", key)
		}
	}

	err := writer.Close()
	if err != nil {
		return nil, "", err
	}
	return body, writer.FormDataContentType(), nil
}

func hasFile(form map[string]FormValue) bool {
	for _, v := range form {
		if v.IsFile {
			return true
		}
	}
	return false
}
