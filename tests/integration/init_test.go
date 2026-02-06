package integration

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"io"
	"os"
	"salon/tests/integration/httpclient"
	"salon/tests/integration/models"
	"testing"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go/modules/compose"
	"github.com/testcontainers/testcontainers-go/wait"
)

const (
	Timeout = 5 * time.Minute
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
	t := &ClientSuite{base: s}
	suite.Run(s.T(), t)
	s.SaveAppLogs("client")
}

func (s *BaseTestSuite) TestEmployee() {
	t := &EmployeeSuite{base: s}
	suite.Run(s.T(), t)
	s.SaveAppLogs("employee")
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

func (suite *BaseTestSuite) SaveAppLogs(testGroup string) {
	logFile, err := os.Create(fmt.Sprintf("logs/%s.txt", testGroup))
	suite.Require().NoError(err)
	defer logFile.Close()
	container, err := suite.compose.ServiceContainer(suite.ctx, "app-test")
	suite.Require().NoError(err)

	reader, err := container.Logs(suite.ctx)
	suite.Require().NoError(err)

	_, err = io.Copy(logFile, reader)
	suite.Require().NoError(err)
}
