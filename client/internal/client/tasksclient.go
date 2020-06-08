package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/go-playground/validator/v10"
	"github.com/pkg/errors"
	"github.com/tsvobo/go-workshop-http/client/internal/model"
)

var validate = validator.New()

type unexpectedStatusCode struct {
	statusCode int
	body       string
}

func (err *unexpectedStatusCode) Error() string {
	return fmt.Sprintf("Unexpected status code '%v' with body '%v'.", err.statusCode, err.body)
}

type Task struct {
	baseURL *url.URL
	client  *http.Client
}

func NewTask(baseUrl string, client *http.Client) (Task, error) {
	u, err := url.Parse(baseUrl)
	if err != nil {
		return Task{}, errors.Wrap(err, "unable to initialize HTTP client")
	}
	return Task{u, client}, nil
}

// TODO TASK-3: Implement task creation using POST request
func (c *Task) Create(ctx context.Context, task model.Task) (model.Task, error) {
	return model.Task{}, nil
}

// TODO TASK-4: Implement task retrieval using GET request
func (c *Task) Find(ctx context.Context, id string) (model.Task, error) {
	return model.Task{}, nil
}

func (c *Task) newRequest(method, path string, body interface{}) (*http.Request, error) {
	p := &url.URL{Path: path}
	u := c.baseURL.ResolveReference(p)

	var buf io.ReadWriter
	if body != nil {
		buf = new(bytes.Buffer)
		if err := json.NewEncoder(buf).Encode(body); err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, err
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("Accept", "application/json")
	return req, nil
}

func (c *Task) do(ctx context.Context, req *http.Request, val interface{}) error {
	req = req.WithContext(ctx)

	resp, err := c.client.Do(req)
	if err != nil {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			return err
		}
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || 300 <= resp.StatusCode {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		return &unexpectedStatusCode{
			statusCode: resp.StatusCode,
			body:       string(body),
		}
	}

	if val != nil {
		return json.NewDecoder(resp.Body).Decode(val)
	}
	return nil
}
