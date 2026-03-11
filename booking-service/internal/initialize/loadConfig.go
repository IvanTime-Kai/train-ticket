package initialize

import (
	"fmt"
	"os"

	"github.com/leminhthai/train-ticket/booking-service/global"
	"github.com/spf13/viper"
)

func LoadConfig() {
	viper := viper.New()

	viper.AddConfigPath("./configs")

	configName := os.Getenv("CONFIG_NAME")

	if configName == "" {
		configName = "local"
	}

	viper.SetConfigName(configName)
	viper.SetConfigType("yaml")

	err := viper.ReadInConfig()

	if err != nil {
		panic(fmt.Errorf("failed to read configuration: %w", err))
	}

	if err := viper.Unmarshal(&global.Config); err != nil {
		fmt.Printf("Unable to decode configuration: %v", err)
	}
}
