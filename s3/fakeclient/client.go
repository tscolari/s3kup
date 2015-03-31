package fakeclient

import "github.com/tscolari/s3up/s3"

type Client struct {
	StoreCall  func(path string, content []byte) error
	ListCall   func(path string) (s3.Versions, error)
	DeleteCall func(path string) error
}

func (c *Client) Store(path string, content []byte) error {
	if c.StoreCall != nil {
		return c.StoreCall(path, content)
	}
	return nil
}

func (c *Client) List(path string) (files s3.Versions, err error) {
	if c.ListCall != nil {
		return c.ListCall(path)
	}
	return nil, nil
}

func (c *Client) Delete(path string) error {
	if c.DeleteCall != nil {
		return c.DeleteCall(path)
	}
	return nil
}
