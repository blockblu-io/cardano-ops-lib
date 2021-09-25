package topology

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
)

// Topology is a list of cardano-nodes expected by a cardano-node at startup. It
// specifies to which nodes a connection attempt shall be started.
type Topology struct {
	// Producers is a list of ConnectionConfig to cardano-nodes
	Producers []ConnectionConfig `json:"Producers"`
}

// ConnectionConfig specifies the details needed by a cardano-node to connect to
// another cardano-node plus optional metadata.
type ConnectionConfig struct {
	// HostAddress specifies the DNS name or IP address of the node(s).
	HostAddress string `json:"addr"`
	// Port specifies the TCP port of the node(s) to connect to.
	Port int `json:"port"`
	// Valency specifies the number of nodes that can be resolved from the given
	// DNS name. Valency should be 1 for IP addresses.
	Valency *int `json:"valency,omitempty"`
	// NodeName is a metadata field for this config, giving it a name.
	NodeName *string `json:"node,omitempty"`
	// FriendlyName is a metadata field for this config, giving it a friendly name.
	FriendlyName *string `json:"friendly_name,omitempty"`
	// Operator is a metadata field for this config, specifying the operator.
	Operator *string `json:"operator,omitempty"`
}

// SameAs checks whether this ConnectionConfig is referring to the same node as the
// given ConnectionConfig. Two ConnectionConfig are the same, if the host address
// and the port are strictly equal. However, the other fields might not match.
func (top *ConnectionConfig) SameAs(other ConnectionConfig) bool {
	topLocation := fmt.Sprintf("%s:%d", top.HostAddress, top.Port)
	otherLocation := fmt.Sprintf("%s:%d", other.HostAddress, other.Port)
	return topLocation == otherLocation
}

// Merge merges the ConnectionConfig of producers listed in the given Topology objects.
//
// If all of the given Topology objects are nil, then nil will be returned. Is only one
// of the given Topology objects isn't nil, then this Topology object is returned.
// Otherwise a merge of all non-nil Topology objects is performed.
//
// In a merge, the first given value for a Topology field is prioritized, and all the
// same fields in the following Topology objects are ignored. If a field only
// occurs in one Topology, then it will occur equally in the merged topology.
func Merge(a *Topology, list ...*Topology) *Topology {
	if a == nil {
		all := true
		for _, top := range list {
			if top != nil {
				all = false
				break
			}
		}
		if all {
			return nil
		}
	}
	topMap := map[string][]ConnectionConfig{}
	fillMergeMap(a, &topMap)
	for _, top := range list {
		fillMergeMap(top, &topMap)
	}
	n := 0
	top := Topology{}
	producers := make([]ConnectionConfig, len(topMap))
	for _, ccList := range topMap {
		conConfig := ConnectionConfig{}
		conConfig.HostAddress = ccList[0].HostAddress
		conConfig.Port = ccList[0].Port
		conConfig.Valency = selectInt(mapToInt(ccList, func(config ConnectionConfig) *int {
			return config.Valency
		}))
		conConfig.FriendlyName = selectString(mapToString(ccList, func(config ConnectionConfig) *string {
			return config.FriendlyName
		}))
		conConfig.NodeName = selectString(mapToString(ccList, func(config ConnectionConfig) *string {
			return config.NodeName
		}))
		conConfig.Operator = selectString(mapToString(ccList, func(config ConnectionConfig) *string {
			return config.Operator
		}))
		producers[n] = conConfig
		n++
	}
	top.Producers = producers
	return &top
}

// ReadTopology reads the data from the given io.Reader, which should be a Topology
// serialized to JSON. If the Topology data can be successfully read and deserialized,
// then the corresponding Topology object is returned. Otherwise, nil and an error will
// be returned.
func ReadTopology(topologyReader io.Reader) (*Topology, error) {
	var topology Topology
	var data, err = ioutil.ReadAll(topologyReader)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(data, &topology)
	if err != nil {
		return nil, err
	}
	return &topology, nil
}

// WriteTopology writes the given topology to the given io.Writer. If the write
// fails, then an error will be returned. Otherwise nil will be returned.
func WriteTopology(topology *Topology, topologyWriter io.Writer) error {
	data, err := json.Marshal(topology)
	if err != nil {
		return err
	}
	_, err = topologyWriter.Write(data)
	return err
}
