package aerospike_driver

import (
	"time"

	"github.com/aerospike/aerospike-client-go"
	as "github.com/aerospike/aerospike-client-go"
)

type AerospikeConnection struct {
	CreatedAt time.Time
	Client    *as.Client
	Namespace string
}

func NewSession(host string, port int, ns string) (*AerospikeConnection, error) {
	policy := aerospike.NewClientPolicy()
	policy.Timeout = 5 * time.Second

	client, err := as.NewClientWithPolicy(policy, host, port)
	if err != nil {
		return nil, err
	}

	sess := &AerospikeConnection{
		CreatedAt: time.Now(),
		Client:    client,
		Namespace: ns,
	}
	return sess, nil
}

// func TempNewSession(host string, port int, ns string) (*AerospikeConnection, error) {
// 	policy := aerospike.NewClientPolicy()
// 	policy.Timeout = 5 * time.Second

// 	client, err := as.NewClientWithPolicy(policy, host, port)
// 	if err != nil {
// 		return nil, err
// 	}

// 	// init policy
// 	defaultWritePolicy = as.NewWritePolicy(0, 0)
// 	defaultWritePolicy.SendKey = true
// 	defaultListPolicy = as.NewListPolicy(as.ListOrderUnordered, as.ListWriteFlagsAddUnique|as.ListWriteFlagsNoFail|as.ListWriteFlagsPartial)
// 	defaultScanPolicy = as.NewScanPolicy()

// 	sess := &AerospikeConnection{
// 		CreatedAt: time.Now(),
// 		Client:    client,
// 		Namespace: ns,
// 	}
// 	return sess, nil
// }

func GetConnection(host string, port int, ns string) (*AerospikeConnection, error) {
	session, err := NewSession(host, port, ns)
	// session, err := TempNewSession(host, port, ns)
	if err != nil {
		return nil, err
	}

	return session, nil
}

func (conn *AerospikeConnection) Get(nameSpace, setName, key string) (*as.Record, error) {
	asKey, err := as.NewKey(nameSpace, setName, key)
	if err != nil {
		return nil, err
	}
	// policy := as.NewPolicy()
	// policy.SendKey = true
	return conn.Client.Get(nil, asKey)
}

// func (conn *AerospikeConnection) Put(nameSpace, setName, key string) error {
// 	// spew.Dump(nameSpace, setName, key)_
// 	k, err := as.NewKey(nameSpace, setName, key)
// 	if err != nil {
// 		return err
// 	}
// 	// policy := as.NewPolicy()
// 	// policy.SendKey = true
// 	// spew.Dump(policy)
// 	bin := &as.BinMap{
// 		"ad":      "a1",
// 		"content": "c1",
// 	}
// 	return conn.Client.Put(defaultWritePolicy, k, *bin)
// }
