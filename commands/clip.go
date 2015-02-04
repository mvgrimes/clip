package commands

import (
	// "fmt"
	"github.com/spf13/cobra" // cli
	"github.com/spf13/viper" // Config file parsing
)

var ClipCmd = &cobra.Command{
	Use:   "clip",
	Short: "Copy and paste between remote systems",
	Long: `Copy and paste between remote systems
	Using RabbitMQ. 
	`,
	// Run: func(cmd *cobra.Command, args []string) {
	// 	InitializeConfig()
	// },
}

var clipCmdV *cobra.Command
var Verbose bool
var Server, Exchange, Key, PasteCmd, CopyCmd string

func Execute() {
	ClipCmd.AddCommand(copyCmd)
	ClipCmd.AddCommand(pasteCmd)
	ClipCmd.Execute()
}

func init() {
	ClipCmd.PersistentFlags().BoolVarP(&Verbose, "verbose", "v", false, "Let us know what is going on")
	ClipCmd.PersistentFlags().StringVarP(&Server, "server", "s", "", "RabbitMQ server connect string")
	ClipCmd.PersistentFlags().StringVarP(&Exchange, "exchange", "e", "", "RabbitMQ exchange")
	ClipCmd.PersistentFlags().StringVarP(&Key, "key", "k", "", "Private encryption key")
	ClipCmd.PersistentFlags().StringVarP(&PasteCmd, "paste-cmd", "p", "", "Path to pbpaste")
	clipCmdV = ClipCmd
}

func InitializeConfig() {
	viper.SetConfigName("config")

	viper.AddConfigPath("/etc/clip")
	viper.AddConfigPath("$HOME/.config/clip")

	err := viper.ReadInConfig()
	failOnError(err, "Error reading config file")

	viper.SetDefault("Verbose", false)
	viper.SetDefault("Server", "amqp://guest:guest@localhost:5672/")
	viper.SetDefault("Exchange", "clip")
	viper.SetDefault("Key", "a very very very very secret key") // 32 bytes
	viper.SetDefault("PasteCmd", "pbpaste")

	if clipCmdV.PersistentFlags().Lookup("verbose").Changed {
		viper.Set("Verbose", Verbose)
	}
	if clipCmdV.PersistentFlags().Lookup("server").Changed {
		viper.Set("Server", Server)
	}
	if clipCmdV.PersistentFlags().Lookup("exchange").Changed {
		viper.Set("Exchange", Exchange)
	}
	if clipCmdV.PersistentFlags().Lookup("key").Changed {
		viper.Set("Key", Key)
	}
	if clipCmdV.PersistentFlags().Lookup("paste-cmd").Changed {
		viper.Set("PasteCmd", PasteCmd)
	}
}
