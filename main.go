package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"reflect"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"

	"alvinr.ca/learn-go/debug-api/models"
)

var (
	// LogInfo to use when logging information
	LogInfo *log.Logger
	// LogWarn to use when logging warnings which impact API performance
	LogWarn *log.Logger
	// LogError to use when the situation full on breaks the API
	LogError *log.Logger
)

func main() {

	fetchConfig()

	initLogging()

}

// initLogging prepares internal logging engine
// Logs to both engine and
func initLogging() {

	file, err := os.OpenFile(viper.GetString("LogPath"), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	multiWriter := io.MultiWriter(os.Stdout, file)
	if err != nil {
		log.Fatal(err)
	}

	LogInfo = log.New(file, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	LogWarn = log.New(file, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
	LogError = log.New(file, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)

	LogInfo.SetOutput(multiWriter)
	LogWarn.SetOutput(multiWriter)
	LogError.SetOutput(multiWriter)
}

// fetchConfig sets up viper for config management
func fetchConfig() {

	viper.SetConfigName("api-configs.yaml")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()

	if err != nil {

		log.Println("Config file doesn't exist, creating ./api-configs.yaml")

		newConf := models.NewConfig()
		newConfVals := reflect.ValueOf(newConf)
		newConfType := newConfVals.Type()

		for i := 0; i < newConfVals.NumField(); i++ {
			viper.SetDefault(newConfType.Field(i).Name, newConfVals.Field(i).Interface())
		}

		viper.WriteConfigAs("./api-configs.yaml")
	}

	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		restart()
	})
}

// restart the API server
// Used for hot config
func restart() {
	fmt.Println("Config change detected, restarting")
}
