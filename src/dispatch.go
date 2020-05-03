package dispatch

import (
	"fmt"

	"github.com/spf13/cobra"
)

func Dispatch(cmd *cobra.Command, args []string) {
	fmt.Print("hello world!")
}
