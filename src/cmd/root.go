package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"repo.abanicon.com/abantheter-microservices/websocket/app/configs"
	"repo.abanicon.com/public-library/glogger"
	"strings"
)

var (
	configFile string

	rootCmd = &cobra.Command{
		Use:   "websocket",
		Short: "short description",
		Long:  `long description`,
	}
)

func Execute() {
	rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().StringP("author", "a", "Mahdi Jedari", "i.jedari@gmail.com")
	cobra.OnInitialize(initConfig)
}

func initConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("../config")

	viper.AutomaticEnv()
	viper.SetEnvPrefix("websocket")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := viper.ReadInConfig(); err != nil {
		glogger.Fatalf("Fatal error config file: %s \n", err)
	}
	viper.Unmarshal(&configs.Config)
	glogger.Println("configuration initialized! (Notice: configurations may be initialised from OS ENV)")
}
