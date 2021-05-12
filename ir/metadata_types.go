package ir

import "time"

const MaxUINT64 = 1<<64 - 1

// MetadataType interface implemented by metadata field types
type metadataType interface {
	update(with metadataType) metadataType // Determines how the metadata is updated when the node is updated.
}

// expirationTime determines the expiration time of a node.
type expirationTime struct {
	value uint64
}

// TTL sets a TTL in seconds to the node in metadata. It triggers a change
// to the node's expirationTime
func TTL(value uint64) Metadata {
	return func(m *metadataContext) error {
		m.expirationTime.value = uint64(time.Now().Unix()) + value
		return nil
	}
}

// update logic for expirationTime metadata type
func (t expirationTime) update(with metadataType) metadataType {
	withT, ok := with.(expirationTime)
	// If entered wrong type to update do nothing and return metadata as-is
	if !ok {
		return t
	}
	// If new expirationTime below the current one, do not update.
	if withT.value < t.value {
		return t
	}
	t.value = withT.value
	return t
}

// NoExpiration prevents the node from having an expirationTime
func NoExpiration() Metadata {
	return func(m *metadataContext) error {
		m.expirationTime.value = uint64(MaxUINT64)
		return nil
	}

}
