package kqueuey

// CompressionOrdinal specifies how a block in badgerDB should be compressed.
type CompressionOrdinal = uint8

const (
	None CompressionOrdinal = iota
	Snappy
	ZSTD
)

// Compression types badgerDB supports.
const (
	CompressionNone   = "None"
	CompressionSnappy = "Snappy"
	CompressionZSTD   = "ZSTD"
)

func checkCompressionType(compressionType string) CompressionOrdinal {
	if compressionType == CompressionNone {
		return None
	}

	if compressionType == CompressionSnappy {
		return Snappy
	}

	if compressionType == CompressionZSTD {
		return ZSTD
	}

	// If the compression type set to badgers default compression type.
	return Snappy
}
