package config

import (
	"os"
	"strconv"
	"testing"

	"github.com/joho/godotenv"
	. "github.com/smartystreets/goconvey/convey"
)

func TestNewConfig(t *testing.T) {
	Convey("Given a set of environment variables", t, func() {
		err := godotenv.Load("../.env.test")
		if err != nil {
			t.Fatalf("Error loading .env.test file")
		}

		expectedPort := "8080"
		expectedInterval := "30s"
		expectedEthRpcUrl := "http://localhost:8545"
		expectedCacheCap := "100"
		expectedProxyNftAddress := "0xa2133A0F6B8A70D70E75238C7C0d84A0cD3F1Db3"
		expectedDBHost := "eth-proxy-pg"
		expectedDBPort := 5432
		expectedDBUser := "user"
		expectedDBPassword := "password"
		expectedDBName := "PROXY_DB"
		expectedDBSslMode := "disable"

		os.Setenv("PORT", expectedPort)
		os.Setenv("INTERVAL", expectedInterval)
		os.Setenv("ETH_RPC_URL", expectedEthRpcUrl)
		os.Setenv("CACHE_CAP", expectedCacheCap)
		os.Setenv("PROXY_NFT_ADDRESS", expectedProxyNftAddress)
		os.Setenv("DB_USER", expectedDBUser)
		os.Setenv("DB_PASSWORD", expectedDBPassword)
		os.Setenv("DB_NAME", expectedDBName)
		os.Setenv("DB_HOST", expectedDBHost)
		os.Setenv("DB_PORT", strconv.Itoa(expectedDBPort))
		os.Setenv("DB_SSL", expectedDBSslMode)

		Convey("When NewConfig is called", func() {
			config := NewConfig()

			Convey("Then it should return a Config struct with the correct values", func() {
				So(config, ShouldNotBeNil)
				So(config.Port, ShouldEqual, expectedPort)
				So(config.Interval, ShouldEqual, expectedInterval)
				So(config.EthRpcUrl, ShouldEqual, expectedEthRpcUrl)
				So(config.CacheCap, ShouldEqual, expectedCacheCap)
				So(config.ProxyNftAddress, ShouldEqual, expectedProxyNftAddress)

				So(config.DB, ShouldNotBeNil)
				So(config.DB.Host, ShouldEqual, expectedDBHost)
				So(config.DB.Port, ShouldEqual, expectedDBPort)
				So(config.DB.User, ShouldEqual, expectedDBUser)
				So(config.DB.Password, ShouldEqual, expectedDBPassword)
				So(config.DB.Dbname, ShouldEqual, expectedDBName)
				So(config.DB.SslMode, ShouldEqual, expectedDBSslMode)
			})
		})
	})
}
