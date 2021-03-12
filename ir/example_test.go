package ir

import (
	"bytes"
	"fmt"
	"testing"
)

func TestPretty(t *testing.T) {
	peer1 := "QmS3zcG7LhYZYSJMhyRZvTddvbNUqtt8BJpaSs6mi1K5Va"
	peer2 := "QmS3zcG7LhYZYSJMhyRZvTddvbNUqtt8BJpaSs634rfsa3"
	cid1 := "bafzbeigai3eoy2ccc7ybwjfz5r3rdxqrinwi4rwytly24tdbh6yk7zslrm"

	n := Dict{
		Tag: "foo",
		Pairs: Pairs{
			{String{"cid"}, Cid{
				Dict{
					Tag: "cid",
					Pairs: Pairs{
						{String{cid1}, String{cid1}},
					},
				}}},
			{String{"bar2"}, Int64{13}},
			{String{"providers"}, Peer{
				Dict{
					Tag: "foo",
					Pairs: Pairs{
						{String{peer1}, String{peer1}},
						{String{peer2}, String{peer2}},
						{String{"bar2"}, Int64{13}},
					},
				},
			}},
		},
	}
	var w bytes.Buffer
	n.WritePretty(&w)
	fmt.Println(" ===== Sample Dictionary ======")
	fmt.Println(w.String())
}
