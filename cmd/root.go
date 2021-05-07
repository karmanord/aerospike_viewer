package cmd

import (
	"bytes"
	"encoding/json"
	"os"

	"github.com/karmanord/aerospike_viewer/aerospike_driver"
	"github.com/spf13/cobra"
)

var (
	host       string
	port       int
	nameSpace  string
	set        string
	key        string
	encodeType string
)

func NewCmdRoot() *cobra.Command {
	rootCmd := &cobra.Command{
		// Use:           "aerospike_viewer",
		// Use:           "aerospike_viewer",
		Short:         "A cli that gets and displays the result when you specify the key",
		SilenceErrors: true,
		// SilenceUsage:  true,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := parseFlag(cmd); err != nil {
				return err
			}
			conn, err := aerospike_driver.GetConnection(host, port, nameSpace)
			if err != nil {
				return err
			}
			// if err := conn.Put(nameSpace, set, key); err != nil {
			// 	cmd.PrintErrln(err)
			// }

			r, err := conn.Get(nameSpace, set, key)
			if err != nil {
				return err
			}

			// cmd.Println(conn)
			// cmd.Println(r)
			jsonStr, err := json.Marshal(r.Bins)
			if err != nil {
				return err
			}

			var buf bytes.Buffer
			if err := json.Indent(&buf, []byte(jsonStr), "", "  "); err != nil {
				return err
			}

			cmd.Println(buf.String())

			return nil
		},
	}

	rootCmd.PersistentFlags().StringVar(&host, "host", "127.0.0.1", "Host")
	rootCmd.PersistentFlags().IntVar(&port, "port", 3000, "Port")
	rootCmd.PersistentFlags().StringVar(&nameSpace, "ns", "", "Namespace")
	rootCmd.PersistentFlags().StringVar(&set, "set", "", "Set")
	rootCmd.PersistentFlags().StringVar(&key, "key", "", "Key")
	rootCmd.PersistentFlags().StringVar(&encodeType, "enc", "", "Encode Type [msgpack]")

	return rootCmd
}

func Execute() {
	cmd := NewCmdRoot()
	cmd.SetOutput(os.Stdout)
	if err := cmd.Execute(); err != nil {
		cmd.SetOutput(os.Stderr)
		cmd.PrintErrf("Error: %v", err.Error())
	}

}

func parseFlag(cmd *cobra.Command) error {
	var err error
	if host, err = cmd.PersistentFlags().GetString("host"); err != nil {
		return err
	}
	if port, err = cmd.PersistentFlags().GetInt("port"); err != nil {
		return err
	}
	if nameSpace, err = cmd.PersistentFlags().GetString("ns"); err != nil {
		return err
	}
	if set, err = cmd.PersistentFlags().GetString("set"); err != nil {
		return err
	}
	if key, err = cmd.PersistentFlags().GetString("key"); err != nil {
		return err
	}
	if encodeType, err = cmd.PersistentFlags().GetString("enc"); err != nil {
		return err
	}
	return nil
}
