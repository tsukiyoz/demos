/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/fatih/color"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tsukiyoz/demos/design/chatroom/internal/server"
)

var (
	cfgFile string
	addr    = ":2022"
	banner  = `           __             __                                   
  _____   / /_   ____ _  / /_   _____  ____   ____    ____ ___ 
 / ___/  / __ \ / __  / / __/  / ___/ / __ \ / __ \  / __  __ \
/ /__   / / / // /_/ / / /_   / /    / /_/ // /_/ / / / / / / /
\___/  /_/ /_/ \__,_/  \__/  /_/     \____/ \____/ /_/ /_/ /_/ 
ChatRoom start on %s
                                                              `
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "chatroom",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		color.Yellow(banner+"\n", addr)
		initConfig()
		srv := server.NewServer()
		srv.RegisterHandle()
		log.Fatal(http.ListenAndServe(addr, nil))
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.chatroom.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func initConfig() {
	if len(cfgFile) == 0 {
		panic("config file not found")
	}
	cfgFile = filepath.Clean(cfgFile)
	viper.SetConfigFile(cfgFile)

	viper.AddConfigPath(cfgFile)
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}

	// cfg.SensitiveWords = viper.GetStringSlice("sensitive")
	// cfg.MessageQueueLen = viper.GetInt64("message-queue")

	viper.WatchConfig()
	viper.OnConfigChange(func(in fsnotify.Event) {
		viper.ReadInConfig()
		// cfg.SensitiveWords = viper.GetStringSlice("sensitive")
	})
}
