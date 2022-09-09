/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

const (
	bash_completion_func = `__upx_parse_get()
{
    local upx_output out
    if upx_output=$(upx user ls --name 2>/dev/null); then
        out=($(echo "${upx_output}"))
        COMPREPLY=( $( compgen -W "${out[*]}" -- "$cur" ) )
    fi
}

__upx_get_resource()
{
    __upx_parse_get
    if [[ $? -eq 0 ]]; then
        return 0
    fi
}

__custom_func() {
    case ${last_command} in
        upx_user_cu)
            __upx_get_resource
            return
            ;;
        *)
            ;;
    esac
}
`
)

var (
	cfgFile     = ""
	cfgColorful = false

	rootCmd = &cobra.Command{
		Use:                    "upx", // 如果这里设置成了app，那么completion生成的内容就是_app_root_command
		Short:                  "A brief description of your application",
		BashCompletionFunction: bash_completion_func,
	}
)

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.app.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.app.yaml)")
	rootCmd.PersistentFlags().BoolVarP(&cfgColorful, "colorful", "", false, "console colorful mode")

	// autocomplete set
	// rootCmd.CompletionOptions.DisableDefaultCmd = true
	// rootCmd.CompletionOptions.DisableNoDescFlag = true
	// rootCmd.CompletionOptions.DisableDescriptions = true
}
