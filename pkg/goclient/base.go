package goclient

import (
	"fmt"
	"net/http"
)

type Client struct {
	remoteHost string
	port       int
	client     *http.Client
}

func (c *Client) baseURL() string {
	return fmt.Sprintf("http://%s:%d", c.remoteHost, c.port)
}
