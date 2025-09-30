package updates

import (
	"bytes"
	"compress/zlib"
	"context"
	"fmt"
	"io"
)

// BinaryDiffer implements binary diffing
type BinaryDiffer struct{}

// NewDiffer creates a new differ
func NewDiffer() *BinaryDiffer {
	return &BinaryDiffer{}
}

// Generate generates a diff between two binaries
func (d *BinaryDiffer) Generate(ctx context.Context, oldBinary, newBinary []byte) (*Diff, error) {
	if oldBinary == nil || newBinary == nil {
		return nil, fmt.Errorf("binaries cannot be nil")
	}

	// Simplified diff generation (in production, use bsdiff or similar)
	diffData := computeSimpleDiff(oldBinary, newBinary)

	diff := &Diff{
		Size:       int64(len(diffData)),
		Checksum:   fmt.Sprintf("sha256-%d", len(diffData)),
		Data:       diffData,
		Algorithm:  "simple",
		Compressed: false,
	}

	return diff, nil
}

// Apply applies a diff to a binary
func (d *BinaryDiffer) Apply(ctx context.Context, oldBinary []byte, diff *Diff) ([]byte, error) {
	if oldBinary == nil || diff == nil {
		return nil, fmt.Errorf("oldBinary and diff cannot be nil")
	}

	// Decompress if needed
	data := diff.Data
	if diff.Compressed {
		decompressed, err := d.Decompress(ctx, diff)
		if err != nil {
			return nil, fmt.Errorf("decompression failed: %w", err)
		}
		data = decompressed.Data
	}

	// Apply diff (simplified - in production, use bspatch or similar)
	newBinary := applySimpleDiff(oldBinary, data)

	return newBinary, nil
}

// Compress compresses a diff
func (d *BinaryDiffer) Compress(ctx context.Context, diff *Diff) (*Diff, error) {
	if diff == nil || diff.Compressed {
		return diff, nil
	}

	var buf bytes.Buffer
	w := zlib.NewWriter(&buf)
	if _, err := w.Write(diff.Data); err != nil {
		return nil, fmt.Errorf("compression failed: %w", err)
	}
	if err := w.Close(); err != nil {
		return nil, fmt.Errorf("compression failed: %w", err)
	}

	compressed := &Diff{
		FromVersion: diff.FromVersion,
		ToVersion:   diff.ToVersion,
		Size:        int64(buf.Len()),
		Checksum:    fmt.Sprintf("sha256-%d", buf.Len()),
		Data:        buf.Bytes(),
		Algorithm:   diff.Algorithm,
		Compressed:  true,
	}

	return compressed, nil
}

// Decompress decompresses a diff
func (d *BinaryDiffer) Decompress(ctx context.Context, diff *Diff) (*Diff, error) {
	if diff == nil || !diff.Compressed {
		return diff, nil
	}

	r, err := zlib.NewReader(bytes.NewReader(diff.Data))
	if err != nil {
		return nil, fmt.Errorf("decompression failed: %w", err)
	}
	defer r.Close()

	var buf bytes.Buffer
	if _, err := io.Copy(&buf, r); err != nil {
		return nil, fmt.Errorf("decompression failed: %w", err)
	}

	decompressed := &Diff{
		FromVersion: diff.FromVersion,
		ToVersion:   diff.ToVersion,
		Size:        int64(buf.Len()),
		Checksum:    fmt.Sprintf("sha256-%d", buf.Len()),
		Data:        buf.Bytes(),
		Algorithm:   diff.Algorithm,
		Compressed:  false,
	}

	return decompressed, nil
}

// Helper functions for simple diff (in production, use proper binary diff algorithm)
func computeSimpleDiff(old, new []byte) []byte {
	// Simplified: just store the new binary
	return new
}

func applySimpleDiff(old, diff []byte) []byte {
	// Simplified: just return the diff (which is the new binary)
	return diff
}
