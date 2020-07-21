package gomvc

import "github.com/spf13/cobra"

var root = &cobra.Command{
	Use:   "gomvc",
	Short: "GoMVC is a CLI for generating and modifying golang MVC applications",
	Run: func(cmd *cobra.Command, args []string) {
		if err := cmd.Help(); err != nil {
			panic(err)
		}
	},
}

// Root is the function called when running `gomvc`
func Root() *cobra.Command {
	return root
}
