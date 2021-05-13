package ir

import (
	"testing"
	"time"

	"github.com/libp2p/go-smart-record/xr"
)

const sampleTTL = 123

func TestDictMetadata(t *testing.T) {

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
	now := time.Now().Unix()
	m := ds.Metadata()
	if m.ExpirationTime < uint64(now) {
		t.Fatal("Expiration not set successfully in dict:", m.ExpirationTime, now)
	}

	elm := ds.(*Dict).Pairs[1].Key
	if elm.Metadata().ExpirationTime < uint64(now) {
		t.Fatal("Expiration not set successfully in dict element:", elm.Metadata().ExpirationTime, now)
	}

}

func TestSetMetadata(t *testing.T) {

	d := xr.Set{
		Tag: "aaa",
		Elements: xr.Nodes{
			xr.String{Value: "x"},
		},
	}

	ttl := TTL(uint64(sampleTTL))
	ds, err := SyntacticGrammar.Assemble(AssemblerContext{Grammar: SyntacticGrammar}, d, []Metadata{ttl}...)
	if err != nil {
		t.Fatal(err)
	}
	m := ds.Metadata()
	now := time.Now().Unix()
	if m.ExpirationTime < uint64(now) {
		t.Fatal("Expiration not set successfully in set:", m.ExpirationTime, now)
	}

	elm := ds.(*Set).Elements[0]
	if elm.Metadata().ExpirationTime < uint64(now) {
		t.Fatal("Expiration not set successfully in set element:", elm.Metadata().ExpirationTime, now)
	}
}

func TestDictMetadataUpdate(t *testing.T) {

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

	now := time.Now().Unix()
	// Update
	err = ds1.UpdateWith(DefaultUpdateContext{}, ds2)
	if err != nil {
		t.Fatal(err)
	}
	m := ds1.Metadata()
	// Updated with TTL of ds2
	if ds2.Metadata().ExpirationTime != m.ExpirationTime {
		t.Fatal("Expiration not updated successfully in node:", m.ExpirationTime, now)
	}

	ds1, err = SyntacticGrammar.Assemble(AssemblerContext{Grammar: SyntacticGrammar}, d1, []Metadata{ttl1}...)
	if err != nil {
		t.Fatal(err)
	}
	// Update without ttl being set in updated node
	ds2nottl, err := SyntacticGrammar.Assemble(AssemblerContext{Grammar: SyntacticGrammar}, d2)
	if err != nil {
		t.Fatal(err)
	}
	err = ds1.UpdateWith(DefaultUpdateContext{}, ds2nottl)
	if err != nil {
		t.Fatal(err)
	}
	m = ds1.Metadata()
	// Expiration should have not been updated.
	if m.ExpirationTime != ds1.Metadata().ExpirationTime {
		t.Fatal("Update with no ttl not updated successfully in node:", m.ExpirationTime, ds1.Metadata().ExpirationTime, now)
	}
}

func TestBasicMetadataUpdate(t *testing.T) {

	d1 := xr.String{Value: "test"}
	d2 := xr.String{Value: "test2"}

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

	now := time.Now().Unix()
	// Update
	err = ds1.UpdateWith(DefaultUpdateContext{}, ds2)
	if err != nil {
		t.Fatal(err)
	}
	m := ds1.Metadata()
	if ds2.Metadata().ExpirationTime != m.ExpirationTime {
		t.Fatal("Expiration not updated successfully in node:", m.ExpirationTime, now)
	}

	// Update without ttl being set in updated node
	ds1, err = SyntacticGrammar.Assemble(AssemblerContext{Grammar: SyntacticGrammar}, d1, []Metadata{ttl1}...)
	if err != nil {
		t.Fatal(err)
	}
	ds2nottl, err := SyntacticGrammar.Assemble(AssemblerContext{Grammar: SyntacticGrammar}, d2)
	if err != nil {
		t.Fatal(err)
	}
	err = ds1.UpdateWith(DefaultUpdateContext{}, ds2nottl)
	if err != nil {
		t.Fatal(err)
	}
	m = ds1.Metadata()
	if m.ExpirationTime != ds1.Metadata().ExpirationTime {
		t.Fatal("Update with no ttl not updated successfully in node:", m.ExpirationTime, ds1.Metadata().ExpirationTime, now)
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

	now := time.Now().Unix()
	// Update
	err = ds1.UpdateWith(DefaultUpdateContext{}, ds2)
	m := ds1.Metadata()
	if ds2.Metadata().ExpirationTime != m.ExpirationTime {
		t.Fatal("Expiration not updated successfully in node:", m.ExpirationTime, now)
	}

	ds1, err = SyntacticGrammar.Assemble(AssemblerContext{Grammar: SyntacticGrammar}, d1, []Metadata{ttl1}...)
	if err != nil {
		t.Fatal(err)
	}
	// Update without ttl being set in updated node
	ds2nottl, err := SyntacticGrammar.Assemble(AssemblerContext{Grammar: SyntacticGrammar}, d2)
	if err != nil {
		t.Fatal(err)
	}
	err = ds1.UpdateWith(DefaultUpdateContext{}, ds2nottl)
	if err != nil {
		t.Fatal(err)
	}
	m = ds1.Metadata()
	if m.ExpirationTime != ds1.Metadata().ExpirationTime {
		t.Fatal("Update with no ttl not updated successfully in node:", m.ExpirationTime, ds1.Metadata().ExpirationTime, now)
	}
}
