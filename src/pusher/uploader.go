package main

import (
    "github.com/ncw/swift"
    "fmt"
    "os"
    "path/filepath"
    "github.com/fast2/wpauswiftcommons"
    "github.com/fast2/wpaumetadata"
)

func main() {
    // Create a connection
    c := swift.Connection{
        UserName: os.Getenv("SWIFT_API_USER"),
        ApiKey:   os.Getenv("SWIFT_API_KEY"),
        AuthUrl:  os.Getenv("SWIFT_AUTH_URL"),
        Domain:   "",  // Name of the domain (v3 auth only)
        Tenant:   "",  // Name of the tenant (v2 auth only)
    }

    if (len(os.Args) < 5) {
        fmt.Printf("Usage: %s customer jobId status path\n", os.Args[0])
        os.Exit(1)
    }

    println("Starting objekt storage uploader")

    var (
        customer = os.Args[1]
        jobId = os.Args[2]
        status = os.Args[3]
        path = os.Args[4]
        )

    container_name := "jobs-" + customer
    file_prefix := jobId

    // Authenticate
    err := c.Authenticate()
    if err != nil {
        panic(err)
    }

    uploadContentsInFolder(path, file_prefix, container_name, c)
    wpaumetadata.Add(c, container_name, jobId, status)
}

func uploadContentsInFolder(path string, prefix string, container string, c swift.Connection) {
    wpauswiftcommons.CreatePublicContainer(container, c)

    err := filepath.Walk(path, func(subpath string, f os.FileInfo, err error) error {
        if (!f.IsDir()) {
            wpauswiftcommons.UploadFile(container, prefix, subpath, c)
        }
        return nil
    })
    if (err != nil) {
        println(err.Error())
    }
}

