package ir

import (
	"fmt"
)

// metadataContext includes all the metadata fields supported by semantic nodes.
// This context is attached to semantic nodes. The context supports private and
// public metadata fields. For a metadata field to be reported publicly in a node
// it needs to be registered in MetadataInfo.
// NOTE: In the future, we should consider using a more flexible type for
// metadataContext so that anyone can easily register their own metadata in their
// nodes. For instance, we can use a map[string]metadataType and add a .RegisterMetadataType.
type metadataContext struct {
	expirationTime expirationTime // Timestamp of expiration of the node.
}

// MetadataInfo is a container for the reporting of the current
// public metadata of a semantic node. We report directly the metadata
// internal value type, not the metadataType
type MetadataInfo struct {
	ExpirationTime uint64
}

// Metadata option applies metaadata to a smart node.
type Metadata func(*metadataContext) error

// Applies supported metadata items to a metadataContext
// Using this apply function and separating MetdataInfo (reporting purposes)
// from metadataContext (data in semantic node) we avoid someone from being
// able to manipulate metadata in the MetadataCtx directly.
func (m *metadataContext) apply(items ...Metadata) error {
	for i, item := range items {
		if err := item(m); err != nil {
			return fmt.Errorf("error applying metadata value %d: %s", i, err)
		}
	}
	return nil
}

// getMetadata returns public metadata in a context as MetadataInfo
func (m *metadataContext) getMetadata() MetadataInfo {
	if m == nil {
		return MetadataInfo{}
	}

	return MetadataInfo{
		ExpirationTime: m.expirationTime.value,
	}
}

// update the metadata of a node conveniently when it receives an update.
func (m *metadataContext) update(with *metadataContext) {
	// If any of the nodes doesn't have metadata -> return
	if m == nil || with == nil {
		return
	}
	m.expirationTime = m.expirationTime.update(with.expirationTime).(expirationTime)
}

func (m *metadataContext) copy() metadataContext {
	if m == nil {
		return metadataContext{}
	}
	return *m
}
