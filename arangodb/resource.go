package arangodb

import (
	"context"
	"errors"
	"fmt"
	"github.com/arangodb/go-driver"
	"github.com/arangodb/go-driver/http"
	go_shono "github.com/shono-io/go-shono"
	"github.com/shono-io/go-shono/resources"
	"github.com/sirupsen/logrus"
)

type Opt func(config *Config)

func WithEndpoint(endpoint string) Opt {
	return func(config *Config) {
		config.Endpoints = []string{endpoint}
	}
}

func WithAuthentication(auth driver.Authentication) Opt {
	return func(config *Config) {
		config.Authentication = auth
	}
}

func WithDatabaseName(name string) Opt {
	return func(config *Config) {
		config.DatabaseName = name
	}
}

type Config struct {
	Endpoints      []string
	Authentication driver.Authentication
	DatabaseName   string
}

func (sc *Config) Validate() error {
	if sc.Endpoints == nil || len(sc.Endpoints) == 0 {
		return errors.New("no endpoints specified")
	}

	if sc.DatabaseName == "" {
		return errors.New("no database name specified")
	}

	return nil
}

func MustNewResource(id string, opts ...Opt) go_shono.Resource[any] {
	r, err := NewResource(id, opts...)
	if err != nil {
		panic(err)
	}

	return *r
}

func NewResource(id string, opts ...Opt) (*go_shono.Resource[any], error) {
	cfg := &Config{}
	for _, opt := range opts {
		opt(cfg)
	}

	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	conn, err := http.NewConnection(http.ConnectionConfig{
		Endpoints: cfg.Endpoints,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create arangodb connection: %w", err)
	}
	c, err := driver.NewClient(driver.ClientConfig{
		Connection:     conn,
		Authentication: cfg.Authentication,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create arangodb client: %w", err)
	}

	db, err := c.Database(context.Background(), cfg.DatabaseName)
	if err != nil {
		return nil, fmt.Errorf("failed to get arangodb database: %w", err)
	}

	return &go_shono.Resource[any]{
		Id: id,
		ClientFactory: func() any {
			return &Adb{db: db}
		},
	}, nil
}

type Adb struct {
	db driver.Database
}

func (a *Adb) MustQuery(ctx context.Context, query Query) resources.Cursor {
	c, err := a.Query(ctx, query)
	if err != nil {
		panic(err)
	}

	return c
}

func (a *Adb) Query(ctx context.Context, query Query) (resources.Cursor, error) {
	logrus.Debugf("QRY >> %s with %s", query.Statement(), query.Params())

	cursor, err := a.db.Query(driver.WithQueryFullCount(ctx), query.Statement(), query.params)
	if err != nil {
		return nil, err
	}
	defer cursor.Close()

	return &Cursor{c: cursor, ctx: ctx}, nil
}

func (a *Adb) MustGet(ctx context.Context, kind string, key string, target any) bool {
	exists, err := a.Get(ctx, kind, key, target)
	if err != nil {
		panic(err)
	}

	return exists
}

func (a *Adb) Get(ctx context.Context, kind string, key string, target any) (bool, error) {
	col, err := a.db.Collection(ctx, kind)
	if err != nil {
		return false, err
	}

	if _, err := col.ReadDocument(ctx, key, target); err != nil {
		if driver.IsNotFound(err) {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

func (a *Adb) MustSet(ctx context.Context, kind string, key string, value any) {
	err := a.Set(ctx, kind, key, value)
	if err != nil {
		panic(err)
	}
}

func (a *Adb) Set(ctx context.Context, kind string, key string, value any) error {
	col, err := a.db.Collection(ctx, kind)
	if err != nil {
		return err
	}

	// -- check if the entity already exists
	exists, err := col.DocumentExists(ctx, key)
	if err != nil {
		return fmt.Errorf("unable to check if document exists: %w", err)
	}
	if exists {
		if _, err := col.ReplaceDocument(ctx, key, value); err != nil {
			return err
		}
	} else {
		if _, err := col.CreateDocument(ctx, value); err != nil {
			return err
		}
	}

	return nil
}

func (a *Adb) MustDelete(ctx context.Context, kind string, key ...string) {
	if err := a.Delete(ctx, kind, key...); err != nil {
		panic(err)
	}
}

func (a *Adb) Delete(ctx context.Context, kind string, key ...string) error {
	col, err := a.db.Collection(ctx, kind)
	if err != nil {
		return err
	}

	_, _, err = col.RemoveDocuments(ctx, key)
	return err
}
