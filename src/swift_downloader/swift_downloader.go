package main

import (
	"context"
	"fmt"
	"os"

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

	if len(os.Args) < 3 {
		fmt.Printf("Usage: %s project files...\n", os.Args[0])
		os.Exit(1)
	}

	println("Starting objekt storage downloader...")

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

	downloadFiles(ctx, container_name, files, *c)

}

func downloadFiles(ctx context.Context, container string, files []string, c swift.Connection) {
	wpauswiftcommons.CreatePublicContainer(ctx, container, c)
	names, err := c.ObjectNames(ctx, container, nil)

	if err != nil {
		println(err.Error())
		return
	}

	for _, name := range files {
		if !exists(names, name) {
			println("Object does not exists: " + name)
			continue
		}

		f, err := os.Create(name)
		if err != nil {
			println("Couldn't create file: " + name)
			continue
		}
		defer f.Close()

		_, err = c.ObjectGet(ctx, container, name, f, true, nil)
		if err != nil {
			println("Couldn't download file: " + name)
		} else {
			println("Downloaded file: " + name)
		}
	}
}

func exists(names []string, name string) bool {
	for _, e := range names {
		if e == name {
			return true
		}
	}
	return false
}
