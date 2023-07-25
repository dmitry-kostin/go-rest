package ping_test

import (
	"github.com/cockroachdb/errors"
	"github.com/dmitry-kostin/go-rest/src/db"
	"github.com/dmitry-kostin/go-rest/src/pkg"
	"github.com/dmitry-kostin/go-rest/src/services/ping"
	"github.com/dmitry-kostin/go-rest/src/services/ping/models"
	"github.com/jackc/pgx/v5/pgxpool"
	"reflect"
	"testing"
)

func TestService_Ping(t *testing.T) {
	logger := pkg.NewEmptyLogger()
	config := pkg.NewConfig()

	whenConnectionEstablished := func() *pgxpool.Pool {
		return db.InitPostgresConnectionPool(config, logger)
	}

	whenConnectionClosed := func() *pgxpool.Pool {
		conn := db.InitPostgresConnectionPool(config, logger)
		conn.Close()
		return conn
	}

	tests := []struct {
		name    string
		conn    *pgxpool.Pool
		want    interface{}
		wantErr error
	}{
		{
			name: "when db connection established",
			conn: whenConnectionEstablished(),
			want: &models.Ping{
				Pong:    "You reached the destination. Pong.",
				Version: config.Version,
				Env:     config.AppEnv,
			},
			wantErr: nil,
		},
		{
			name:    "when db connection closed",
			conn:    whenConnectionClosed(),
			want:    nil,
			wantErr: pkg.ErrDatabaseError,
		},
	}

	for _, tt := range tests {
		s := ping.NewService(tt.conn, config)
		got, err := s.Ping(nil, nil)
		if tt.wantErr != nil {
			if errors.Is(err, pkg.ErrDatabaseError) {
				return
			}
			t.Errorf("Ping() error = %v, wantErr %v", err, tt.wantErr)
			return
		}
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("Ping() got = %v, want %v", got, tt.want)
		}
	}
}
