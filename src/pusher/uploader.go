package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/FAST2/wpaumetadata"
	"github.com/FAST2/wpauswiftcommons"
	"github.com/ncw/swift/v2"
)

func main() {
	ctx := context.Background()
	// Create a connection using openstack v3applicationcredential
	c := &swift.Connection{
		ApplicationCredentialId:     os.Getenv("OS_APPLICATION_CREDENTIAL_ID"),
		ApplicationCredentialSecret: os.Getenv("OS_APPLICATION_CREDENTIAL_SECRET"),
		AuthUrl:                     os.Getenv("OS_AUTH_URL"),
	}

	if len(os.Args) < 5 {
		fmt.Printf("Usage: %s [customer] [jobId] [status] [path]\n", os.Args[0])
		os.Exit(1)
	}

	println("Starting object storage uploader...")

	var (
		customer = os.Args[1]
		jobId    = os.Args[2]
		status   = os.Args[3]
		path     = os.Args[4]
	)

	container_name := "jobs-" + customer
	file_prefix := jobId

	// Authenticate
	err := c.Authenticate(ctx)
	if err != nil {
		fmt.Printf("Error authenticating: %s\n", err)
		os.Exit(1)
	}

	uploadContentsInFolder(ctx, path, file_prefix, container_name, *c)
	wpaumetadata.Add(ctx, *c, container_name, jobId, status)
}

func uploadContentsInFolder(ctx context.Context, path string, prefix string, container string, c swift.Connection) {
	wpauswiftcommons.CreatePublicContainer(ctx, container, c)

	err := filepath.Walk(path, func(subpath string, f os.FileInfo, err error) error {
		if !f.IsDir() {
			wpauswiftcommons.UploadFile(ctx, container, prefix, subpath, c)
		}
		return nil
	})
	if err != nil {
		println(err.Error())
	}
}
