package ir

import (
	"bytes"
	"fmt"
	"testing"
)

func TestPeer(t *testing.T) {

	peer1 := "QmS3zcG7LhYZYSJMhyRZvTddvbNUqtt8BJpaSs6mi1K5Va"
	peer2 := "QmS3zcG7LhYZYSJMhyRZvTddvbNUqtt8BJpaSs634rfsa3"

	n := Peer{
		Dict{
			Tag: "foo",
			Pairs: Pairs{
				{String{peer1}, String{peer1}},
				{String{peer2}, String{peer2}},
				{String{"bar2"}, Int64{13}},
			},
		},
	}
	var w bytes.Buffer
	n.WritePretty(&w)
	fmt.Println(" ===== Sample Peer ======")
	fmt.Println(w.String())
}
