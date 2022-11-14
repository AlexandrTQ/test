package cmd

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	path    string
	rootCmd = &cobra.Command{
		Use: "TransactionServer",
	}
)

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&path, "config", "", "config's path")
}

func initConfig() {
	if path == "" {
		log.Fatal("empty path to config file")
	}
	fileDir := filepath.Dir(path)
	file := filepath.Base(path)
	fileName, fileExec := strings.Split(file, ".")[0], strings.Split(file, ".")[1]
	if err := InitConfig(fileDir, fileName, fileExec); err != nil {
		log.Fatal(err)
	}
}

func InitConfig(fileDir string, fileName string, fileExec string) error {
	viper.AddConfigPath(fileDir)
	viper.SetConfigName(fileName)
	viper.SetConfigType(fileExec)

	err := viper.ReadInConfig()
	if err != nil {
		return err
	}

	return nil
}
