package swifthelper

import (
    "github.com/ncw/swift"
    "fmt"
    "crypto/md5"
    "io/ioutil"
    "encoding/hex"
    //"os"
    "path/filepath"
    //"bytes"
    //"encoding/json"
    //"swift_helper"
    //"strings"
)

func CreatePublicContainer(name string, c swift.Connection) {
    headers := map[string]string{
        "X-Container-Read": ".r:*",
    }
    c.ContainerCreate(name, headers)
}

func UploadFile(container string, prefix string, path string, c swift.Connection) {
    dat, err := ioutil.ReadFile(path)
    if (err != nil) {
        println(err.Error())
    } else {
        name := filepath.Base(path)
        if len(prefix) > 0 {
            name = prefix + "-" + name
        }
        ext := filepath.Ext(path)
        hasher := md5.New()
        hasher.Write(dat)
        md5hash := hex.EncodeToString(hasher.Sum(nil))

        fmt.Printf("Uploading %s to container %s\n", name, container)
        file, err := c.ObjectCreate(container, name, false, md5hash, ext, nil)
        if (err != nil) {
            println(err.Error())
        } else {
            file.Write(dat)
        }
        file.Close()
    }
}