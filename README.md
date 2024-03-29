# aerospike_viewer
A cli that gets and displays the result when you specify the key

## Installation
`$ go get -u github.com/karmanord/aerospike_viewer`

## Usage

`ex. $ aerospike_viewer --host 127.0.0.1 --ns test --set set1 --key 111 --bin`

### Output result(Bin Name And Value)
```
{
  "bin1": "b1",
  "bin2": "b2"
}
```

`ex. $ aerospike_viewer --host 127.0.0.1 --ns test --set set1 --key 111 --list-bins`

### Output result(Bin Name Only)
```
bin1
bin2
```

`ex. $ aerospike_viewer --host 127.0.0.1 --ns test --set set1 --list-keys`

### Output result(Key List)
```
key1
key2
key3
```

`ex. $ aerospike_viewer --host 127.0.0.1 --ns test --set set1 --key 222 --bin`

### Output result(Do not decode MessagePack)
```
{
  "json": "g6VUYWdJZM4AAbIHq0FkQWNjb3VudElkzgADZA6pUmVsYXRpb25zkoSiSWTOAAUWFaNVcmykdXJsMalNYANkDg="
}
```

`ex. $ aerospike_viewer --host 127.0.0.1 --ns test --set set1 --key 222 --bin --enc msgpack`

### Output result(Decode MessagePack)
```
{
  "json": {
    "AId": 111111,
    "BId": 222222
    "Internal": [
      {
        "Name": name1,
        "Url": "url1"
      },
      {
        "Name": name2,
        "Url": "url2"
      }
    ],
  }
}
```

### Flags
```
    --bin           Display the value of bin
    --enc string    Encode Type [msgpack]
-h, --help          help for this command
    --host string   Host (default "127.0.0.1")
    --key string    Key
    --list-bins     Display only bin name
    --list-keys     Display Key List
    --ns string     Namespace
    --port int      Port (default 3000)
    --set string    Set
```