package topology

import (
	"reflect"
	"testing"
)

var (
	val1             = 1
	val2             = 2
	friendlyNameLin  = "SOBIT Lin"
	operatorNameLove = "LOVE"
	nodeNameLove     = "LOVE US"
	lin1             = ConnectionConfig{
		FriendlyName: &friendlyNameLin,
		HostAddress:  "114.14.14.1",
		Port:         1414,
		Valency:      &val2,
	}
	lin2 = ConnectionConfig{
		Operator:    &friendlyNameLin,
		HostAddress: "114.14.14.1",
		Port:        1414,
		Valency:     &val1,
	}
	love = ConnectionConfig{
		NodeName:    &nodeNameLove,
		Operator:    &operatorNameLove,
		HostAddress: "112.12.12.1",
		Port:        1313,
		Valency:     &val2,
	}
	aa = Topology{
		Producers: []ConnectionConfig{
			love,
			lin1,
		},
	}
	bb = Topology{
		Producers: []ConnectionConfig{
			love,
			lin2,
		},
	}
)

func TestSameAs_TwoStrictlyEqualConfigs_ReturnTrue(t *testing.T) {
	r1 := lin1.SameAs(lin2)
	r2 := lin2.SameAs(lin1)
	if !r1 || !r2 {
		t.Errorf("two strictly equal connection configs are the same, but returned false.")
	}
}

func TestSameAs_TwoSameConfigs_ReturnTrue(t *testing.T) {
	r1 := lin1.SameAs(lin2)
	r2 := lin2.SameAs(lin1)
	if !r1 || !r2 {
		t.Errorf("two connection configs with equal host address and port are the same, but returned false.")
	}
}

func TestSameAs_TwoDifferentConfigs_ReturnFalse(t *testing.T) {
	r1 := lin1.SameAs(love)
	r2 := love.SameAs(lin1)
	if r1 || r2 {
		t.Errorf("the two connection configs are different, but returned true.")
	}
}

func TestMerge_Merge2NilTopology_ReturnNil(t *testing.T) {
	c := Merge(nil, nil)
	if !reflect.ValueOf(c).IsNil() {
		t.Errorf("merge of two nil topologies must return nil.")
	}
}

// todo: add more test cases
