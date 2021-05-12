package cmd

import (
	"bytes"
	"encoding/json"
	"errors"
	"os"
	"reflect"
	"sort"

	"github.com/aerospike/aerospike-client-go"
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
	listKeysFlag   bool
	listBinsFlag   bool
)

type encodeType string

const (
	messagePack encodeType = "msgpack"
)

func (e encodeType) String() string {
	return string(e)
}

func NewCmdRoot() *cobra.Command {
	rootCmd := &cobra.Command{
		Short:         "A cli that gets and displays the result when you specify the key",
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			if !binFlag && !listBinsFlag && !listKeysFlag {
				return errors.New("Specify --bin or --list-keys or --list-bins ")
			}

			conn, err := aerospike_driver.NewConnection(hostFlag, portFlag, nameSpaceFlag)
			if err != nil {
				return err
			}

			switch {
			case binFlag:
				r, err := conn.Get(nameSpaceFlag, setFlag, keyFlag)
				if err != nil {
					return err
				}

				var jsonStr []byte
				if encodeTypeFlag == messagePack.String() {
					bins := make(map[string]interface{})
					for k, v := range r.Bins {
						if reflect.TypeOf(v).String() == "[]uint8" {
							bins[k] = aerospike_driver.MessagePackDecode(v.([]byte))
						} else {
							bins[k] = v
						}
					}

					if jsonStr, err = json.Marshal(bins); err != nil {
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

			case listKeysFlag:
				defaultScanPolicy := aerospike.NewScanPolicy()
				defaultScanPolicy.MultiPolicy.IncludeBinData = false
				recordSets, _ := conn.Client.ScanAll(defaultScanPolicy, nameSpaceFlag, setFlag)

				var keys []string
				for v := range recordSets.Results() {
					keys = append(keys, v.Record.Key.Value().String())
				}

				sortedCmdPrintln(cmd, keys)

			case listBinsFlag:
				r, err := conn.Get(nameSpaceFlag, setFlag, keyFlag)
				if err != nil {
					return err
				}

				names := make([]string, 0, len(r.Bins))
				for name := range r.Bins {
					names = append(names, name)
				}

				sortedCmdPrintln(cmd, names)

			default:
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
	rootCmd.PersistentFlags().BoolVar(&listBinsFlag, "list-bins", false, "Display only bin name")
	rootCmd.PersistentFlags().BoolVar(&listKeysFlag, "list-keys", false, "Display Key List")

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

func sortedCmdPrintln(cmd *cobra.Command, sliceString []string) {
	sort.Strings(sliceString)
	for _, s := range sliceString {
		cmd.Println(s)
	}
}
