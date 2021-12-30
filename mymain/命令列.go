package mymain

import (
	"crypto/des"

	"github.com/spf13/cobra"
)

var (
	rootCommand = &cobra.Command{
		SilenceUsage: true,
		Args:         checkArguments(cobra.ExactArgs(1)),
		Run: func(cmd *cobra.Command, args []string) {
			desKey = args[0]
		},
	}
)

// 檢查命令列參數
func checkArguments(checks ...cobra.PositionalArgs) cobra.PositionalArgs {

	return func(cmd *cobra.Command, args []string) error {

		for _, check := range checks {

			if err := check(cmd, args); err != nil {
				return err
			}

		}

		if _, err := des.NewCipher([]byte(args[0])); err != nil {
			return err
		} else {
			return nil
		}

	}
}
