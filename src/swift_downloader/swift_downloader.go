package main

import (
    "github.com/ncw/swift"
    "fmt"
    "os"
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

    println("Starting objekt storage downloader")

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

    downloadFiles(container_name, files, c)

    
}

func downloadFiles(container string, files []string, c swift.Connection) {
    swifthelper.CreatePublicContainer(container, c)
    names, err := c.ObjectNames(container, nil)

    if (err != nil) {
        println(err.Error())
        return
    }

    for _, name := range files {
        if (!exists(names, name)) {
            println("Object does not exists: " + name)
            continue
        }
        
        f, err := os.Create(name)
        if (err != nil) {
            println("Couldn't create file: " + name)
            continue
        }
        defer f.Close()

        _, err = c.ObjectGet(container, name, f, true, nil)
        if (err != nil) {
            println("Couldn't download file: " + name)
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
