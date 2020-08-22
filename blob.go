// Package blobbucket provides an implementation of FSDB Bucket using Go-Cloud
// Blob interface.
package blobbucket

import (
	"context"
	"io"

	"github.com/fishy/errbatch"
	"github.com/fishy/fsdb/bucket"
	"gocloud.dev/blob"
)

// Make sure *BlobBucket satisifies bucket.Bucket interface.
var _ bucket.Bucket = (*BlobBucket)(nil)

// BlobBucket is an implementation of FSDB Bucket with Go-Cloud Blob interface.
type BlobBucket struct {
	bkt *blob.Bucket
}

// Open opens a blob bucket.
func Open(bkt *blob.Bucket) *BlobBucket {
	return &BlobBucket{
		bkt: bkt,
	}
}

func (bkt *BlobBucket) Read(
	ctx context.Context,
	name string,
) (io.ReadCloser, error) {
	return bkt.bkt.NewReader(ctx, name, nil)
}

func (bkt *BlobBucket) Write(
	ctx context.Context,
	name string,
	data io.Reader,
) error {
	writer, err := bkt.bkt.NewWriter(ctx, name, &blob.WriterOptions{
		ContentType: "application/octet-stream",
	})
	if err != nil {
		return err
	}
	if _, err := io.Copy(writer, data); err != nil {
		var batch errbatch.ErrBatch
		batch.Add(err)
		batch.Add(writer.Close())
		return batch.Compile()
	}
	return writer.Close()
}

// Delete deletes an object from the bucket.
func (bkt *BlobBucket) Delete(ctx context.Context, name string) error {
	return bkt.bkt.Delete(ctx, name)
}

// IsNotExist checks whether err means the object does not exist on the blob.
func (bkt *BlobBucket) IsNotExist(err error) bool {
	return blob.IsNotExist(err)
}
