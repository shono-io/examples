package arangodb

import (
	"context"
	"github.com/arangodb/go-driver"
)

type Cursor struct {
	ctx context.Context
	c   driver.Cursor
}

func (c *Cursor) HasNext() bool {
	return c.c.HasMore()
}

func (c *Cursor) Next(v any) error {
	_, err := c.c.ReadDocument(c.ctx, v)
	if err != nil {
		if driver.IsNoMoreDocuments(err) {
			return nil
		} else if err != nil {
			return err
		}
	}

	return nil
}

func (c *Cursor) Count() int64 {
	return c.c.Statistics().FullCount()
}

func (c *Cursor) Close() error {
	return c.c.Close()
}
