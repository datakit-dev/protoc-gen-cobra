package cmd

import (
	"github.com/datakit-dev/protoc-gen-cobra/example/pb"

	_ "github.com/datakit-dev/protoc-gen-cobra/auth/jwt"
	_ "github.com/datakit-dev/protoc-gen-cobra/auth/oauth"
	_ "github.com/datakit-dev/protoc-gen-cobra/iocodec/yaml"
)

func init() {
	rootCmd.AddCommand(pb.BankClientCommand())
	rootCmd.AddCommand(pb.CRUDClientCommand())
	rootCmd.AddCommand(pb.CacheClientCommand())
	rootCmd.AddCommand(pb.CyclicalClientCommand())
	rootCmd.AddCommand(pb.DeprecatedClientCommand())
	rootCmd.AddCommand(pb.NestedClientCommand())
	rootCmd.AddCommand(pb.OneofClientCommand())
	rootCmd.AddCommand(pb.Proto2ClientCommand())
	rootCmd.AddCommand(pb.TimerClientCommand())
	rootCmd.AddCommand(pb.TypesClientCommand())
}
