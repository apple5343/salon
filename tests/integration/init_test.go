package integration

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"os"
	"salon/tests/integration/httpclient"
	"salon/tests/integration/models"
	"testing"
	"time"

	dockercontainer "github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/ilyakaznacheev/cleanenv"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go/modules/compose"
	"github.com/testcontainers/testcontainers-go/wait"
)

const (
	Timeout = 10 * time.Minute
)

type TestConfig struct {
	PostgresUser          string `env:"POSTGRES_USER" required:"true"`
	PostgresPassword      string `env:"POSTGRES_PASSWORD" required:"true"`
	PostgresDB            string `env:"POSTGRES_DB" required:"true"`
	PostgresPort          string `env:"POSTGRES_PORT" required:"true"`
	PostgresDSN           string `env:"POSTGRES_DSN" required:"true"`
	ServerAddr            string `env:"SERVER_ADDR" required:"true"`
	HttpServerAddr        string `env:"HTTP_SERVER_ADDR" required:"true"`
	AdminEmail            string `env:"ADMIN_EMAIL" required:"true"`
	AdminPassword         string `env:"ADMIN_PASSWORD" required:"true"`
	AdminPhone            string `env:"ADMIN_PHONE" required:"true"`
	AdminFullName         string `env:"ADMIN_NAME" required:"true"`
	AdminPassportSeries   string `env:"ADMIN_PASSPORT_SERIES" required:"true"`
	AdminPassportNumber   string `env:"ADMIN_PASSPORT_NUMBER" required:"true"`
	AdminPassportIssuedBy string `env:"ADMIN_PASSPORT_ISSUED_BY" required:"true"`
}

type BaseTestSuite struct {
	suite.Suite
	config      *TestConfig
	compose     compose.ComposeStack
	composePath string
	db          *sql.DB
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

func (suite *BaseTestSuite) RunGroupWithLogs(testGroup string, run func()) {
	suite.MustRestartApp()
	since := time.Now()

	run()

	suite.SaveAppLogs(testGroup, since)
}

func (suite *BaseTestSuite) MustRestartApp() {
	c, err := suite.compose.ServiceContainer(suite.ctx, "app-test")
	suite.Require().NoError(err)

	err = c.Stop(suite.ctx, nil)
	suite.Require().NoError(err)

	err = c.Start(suite.ctx)
	suite.Require().NoError(err)

	err = wait.NewHealthStrategy().WaitUntilReady(suite.ctx, c)
	suite.Require().NoError(err)
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
}

func (suite *BaseTestSuite) TearDownSuite() {
	suite.cancel()
	suite.compose.Down(suite.ctx)
}

func (suite *BaseTestSuite) GetEnv() map[string]string {
	return map[string]string{
		"POSTGRES_USER":            suite.config.PostgresUser,
		"POSTGRES_PASSWORD":        suite.config.PostgresPassword,
		"POSTGRES_DB":              suite.config.PostgresDB,
		"POSTGRES_PORT":            suite.config.PostgresPort,
		"POSTGRES_DSN":             suite.config.PostgresDSN,
		"SERVER_ADDR":              suite.config.ServerAddr,
		"HTTP_SERVER_ADDR":         suite.config.HttpServerAddr,
		"ADMIN_EMAIL":              suite.config.AdminEmail,
		"ADMIN_PASSWORD":           suite.config.AdminPassword,
		"ADMIN_PHONE":              suite.config.AdminPhone,
		"ADMIN_NAME":               suite.config.AdminFullName,
		"ADMIN_PASSPORT_SERIES":    suite.config.AdminPassportSeries,
		"ADMIN_PASSPORT_NUMBER":    suite.config.AdminPassportNumber,
		"ADMIN_PASSPORT_ISSUED_BY": suite.config.AdminPassportIssuedBy,
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
