package ir

import (
	"testing"
	"time"

	"github.com/libp2p/go-smart-record/xr"
)

const sampleTTL = 123

func TestSetMetadata(t *testing.T) {
	now := time.Now().Unix()

	d := xr.Dict{
		Tag: "aaa",
		Pairs: xr.Pairs{
			xr.Pair{Key: xr.String{Value: "x"}, Value: xr.NewInt64(1)},
			xr.Pair{Key: xr.String{Value: "w"}, Value: xr.NewInt64(1)},
		},
	}

	ttl := TTL(uint64(sampleTTL))
	ds, err := SyntacticGrammar.Assemble(AssemblerContext{Grammar: SyntacticGrammar}, d, []Metadata{ttl}...)
	if err != nil {
		t.Fatal(err)
	}
	m := ds.Metadata()
	if m.TTL != uint64(sampleTTL) {
		t.Fatal("TTL not set successfully in node:", m.TTL, sampleTTL)
	}
	if m.AssemblyTime < uint64(now) {
		t.Fatal("Assembly time set failed:", m.AssemblyTime, now)
	}

}
func TestUpdateMetadata(t *testing.T) {

	d1 := xr.Dict{
		Tag: "aaa",
		Pairs: xr.Pairs{
			xr.Pair{Key: xr.String{Value: "x"}, Value: xr.NewInt64(1)},
		},
	}
	d2 := xr.Dict{
		Tag: "aaa",
		Pairs: xr.Pairs{
			xr.Pair{Key: xr.String{Value: "w"}, Value: xr.NewInt64(1)},
		},
	}

	ttl1 := TTL(uint64(sampleTTL))
	ttlval2 := sampleTTL + 3
	ttl2 := TTL(uint64(ttlval2))
	// Assemble first
	ds1, err := SyntacticGrammar.Assemble(AssemblerContext{Grammar: SyntacticGrammar}, d1, []Metadata{ttl1}...)
	if err != nil {
		t.Fatal(err)
	}
	// Assemble second
	ds2, err := SyntacticGrammar.Assemble(AssemblerContext{Grammar: SyntacticGrammar}, d2, []Metadata{ttl2}...)
	if err != nil {
		t.Fatal(err)
	}

	// Update
	dsu, err := ds1.UpdateWith(DefaultUpdateContext{}, ds2)
	m := dsu.Metadata()
	if m.TTL != uint64(ttlval2) {
		t.Fatal("TTL not updated successfully in node:", m.TTL, ttlval2)
	}
	if m.AssemblyTime > ds1.Metadata().AssemblyTime {
		t.Fatal("Assembly time update failed:", m.AssemblyTime, ds1.Metadata().AssemblyTime)
	}

	// Update without ttl being set in updated node
	ds2nottl, err := SyntacticGrammar.Assemble(AssemblerContext{Grammar: SyntacticGrammar}, d2)
	if err != nil {
		t.Fatal(err)
	}
	dsu, err = ds1.UpdateWith(DefaultUpdateContext{}, ds2nottl)
	if err != nil {
		t.Fatal(err)
	}
	m = dsu.Metadata()
	if m.TTL != uint64(sampleTTL) {
		t.Fatal("Update with no ttl not updated successfully in node:", m.TTL, sampleTTL)
	}
	if m.AssemblyTime > ds1.Metadata().AssemblyTime {
		t.Fatal("Assembly time update failed:", m.AssemblyTime, ds1.Metadata().AssemblyTime)
	}
}

func TestSetMetadataUpdate(t *testing.T) {

	d1 := xr.Set{
		Tag: "aaa",
		Elements: xr.Nodes{
			xr.String{Value: "x"},
		},
	}
	d2 := xr.Set{
		Tag: "aaa",
		Elements: xr.Nodes{
			xr.String{Value: "w"},
		},
	}
	ttl1 := TTL(uint64(sampleTTL))
	ttlval2 := sampleTTL + 3
	ttl2 := TTL(uint64(ttlval2))
	// Assemble first
	ds1, err := SyntacticGrammar.Assemble(AssemblerContext{Grammar: SyntacticGrammar}, d1, []Metadata{ttl1}...)
	if err != nil {
		t.Fatal(err)
	}
	// Assemble second
	ds2, err := SyntacticGrammar.Assemble(AssemblerContext{Grammar: SyntacticGrammar}, d2, []Metadata{ttl2}...)
	if err != nil {
		t.Fatal(err)
	}

	// Update
	dsu, err := ds1.UpdateWith(DefaultUpdateContext{}, ds2)
	m := dsu.Metadata()
	if m.TTL != uint64(ttlval2) {
		t.Fatal("TTL not updated successfully in node:", m.TTL, ttlval2)
	}
	if m.AssemblyTime > ds1.Metadata().AssemblyTime {
		t.Fatal("Assembly time update failed:", m.AssemblyTime, ds1.Metadata().AssemblyTime)
	}

	// Update without ttl being set in updated node
	ds2nottl, err := SyntacticGrammar.Assemble(AssemblerContext{Grammar: SyntacticGrammar}, d2)
	if err != nil {
		t.Fatal(err)
	}
	dsu, err = ds1.UpdateWith(DefaultUpdateContext{}, ds2nottl)
	if err != nil {
		t.Fatal(err)
	}
	m = dsu.Metadata()
	if m.TTL != uint64(sampleTTL) {
		t.Fatal("Update with no ttl not updated successfully in node:", m.TTL, sampleTTL)
	}
	if m.AssemblyTime > ds1.Metadata().AssemblyTime {
		t.Fatal("Assembly time update failed:", m.AssemblyTime, ds1.Metadata().AssemblyTime)
	}
}
