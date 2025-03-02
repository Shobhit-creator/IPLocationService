package main

import (
	"fmt"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
	"github.com/shobhit-Creator/IPLocationService/internal/service"
	"github.com/spf13/viper"
)

func initConfig() {
	viper.SetConfigName("config")  
	viper.SetConfigType("yaml")     
	viper.AddConfigPath("../../") 

	path, err := filepath.Abs("../../config")
	if err != nil {
		panic("config dire not found")
	}
	viper.AddConfigPath(path)
	viper.WatchConfig()
	err = viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("can not read config: %w", err))
	}
	viper.OnConfigChange(func (e fsnotify.Event) {
		fmt.Println("Config file changed:", e.Name)
		service.LoadProvidersFromConfig()
	})
}