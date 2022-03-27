package main

import (
	"context"
	"fmt"
	"os"

	swift "github.com/ncw/swift/v2"
	"github.com/oquinena/wpauswiftcommons"
)

func main() {
	ctx := context.Background()
	// Create a connection using openstack v3applicationcredential
	c := &swift.Connection{
		ApplicationCredentialId:     os.Getenv("OS_APPLICATION_CREDENTIAL_ID"),
		ApplicationCredentialSecret: os.Getenv("OS_APPLICATION_CREDENTIAL_SECRET"),
		AuthUrl:                     os.Getenv("OS_AUTH_URL"),
	}

	// We need at least two arguments. Project to upload to and file(s) to upload
	if len(os.Args) < 3 {
		fmt.Printf("Usage: %s [project] [files]...\n", os.Args[0])
		os.Exit(1)
	}

	fmt.Println("Starting objekt storage uploader...")

	var (
		project = os.Args[1]
	)

	container_name := "project-" + project

	// Authenticate
	err := c.Authenticate(ctx)
	if err != nil {
		fmt.Printf("Error authenticating: %s\n", err)
		os.Exit(1)
	}

	files := os.Args[2:]

	uploadFiles(ctx, container_name, files, *c)

}

func uploadFiles(ctx context.Context, container string, files []string, c swift.Connection) {
	wpauswiftcommons.CreatePublicContainer(ctx, container, c)

	for _, e := range files {
		wpauswiftcommons.UploadFile(ctx, container, "", e, c)
	}
}
