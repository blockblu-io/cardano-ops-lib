package topology

import (
	"fmt"
)

func fillMergeMap(top *Topology, m *map[string][]ConnectionConfig) {
	if top == nil {
		return
	}
	for _, conConfig := range top.Producers {
		id := fmt.Sprintf("%s:%d", conConfig.HostAddress, conConfig.Port)
		if val, exists := (*m)[id]; exists {
			(*m)[id] = append(val, conConfig)
		} else {
			(*m)[id] = []ConnectionConfig{conConfig}
		}
	}
}

func mapToString(list []ConnectionConfig, f func(ConnectionConfig) *string) []*string {
	stringList := make([]*string, len(list))
	for i, val := range list {
		stringList[i] = f(val)
	}
	return stringList
}

func mapToInt(list []ConnectionConfig, f func(ConnectionConfig) *int) []*int {
	intList := make([]*int, len(list))
	for i, val := range list {
		intList[i] = f(val)
	}
	return intList
}

func selectInt(list []*int) *int {
	for _, elem := range list {
		if elem != nil {
			return elem
		}
	}
	return nil
}

func selectString(list []*string) *string {
	for _, elem := range list {
		if elem != nil {
			return elem
		}
	}
	return nil
}
