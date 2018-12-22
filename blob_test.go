package blobbucket_test

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/fishy/blobbucket"
	"gocloud.dev/blob/fileblob"
)

func Example() {
	dir, err := ioutil.TempDir("", "blob-bucket-test")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(dir)
	fs, err := fileblob.OpenBucket(dir, nil)
	if err != nil {
		log.Fatal(err)
	}
	ctx := context.Background()

	bkt := blobbucket.Open(fs)
	obj := "test/object"
	content := `Lorem ipsum dolor sit amet,
consectetur adipiscing elit,
sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.

Ut enim ad minim veniam,
quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat.

Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur.

Excepteur sint occaecat cupidatat non proident,
sunt in culpa qui officia deserunt mollit anim id est laborum.`

	_, err = bkt.Read(ctx, obj)
	fmt.Println("Read IsNotExist:", bkt.IsNotExist(err))

	if err := bkt.Write(ctx, obj, strings.NewReader(content)); err != nil {
		log.Fatal(err)
	}
	reader, err := bkt.Read(ctx, obj)
	if err != nil {
		log.Fatal(err)
	}
	defer reader.Close()
	buf, err := ioutil.ReadAll(reader)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Content unchanged:", content == string(buf))
	if err := bkt.Delete(ctx, obj); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Delete IsNotExist:", bkt.IsNotExist(bkt.Delete(ctx, obj)))
	// Output:
	// Read IsNotExist: true
	// Content unchanged: true
	// Delete IsNotExist: true
}
