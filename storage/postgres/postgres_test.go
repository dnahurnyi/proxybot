//go:build integration
// +build integration

package postgres_test

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/dnahurnyi/proxybot/storage/postgres"
	"github.com/golang-migrate/migrate/v4"
	"github.com/jmoiron/sqlx"
	"github.com/myles-mcdonnell/blondie"
	"github.com/stretchr/testify/suite"
	gorm_postgres "gorm.io/driver/postgres"
	"gorm.io/gorm"

	// db driver
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

type clientConfig struct {
	Host     string
	Port     int
	Db       string
	User     string
	Password string
	Ssl      bool
}

type RepoTestSuite struct {
	suite.Suite
	pgCfg *clientConfig
	db    *sqlx.DB
	repo  *postgres.Repository
}

func (config *clientConfig) Address() string {
	address := fmt.Sprintf("postgres://%v:%v@%v:%v/%v?sslmode=", config.User, config.Password, config.Host, config.Port, config.Db)
	if !config.Ssl {
		address += "disable"
	}

	return address
}

// ApplySchemaMigration Bootstrap creates and migrates db schema versions
func ApplySchemaMigration(config *clientConfig) error {
	dbConn, err := ConnectClientWithRetry(config, 0, 3)
	if err != nil {
		return err
	}

	db, err := dbConn.DB()
	if err != nil {
		return err
	}

	err = db.Close()
	if err != nil {
		return err
	}

	m, err := migrate.New(
		"file://db/migrations",
		config.Address(),
	)
	if err != nil {
		return err
	}
	defer m.Close()

	if err = migrateUp(m); err != nil {
		return err
	}

	_, _, err = m.Version()
	return err
}

type IMigrate interface {
	Up() error
}

func migrateUp(m IMigrate) error {
	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			return err
		}
	}

	return nil
}

var retryWait = 3 * time.Second

func ConnectClientWithRetry(cfg *clientConfig, preConnectTestTimeout time.Duration, retries int) (db *gorm.DB, err error) {
	blondieOpts := blondie.DefaultOptions()
	blondieOpts.QuietMode = true
	for retries > 0 {
		blondie.WaitForDeps(
			[]blondie.DepCheck{
				blondie.NewTcpCheck(cfg.Host, cfg.Port, preConnectTestTimeout),
			},
			blondieOpts,
		)
		db, err = gorm.Open(gorm_postgres.Open(cfg.Address()), &gorm.Config{})
		if err != nil {
			retries--
			time.Sleep(retryWait)
			continue
		}

		break
	}

	return db, err
}

func (s *RepoTestSuite) SetupSuite() {
	cfg := &clientConfig{
		User:     "user",
		Password: "password",
		Db:       "proxybot-db",
		Host:     "localhost",
		Port:     5432,
		Ssl:      false,
	}

	dbURL := os.Getenv("DB_URL")
	if len(dbURL) != 0 {
		urlParsed, err := url.Parse(dbURL)
		if err != nil {
			s.T().Fatal(fmt.Errorf("can't parse dbURL '%s': %w", dbURL, err))
		}
		cfg.Host = urlParsed.Hostname()
		cfg.Port, err = strconv.Atoi(urlParsed.Port())
		if err != nil {
			s.T().Fatal(fmt.Errorf("invalid port '%s': %w", urlParsed.Port(), err))
		}
		cfg.User = urlParsed.User.Username()
		if pass, ok := urlParsed.User.Password(); ok {
			cfg.Password = pass
		}
		cfg.Db = strings.ReplaceAll(urlParsed.Path, "/", "")
	}

	s.pgCfg = cfg
	var err error
	dbConn, err := gorm.Open(gorm_postgres.Open(cfg.Address()), &gorm.Config{})
	if err != nil {
		s.T().Fatal(fmt.Errorf("connect to postgres: %w", err))
	}
	db, err := dbConn.DB()
	if err != nil {
		s.T().Fatal(fmt.Errorf("get db instance from DB connection: %w", err))
	}
	s.db = sqlx.NewDb(db, "postgres").Unsafe()
	err = s.db.Ping()
	if err != nil {
		s.T().Fatal(fmt.Errorf("connecting to test db with '%s': %w", cfg.Address(), err))
	}

	if err = os.Chdir("../.."); err != nil {
		s.T().Fatal(fmt.Errorf("change directory: %w", err))
	}

	if err = ApplySchemaMigration(cfg); err != nil {
		s.T().Fatal(fmt.Errorf("db schema migration: %w", err))
	}

	if s.repo, err = postgres.New(dbConn); err != nil {
		s.T().Fatal(fmt.Errorf("db init: %w", err))
	}
}

func (s *RepoTestSuite) TearDownSuite() {
	err := s.db.Close()
	if err != nil {
		s.T().Fatal(fmt.Errorf("db close failed: %w", err))
	}
}

func (s *RepoTestSuite) TearDownTest() {
	_, err := s.db.Exec(`TRUNCATE subscriptions, tags;`)
	if err != nil {
		s.T().Fatal(fmt.Errorf("test cleanup failed: %w", err))
	}
}

func TestRepoTestSuite(t *testing.T) {
	suite.Run(t, new(RepoTestSuite))
}

func (s *RepoTestSuite) Test_Migrations() {
	m, err := migrate.New("file://db/migrations", s.pgCfg.Address())
	if err != nil {
		s.T().Fatal(fmt.Errorf("migrate new: %w", err))
	}
	defer m.Close()

	if err = m.Down(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			s.Nil(err)
		}
	}

	_, _, err = m.Version()
	s.Equal(migrate.ErrNilVersion, err)

	if err = ApplySchemaMigration(s.pgCfg); err != nil {
		s.T().Fatal(fmt.Errorf("db schema migration: %w", err))
	}
}
