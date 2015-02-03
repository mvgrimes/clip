package commands

import (
	// "fmt"
	"github.com/spf13/cobra" // cli
	"github.com/spf13/viper" // Config file parsing
)

func InitializeConfig() {
	viper.SetConfigName("config")

	viper.AddConfigPath("/etc/clip")
	viper.AddConfigPath("$HOME/.config/clip")

	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	viper.SetDefault("amqp", "amqp://guest:guest@localhost:5672/")
	viper.SetDefault("exchange", "clip")
	viper.SetDefault("key", "a very very very very secret key") // 32 bytes
	viper.SetDefault("verbose", false)
	// viper.SetDefault("log", )

}

var rootCmd = &cobra.Command{
	Use:   "clipx",
	Short: "Copy and paste between remote systems",
	// Long ``,
	Run: func(cmd *cobra.Command, args []string) {
		InitializeConfig()
	},
}

func Execute() {
	rootCmd.AddCommand(copyCmd)
	rootCmd.AddCommand(pasteCmd)
	rootCmd.Execute()
}
