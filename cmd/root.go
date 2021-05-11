package cmd

import (
	"bytes"
	"encoding/json"
	"errors"
	"os"

	"github.com/karmanord/aerospike_viewer/aerospike_driver"
	"github.com/spf13/cobra"
)

var (
	hostFlag       string
	portFlag       int
	nameSpaceFlag  string
	setFlag        string
	keyFlag        string
	encodeTypeFlag string
	binFlag        bool
	listFlag       bool
)

type encodeType string

const (
	MessagePack encodeType = "msgpack"
)

func (e encodeType) String() string {
	return string(e)
}

func NewCmdRoot() *cobra.Command {
	rootCmd := &cobra.Command{
		Short:         "A cli that gets and displays the result when you specify the key",
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			if !binFlag && !listFlag {
				return errors.New("Specify --bin or -l")
			}

			conn, err := aerospike_driver.NewConnection(hostFlag, portFlag, nameSpaceFlag)
			if err != nil {
				return err
			}

			r, err := conn.Get(nameSpaceFlag, setFlag, keyFlag)
			if err != nil {
				return err
			}

			if binFlag {
				var jsonStr []byte
				if encodeTypeFlag == MessagePack.String() {
					if r.Bins["json"] == nil {
						return errors.New("The Bin of the target record is not json")
					}

					decodeMap := aerospike_driver.MessagePackDecode(r.Bins["json"].([]byte))
					bin := map[string]interface{}{
						"json": decodeMap,
					}
					if jsonStr, err = json.Marshal(bin); err != nil {
						return err
					}
				} else {
					if jsonStr, err = json.Marshal(r.Bins); err != nil {
						return err
					}
				}

				var buf bytes.Buffer
				if err := json.Indent(&buf, []byte(jsonStr), "", "  "); err != nil {
					return err
				}
				cmd.Println(buf.String())
			} else if listFlag {
				for name := range r.Bins {
					cmd.Println(name)
				}
			} else {
				// ここには来ない
			}

			return nil
		},
	}

	rootCmd.PersistentFlags().StringVar(&hostFlag, "host", "127.0.0.1", "Host")
	rootCmd.PersistentFlags().IntVar(&portFlag, "port", 3000, "Port")
	rootCmd.PersistentFlags().StringVar(&nameSpaceFlag, "ns", "", "Namespace")
	rootCmd.PersistentFlags().StringVar(&setFlag, "set", "", "Set")
	rootCmd.PersistentFlags().StringVar(&keyFlag, "key", "", "Key")
	rootCmd.PersistentFlags().StringVar(&encodeTypeFlag, "enc", "", "Encode Type [msgpack]")
	rootCmd.PersistentFlags().BoolVar(&binFlag, "bin", false, "Display the value of bin")
	rootCmd.PersistentFlags().BoolVarP(&listFlag, "list", "l", false, "Show only bin name")

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
