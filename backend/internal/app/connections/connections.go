package connections

import (
	"fmt"
	"context"
	"github.com/jackc/pgx/v5"
	"gitlab.com/w0ikid/study-platform/internal/app/config"
)
type Connections struct {
	DB *pgx.Conn
}

func (c *Connections) Close() {
	c.DB.Close(context.Background())
}

func NewConnections(cfg *config.Config) (*Connections, error) {
	conn, err := pgx.Connect(context.Background(), cfg.DB.GetDBConnString())
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}
	return &Connections{DB: conn}, nil
}