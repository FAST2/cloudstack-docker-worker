package main

import (
    "github.com/ncw/swift"
    "fmt"
    "os"
    "github.com/fast2/wpaumetadata"
    "encoding/json"
)

const MAX_JOBS = 109

func main() {
    // Create a connection
    c := swift.Connection{
        UserName: os.Getenv("SWIFT_API_USER"),
        ApiKey:   os.Getenv("SWIFT_API_KEY"),
        AuthUrl:  os.Getenv("SWIFT_AUTH_URL"),
        Domain:   "",  // Name of the domain (v3 auth only)
        Tenant:   "",  // Name of the tenant (v2 auth only)
    }

    println("Starting cleaner")

    println(os.Getenv("SWIFT_API_USER"))

    // Authenticate
    err := c.Authenticate()
    if err != nil {
        panic(err)
    }

    containers := get_projects(c)

    fmt.Println("Found job containers:")
    fmt.Println(containers)

    for _, val := range containers {
        clean_old_jobs(c, val)
    }
}

func get_projects(c swift.Connection) []string {
    opts := new(swift.ContainersOpts)
    //opts.Prefix = "jobs"

    containers, err := c.ContainerNames(opts)
    if (err != nil) {
        panic(err)
    }
    return containers
}

func clean_old_jobs(c swift.Connection, container string) {
    buf, err := wpaumetadata.GetMetadata(c, container)
    if (err != nil) {
        fmt.Printf("No metadata file for %s\n", container)
        return
    }
    infos, err := wpaumetadata.ParseMetadata(buf)

    if (err != nil) {
        panic(err)
    }

    jobs_count := len(infos.Infos)
    fmt.Printf("Container %s has %d nr of jobs\n", container, jobs_count)

    if (jobs_count > MAX_JOBS) {
        remove_count := jobs_count - MAX_JOBS
        fmt.Printf("Will remove %d jobs from container %s\n", remove_count, container)
        for i := 0; i < remove_count; i++ {
            prefix := infos.Infos[i].JobId
            fmt.Printf("Will iterate over contents of %s with prefix %s for removal of items\n", container, prefix)
            objects_to_remove := get_objects(c, container, prefix)
            fmt.Printf("Removing %s\n", objects_to_remove)
            //remove_objects(c, container, objects_to_remove)
            infos = remove_from_metadata(infos, prefix)
        }
    }

    if (false) {

    
    updated_metadata, err := json.Marshal(infos)
    if (err != nil) {
        fmt.Printf("Couldn't marshal to JSON with updated for %s", container)
    } else {
        wpaumetadata.Upload(c, container, updated_metadata)
    }
}

}

func get_objects(c swift.Connection, container string, prefix string) []string {
    opts := new(swift.ObjectsOpts)
    opts.Prefix = prefix
    names, err := c.ObjectNames(container, opts)
    if (err != nil) {
        panic(err)
    }
    return names
}

func remove_objects(c swift.Connection, container string, names []string) bool {
    res, err := c.BulkDelete(container, names)
    if (err != nil) {
        return res.NumberDeleted == int64(len(names))
    }
    return false
}

func remove_from_metadata(infos wpaumetadata.Jobinfos, id_to_remove string) wpaumetadata.Jobinfos {
    for index, val := range infos.Infos {
        if val.JobId == id_to_remove {
            infos.Infos = append(infos.Infos[:index], infos.Infos[index+1:]...)
            break
        }
    }
    return infos
}
