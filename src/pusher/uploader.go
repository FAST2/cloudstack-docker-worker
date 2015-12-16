package main

import (
    "github.com/ncw/swift"
    "fmt"
    "crypto/md5"
    "io/ioutil"
    "encoding/hex"
    "os"
    "path/filepath"
    "bytes"
    "encoding/json"
    //"strings"
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

    updateStatusfile(container_name, jobId, status, c)
}

type Jobinfo struct {
    JobId string `json:"jobId"`
    Status string `json:"status"`
}


type Jobinfos struct {
    Infos [] Jobinfo `json:"jobs"`
}


func updateStatusfile(container string, jobId string, status string, c swift.Connection) {
    buf := new(bytes.Buffer)
    _, err := c.ObjectGet(container, "jobs.json", buf, true, nil)

    info := Jobinfo{jobId, status}
    s, _ := json.Marshal(info)
    println(string(s))
    var content []byte = nil

    if (err != nil) {
        println("jobs.json doesn't exists, will create new file")
        // Does not exists, corruped or otherwise, recreate it
        infos := Jobinfos{}
        infos.Infos = append(infos.Infos, info)
        b, err := json.Marshal(infos)
        if (err != nil) {
            println("Couldn't marshal json data, won't create jobs.json")
            return
        } else {
            content = b
        }
    } else {
        println("File exists, will append to it")
        infos := Jobinfos{}
        err := json.Unmarshal(buf.Bytes(), &infos)
        if (err != nil) {
            println("Couldn't unmarshal current json data, won't update jobs.json")
            return
        }
        infos.Infos = append(infos.Infos, info)
        b, err := json.Marshal(infos)
        if (err != nil) {
            println("Couldn't marshal json data, won't create jobs.json")
            return
        } else {
            content = b
        }
    }

    hasher := md5.New()
    hasher.Write(content)
    md5hash := hex.EncodeToString(hasher.Sum(nil))

    file, err := c.ObjectCreate(container, "jobs.json", false, md5hash, "json", nil)
    if (err != nil) {
        println(err.Error())
    } else {
        file.Write(content)
    }
    file.Close()
}

func getContentFromContainer(container string, c swift.Connection) {
    names, _ := c.ObjectNamesAll(container, nil)

    for i := range names {
        println(names[i])
        c.ObjectGet(container, names[i], os.Stdout, true, nil)
    }
}


func uploadContentsInFolder(path string, prefix string, container string, c swift.Connection) {
    createContainer(container, c)

    err := filepath.Walk(path, func(subpath string, f os.FileInfo, err error) error {
        if (!f.IsDir()) {
            uploadFile(container, prefix, subpath, c)
        }
        return nil
    })
    if (err != nil) {
        println(err.Error())
    }
}

func createContainer(name string, c swift.Connection) {
    headers := map[string]string{
        "X-Container-Read": ".r:*",
    }
    c.ContainerCreate(name, headers)
}

func uploadFile(container string, prefix string, path string, c swift.Connection) {
    dat, err := ioutil.ReadFile(path)
    if (err != nil) {
        println(err.Error())
    } else {
        name := prefix + "-" + filepath.Base(path)
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
