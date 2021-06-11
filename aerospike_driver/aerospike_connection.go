package aerospike_driver

import (
	"strconv"
	"time"

	"github.com/aerospike/aerospike-client-go"
	as "github.com/aerospike/aerospike-client-go"
	"github.com/shamaton/msgpack"
)

type Connection struct {
	Client *as.Client
}

func NewConnection(host string, port int, ns string) (*Connection, error) {
	policy := aerospike.NewClientPolicy()
	policy.Timeout = 5 * time.Second

	client, err := as.NewClientWithPolicy(policy, host, port)
	if err != nil {
		return nil, err
	}

	conn := &Connection{
		Client: client,
	}

	return conn, nil
}

func (conn *Connection) Get(nameSpace, setName, key string) (*as.Record, error) {
	asKeys, err := convertKey(nameSpace, setName, key)
	if err != nil {
		return nil, err
	}

	var res *as.Record
	for _, asKey := range asKeys {
		res, err = conn.Client.Get(nil, asKey)
		if res != nil {
			return res, nil
		}
	}
	return nil, err
}

func convertKey(nameSpace, setName, key string) ([]*as.Key, error) {
	askeys := make([]*as.Key, 0, 3)
	if asKey, err := as.NewKey(nameSpace, setName, key); err != nil {
		return nil, err
	} else {
		askeys = append(askeys, asKey)
	}

	intKey, err := strconv.Atoi(key)
	if err == nil {
		if asKey, err := as.NewKey(nameSpace, setName, intKey); err != nil {
			return nil, err
		} else {
			askeys = append(askeys, asKey)
		}
	}

	floatKey, err := strconv.ParseFloat(key, 64)
	if err == nil {
		if asKey, err := as.NewKey(nameSpace, setName, floatKey); err != nil {
			return nil, err
		} else {
			askeys = append(askeys, asKey)
		}
	}

	return askeys, nil
}

func MessagePackDecode(data []byte) map[string]interface{} {
	var decodeMap map[string]interface{}
	msgpack.Unmarshal(data, &decodeMap)

	for _, v := range decodeMap {
		recursiveToJSON(v)
	}

	return decodeMap
}

// NOTE: 2階層以上あるjsonをmsgpack.Unmarshal()すると、map[interface{}]interface{}となりjson.Marshal()出来なくなる為、以下ロジックで防ぐ
type (
	jsonArray []interface{}
	jsonMap   map[string]interface{}
)

func recursiveToJSON(v interface{}) interface{} {
	var r interface{}

	switch v := v.(type) {
	case []interface{}:
		for i, e := range v {
			v[i] = recursiveToJSON(e)
		}
		r = jsonArray(v)
	case map[interface{}]interface{}:
		newMap := make(map[string]interface{}, len(v))
		for k, e := range v {
			newMap[k.(string)] = recursiveToJSON(e)
		}
		r = jsonMap(newMap)
	default:
		r = v
	}

	return r
}
