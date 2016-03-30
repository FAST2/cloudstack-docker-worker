package main

import (
    "github.com/ncw/swift"
    "fmt"
    //"crypto/md5"
    //"io/ioutil"
    //"encoding/hex"
    "os"
    //"path/filepath"
    //"bytes"
    //"encoding/json"
    "swifthelper"
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

    if (len(os.Args) < 3) {
        fmt.Printf("Usage: %s project files...\n", os.Args[0])
        os.Exit(1)
    }

    println("Starting objekt storage uploader")

    var (
        project = os.Args[1]
    )

    container_name := "project-" + project

    // Authenticate
    err := c.Authenticate()
    if err != nil {
        panic(err)
    }

    files := os.Args[2:]

    uploadFiles(container_name, files, c)

    
}

func uploadFiles(container string, files []string, c swift.Connection) {
    swifthelper.CreatePublicContainer(container, c)

    for _, e := range files {
        swifthelper.UploadFile(container, "", e, c)
    }
}
