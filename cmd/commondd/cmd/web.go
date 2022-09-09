/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"exmaple/cobra_demo/app"
	"github.com/spf13/cobra"
)

var (
	config_path  string
	casbin_model string
	debug        string
)

// webCmd represents the web command
var webCmd = &cobra.Command{
	Use:   "web",
	Short: "run project",
	Run: func(cmd *cobra.Command, args []string) {
		app.Run(config_path, casbin_model, debug)
	},
}

func init() {
	rootCmd.AddCommand(webCmd)

	// Specifying a configuration file
	webCmd.Flags().StringVarP(&config_path, "config_path", "f", "", "指定配置文件")

	// Specifying a casbin file
	webCmd.Flags().StringVarP(&casbin_model, "casbin_model", "c", "", "指定casbin文件")

	// Whether open debug mode
	webCmd.Flags().StringVarP(&debug, "debug", "d", "", "是否以debug的模式运行")

	// require
	webCmd.MarkFlagRequired("config_path")
}
