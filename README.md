# aerospike_viewer
A cli that gets and displays the result when you specify the key

## Installation
`$ go get -u https://github.com/karmanord/aerospike_viewer`

`$ go install`
## Usage

`ex. $ aerospike_viewer --host 127.0.0.1 --ns test --set set1 --key 123`

### Output result
```
{
  "bin1": "b1",
  "bin2": "b2"
}
```

### Flags
```
    --enc string    Encode Type [msgpack]
-h, --help          help for this command
    --host string   Host (default "127.0.0.1")
    --key string    Key
    --ns string     Namespace
    --port int      Port (default 3000)
    --set string    Set
```