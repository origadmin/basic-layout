package utils

import (
	"strings"

	"github.com/spf13/cobra"
)

func ToLower(cmd *cobra.Command) string {
	return strings.ToLower(cmd.Root().Name())
}
