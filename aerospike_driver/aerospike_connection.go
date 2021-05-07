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

var (
	defaultScanPolicy *as.ScanPolicy

	defaultWritePolicy *as.WritePolicy
	defaultListPolicy  *as.ListPolicy
)

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

func GetConnection(host string, port int, ns string) (*AerospikeConnection, error) {
	session, err := NewSession(host, port, ns)
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

	return conn.Client.Get(nil, asKey)
}
