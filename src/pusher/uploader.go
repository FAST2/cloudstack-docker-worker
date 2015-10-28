package main

import (
    "github.com/ncw/swift"
    "fmt"
    "crypto/md5"
    "io/ioutil"
    "encoding/hex"
    "os"
    "path/filepath"
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
        fmt.Printf("Usage: %s path container_name\n", os.Args[0])
        os.Exit(1)
    }

    println("Starting objekt storage uploader")

    var (
        path = os.Args[1]
        container_name = os.Args[2]
        )

    // Authenticate
    err := c.Authenticate()
    if err != nil {
        panic(err)
    }

    uploadContentsInFolder(path, container_name, c)

    // List all the containers
    containers, err := c.ContainerNames(nil)
    fmt.Println(containers)
}

func getContentFromContainer(container string, c swift.Connection) {
    names, _ := c.ObjectNamesAll(container, nil)

    for i := range names {
        println(names[i])
        c.ObjectGet(container, names[i], os.Stdout, true, nil)
    }
}


func uploadContentsInFolder(path string, container string, c swift.Connection) {
    createContainer(container, c)

    err := filepath.Walk(path, func(path string, f os.FileInfo, err error) error {
        uploadFile(container, path, c)
        return nil
    })
    if (err != nil) {
        println(err.Error())
    }
}

func createContainer(name string, c swift.Connection) {
    c.ContainerCreate(name, nil)
}

func uploadFile(container string, path string, c swift.Connection) {
    dat, err := ioutil.ReadFile(path)
    if (err != nil) {
        println(err.Error())
    } else {
        name := filepath.Base(path)
        ext := filepath.Ext(path)
        hasher := md5.New()
        hasher.Write(dat)
        md5hash := hex.EncodeToString(hasher.Sum(nil))

        fmt.Printf("Uploading %s to container %s", name, container)
        file, err := c.ObjectCreate(container, name, false, md5hash, ext, nil)
        if (err != nil) {
            println(err.Error())
        } else {
            file.Write(dat)
        }
        file.Close()
    }
}
