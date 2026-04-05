package integration

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"os"
	"salon/tests/integration/httpclient"
	"salon/tests/integration/models"
	"strconv"
	"strings"
	"testing"
	"time"

	dockercontainer "github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/ilyakaznacheev/cleanenv"
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go/modules/compose"
	"github.com/testcontainers/testcontainers-go/wait"
)

const (
	Timeout = 10 * time.Minute
)

type TestConfig struct {
	PostgresUser          string `env:"POSTGRES_USER" env-required:"true"`
	PostgresPassword      string `env:"POSTGRES_PASSWORD" env-required:"true"`
	PostgresDB            string `env:"POSTGRES_DB" env-required:"true"`
	PostgresPort          string `env:"POSTGRES_PORT" renv-equired:"true"`
	PostgresDSN           string `env:"POSTGRES_DSN" env-required:"true"`
	RedisPort             string `env:"REDIS_PORT" env-required:"true"`
	RedisAddr             string `env:"REDIS_ADDR" env-required:"true"`
	RedisPassword         string `env:"REDIS_PASSWORD" env-required:"true"`
	RedisDB               int    `env:"REDIS_DB" env-required:"true"`
	ServerAddr            string `env:"SERVER_ADDR" env-required:"true"`
	HttpServerAddr        string `env:"HTTP_SERVER_ADDR" env-required:"true"`
	AdminEmail            string `env:"ADMIN_EMAIL" env-required:"true"`
	AdminPassword         string `env:"ADMIN_PASSWORD" env-required:"true"`
	AdminPhone            string `env:"ADMIN_PHONE" env-required:"true"`
	AdminFullName         string `env:"ADMIN_NAME" env-required:"true"`
	AdminPassportSeries   string `env:"ADMIN_PASSPORT_SERIES" env-required:"true"`
	AdminPassportNumber   string `env:"ADMIN_PASSPORT_NUMBER" env-required:"true"`
	AdminPassportIssuedBy string `env:"ADMIN_PASSPORT_ISSUED_BY" env-required:"true"`
	RefreshSecret         string `env:"REFRESH_SECRET" env-required:"true"`
	AccessSecret          string `env:"ACCESS_SECRET" env-required:"true"`
}

type BaseTestSuite struct {
	suite.Suite
	config      *TestConfig
	compose     compose.ComposeStack
	composePath string
	db          *sql.DB
	rdb         *redis.Client
	client      *httpclient.Client
	ctx         context.Context
	cancel      context.CancelFunc
	adminCreds  *models.EmployeeCreds
}

func TestBaseSuite(t *testing.T) {
	suite.Run(t, new(BaseTestSuite))
}

func (s *BaseTestSuite) TestClient() {
	s.RunGroupWithLogs("client", func() {
		t := &ClientSuite{base: s}
		suite.Run(s.T(), t)
	})
}

func (s *BaseTestSuite) TestEmployee() {
	s.RunGroupWithLogs("employee", func() {
		t := &EmployeeSuite{base: s}
		suite.Run(s.T(), t)
	})
}

func (s *BaseTestSuite) TestSupplier() {
	s.RunGroupWithLogs("supplier", func() {
		t := &SupplierSuite{base: s}
		suite.Run(s.T(), t)
	})
}

func (s *BaseTestSuite) TestBrand() {
	s.RunGroupWithLogs("brand", func() {
		t := &BrandSuite{base: s}
		suite.Run(s.T(), t)
	})
}

func (s *BaseTestSuite) TestModel() {
	s.RunGroupWithLogs("model", func() {
		t := &ModelSuite{base: s}
		suite.Run(s.T(), t)
	})
}

func (s *BaseTestSuite) TestCar() {
	s.RunGroupWithLogs("car", func() {
		t := &CarSuite{base: s}
		suite.Run(s.T(), t)
	})
}

func (s *BaseTestSuite) TestSale() {
	s.RunGroupWithLogs("sale", func() {
		t := &SaleSuite{base: s}
		suite.Run(s.T(), t)
	})
}

func (suite *BaseTestSuite) RunGroupWithLogs(testGroup string, run func()) {
	since := suite.MustRestartApp()

	run()

	suite.SaveAppLogs(testGroup, since)
}

func (suite *BaseTestSuite) MustRestartApp() time.Time {
	c, err := suite.compose.ServiceContainer(suite.ctx, "app-test")
	suite.Require().NoError(err)

	err = c.Stop(suite.ctx, nil)
	suite.Require().NoError(err)

	startedAt := time.Now()
	err = c.Start(suite.ctx)
	suite.Require().NoError(err)

	err = wait.NewHealthStrategy().WaitUntilReady(suite.ctx, c)
	if err != nil {
		logs := suite.readContainerLogs(c.GetContainerID(), startedAt)
		suite.Require().NoError(fmt.Errorf("%w\n--- app-test logs ---\n%s", err, logs))
	}

	return startedAt
}

func (suite *BaseTestSuite) readContainerLogs(containerID string, since time.Time) string {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return "failed to create docker client: " + err.Error()
	}
	defer cli.Close()

	reader, err := cli.ContainerLogs(suite.ctx, containerID, dockercontainer.LogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Since:      fmt.Sprintf("%d.%09d", since.Unix(), int64(since.Nanosecond())),
	})
	if err != nil {
		return "failed to read container logs: " + err.Error()
	}
	defer reader.Close()

	var buf strings.Builder
	stdcopy.StdCopy(&buf, &buf, reader)
	return buf.String()
}

func (suite *BaseTestSuite) SetupSuite() {
	var envPath, composePath, migrationsPath string
	flag.StringVar(&envPath, "env", "../../config/test.env", "Path to test env file")
	flag.StringVar(&composePath, "compose", "../../docker-compose.test.yml", "Path to compose file")
	flag.StringVar(&migrationsPath, "migrations", "../../migrations", "Path to migrations dir")
	flag.Parse()

	ctx, cancel := context.WithTimeout(context.Background(), Timeout)
	suite.ctx = ctx
	suite.cancel = cancel

	suite.MustInitConfig(envPath)
	suite.MustInitCompose(composePath)
	suite.compose.Up(ctx)
	suite.MustInitDB()
	suite.composePath = composePath
	suite.client = httpclient.NewClient(suite.config.ServerAddr)
	suite.adminCreds = &models.EmployeeCreds{
		Email:    suite.config.AdminEmail,
		Password: suite.config.AdminPassword,
	}
}

func (suite *BaseTestSuite) SetupTest() {
	tables := []string{"employees", "clients", "brands", "models", "suppliers", "cars", "sales", "events"}
	for _, table := range tables {
		query := fmt.Sprintf("TRUNCATE TABLE %s RESTART IDENTITY CASCADE", table)
		_, err := suite.db.Exec(query)
		suite.Require().NoError(err)
	}
	c, err := suite.compose.ServiceContainer(suite.ctx, "admin-registry-test")
	suite.Require().NoError(err)
	err = c.Start(suite.ctx)
	suite.Require().NoError(err)
	err = wait.ForExit().WaitUntilReady(suite.ctx, c)
	suite.Require().NoError(err)

	err = suite.rdb.FlushAll(suite.ctx).Err()
	suite.Require().NoError(err)
}

func (suite *BaseTestSuite) MustInitConfig(envPath string) {
	var config TestConfig
	if err := cleanenv.ReadConfig(envPath, &config); err != nil {
		suite.FailNow("read config: " + err.Error())
	}
	suite.config = &config
}

func (suite *BaseTestSuite) MustInitCompose(composePath string) {
	c, err := compose.NewDockerCompose(composePath)
	if err != nil {
		suite.FailNow("init compose: " + err.Error())
	}
	suite.compose = c.WithEnv(suite.GetEnv()).WaitForService("app-test", wait.NewHealthStrategy())
}

func (suite *BaseTestSuite) MustInitDB() {
	db, err := sql.Open("postgres", fmt.Sprintf("postgres://%s:%s@localhost:%s/%s?sslmode=disable",
		suite.config.PostgresUser, suite.config.PostgresPassword, suite.config.PostgresPort, suite.config.PostgresDB))
	if err != nil {
		suite.FailNow("init db: " + err.Error())
	}
	suite.db = db
	suite.rdb = redis.NewClient(&redis.Options{
		Addr:     "localhost:" + suite.config.RedisPort,
		Password: suite.config.RedisPassword,
		DB:       suite.config.RedisDB,
	})
	if err = suite.rdb.Ping(suite.ctx).Err(); err != nil {
		suite.FailNow("init redis: " + err.Error())
	}
}

func (suite *BaseTestSuite) TearDownSuite() {
	suite.cancel()
	suite.compose.Down(suite.ctx)
	suite.db.Close()
	suite.rdb.Close()
}

func (suite *BaseTestSuite) GetEnv() map[string]string {
	return map[string]string{
		"POSTGRES_USER":            suite.config.PostgresUser,
		"POSTGRES_PASSWORD":        suite.config.PostgresPassword,
		"POSTGRES_DB":              suite.config.PostgresDB,
		"POSTGRES_PORT":            suite.config.PostgresPort,
		"POSTGRES_DSN":             suite.config.PostgresDSN,
		"REDIS_PORT":               suite.config.RedisPort,
		"REDIS_ADDR":               suite.config.RedisAddr,
		"REDIS_PASSWORD":           suite.config.RedisPassword,
		"REDIS_DB":                 strconv.Itoa(suite.config.RedisDB),
		"SERVER_ADDR":              suite.config.ServerAddr,
		"HTTP_SERVER_ADDR":         suite.config.HttpServerAddr,
		"ADMIN_EMAIL":              suite.config.AdminEmail,
		"ADMIN_PASSWORD":           suite.config.AdminPassword,
		"ADMIN_PHONE":              suite.config.AdminPhone,
		"ADMIN_NAME":               suite.config.AdminFullName,
		"ADMIN_PASSPORT_SERIES":    suite.config.AdminPassportSeries,
		"ADMIN_PASSPORT_NUMBER":    suite.config.AdminPassportNumber,
		"ADMIN_PASSPORT_ISSUED_BY": suite.config.AdminPassportIssuedBy,
		"REFRESH_SECRET":           suite.config.RefreshSecret,
		"ACCESS_SECRET":            suite.config.AccessSecret,
	}
}

func (suite *BaseTestSuite) SaveAppLogs(testGroup string, since time.Time) {
	logFile, err := os.Create(fmt.Sprintf("logs/%s.txt", testGroup))
	suite.Require().NoError(err)
	defer logFile.Close()

	c, err := suite.compose.ServiceContainer(suite.ctx, "app-test")
	suite.Require().NoError(err)

	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	suite.Require().NoError(err)
	defer cli.Close()

	reader, err := cli.ContainerLogs(suite.ctx, c.GetContainerID(), dockercontainer.LogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Since:      fmt.Sprintf("%d.%09d", since.Unix(), int64(since.Nanosecond())),
	})
	suite.Require().NoError(err)
	defer reader.Close()

	_, err = stdcopy.StdCopy(logFile, logFile, reader)
	suite.Require().NoError(err)
}