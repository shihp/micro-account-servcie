package main

import (
	"fmt"
	"micro-account-service/service"
	"micro-account-service/dbclient"
	"flag"
	"src/github.com/spf13/viper"
	"micro-account-service/config"
)

var appName = "accountservice"
// Init function, runs before main()
func init() {
	// Read command line flags
	profile := flag.String("profile", "test", "Environment profile, something similar to spring profiles")
	configServerUrl := flag.String("configServerUrl", "http://192.168.99.100:8888", "Address to config server")
	configBranch := flag.String("configBranch", "P8", "git branch to fetch configuration from")
	flag.Parse()

	// Pass the flag values into viper.
	viper.Set("profile", *profile)
	viper.Set("configServerUrl", *configServerUrl)
	viper.Set("configBranch", *configBranch)
}

func main() {
	fmt.Printf("Starting %v\n", appName)
	config.LoadConfigurationFromBranch(
		viper.GetString("configServerUrl"),
		appName,
		viper.GetString("profile"),
		viper.GetString("configBranch"))

	initializeBoltClient()
	go config.StartListener(appName, viper.GetString("amqp_server_url"), viper.GetString("config_event_bus"))
	fmt.Printf("port :   "+viper.GetString("server_port"))
	service.StartWebServer(viper.GetString("server_port"))
}

func initializeBoltClient() {
	service.DBClient = &dbclient.BoltClient{}
	service.DBClient.OpenBoltDb()
	service.DBClient.Seed()
}
