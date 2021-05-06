package ir

// MetadataType interface that metadata field types need to
// implement.
type metadataType interface {
	update(with metadataType) metadataType // Determines how the metadata is updated when the node is updated.
}

// ttl determines the expiration time of a node.
type ttl struct {
	value uint64
}

// TTL sets a TTL in seconds to the node in metadata
func TTL(value uint64) Metadata {
	return func(m *metadataContext) error {
		m.ttl = ttl{value}
		return nil
	}
}

// update for ttl
func (t ttl) update(with metadataType) metadataType {
	withT, ok := with.(ttl)
	// If entered wrong type to update do nothing and return metadata as-is
	if !ok {
		return t
	}
	// If with.ttl is zero, it means it is not set. Do not update
	if withT.value == 0 {
		return t
	}
	t.value = withT.value
	return t
}

// assemblyTime specifies that time of assembly of a node
type assemblyTime struct {
	value uint64
}

// update for asasemblyTime
func (a assemblyTime) update(with metadataType) metadataType {
	withT, ok := with.(assemblyTime)
	// If entered wrong type to update do nothing and return metadata as-is
	if !ok {
		return a
	}
	// Use the largest of the two as the final assemblyTime.
	if a.value < withT.value {
		a.value = withT.value
	}
	return a
}
