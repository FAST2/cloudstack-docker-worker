package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/FAST2/wpaumetadata"
	"github.com/ncw/swift/v2"
)

const MAX_JOBS = 60 // Two months

func main() {
	ctx := context.Background()
	// Create a connection using openstack v3applicationcredential
	c := &swift.Connection{
		ApplicationCredentialId:     os.Getenv("OS_APPLICATION_CREDENTIAL_ID"),
		ApplicationCredentialSecret: os.Getenv("OS_APPLICATION_CREDENTIAL_SECRET"),
		AuthUrl:                     os.Getenv("OS_AUTH_URL"),
	}

	println("Starting cleaner")

	// Authenticate
	err := c.Authenticate(ctx)
	if err != nil {
		fmt.Printf("Error authenticating: %s\n", err)
		os.Exit(1)
	}

	containers := get_projects(ctx, *c)

	fmt.Println("Found job containers:")
	fmt.Println(containers)

	for _, val := range containers {
		clean_old_jobs(ctx, *c, val)
	}
}

func get_projects(ctx context.Context, c swift.Connection) []string {
	opts := new(swift.ContainersOpts)
	opts.Prefix = "jobs"

	containers, err := c.ContainerNames(ctx, opts)
	if err != nil {
		panic(err)
	}
	return containers
}

func clean_old_jobs(ctx context.Context, c swift.Connection, container string) {
	buf, err := wpaumetadata.GetMetadata(ctx, c, container)
	if err != nil {
		fmt.Printf("No metadata file for %s\n", container)
		return
	}
	infos, err := wpaumetadata.ParseMetadata(buf)

	if err != nil {
		panic(err)
	}

	jobs_count := len(infos.Infos)
	fmt.Printf("Container %s has %d nr of jobs\n", container, jobs_count)

	if jobs_count > MAX_JOBS {
		remove_count := jobs_count - MAX_JOBS
		fmt.Printf("Will remove %d jobs from container %s\n", remove_count, container)
		for i := 0; i < remove_count; i++ {
			prefix := infos.Infos[i].JobId
			fmt.Printf("Will iterate over contents of %s with prefix %s for removal of items\n", container, prefix)
			objects_to_remove := get_objects(ctx, c, container, prefix)
			fmt.Printf("Items to be removed: %s\n", objects_to_remove)
			remove_objects(ctx, c, container, objects_to_remove)
			infos = remove_from_metadata(infos, prefix)
		}
	}

	updated_metadata, err := json.Marshal(infos)
	if err != nil {
		fmt.Printf("Couldn't marshal to JSON with updated for %s", container)
	} else {
		wpaumetadata.Upload(ctx, c, container, updated_metadata)
	}

}

func get_objects(ctx context.Context, c swift.Connection, container string, prefix string) []string {
	opts := new(swift.ObjectsOpts)
	opts.Prefix = prefix
	names, err := c.ObjectNames(ctx, container, opts)
	if err != nil {
		panic(err)
	}
	return names
}

func remove_objects(ctx context.Context, c swift.Connection, container string, names []string) bool {
	removed_all := true
	for _, val := range names {
		fmt.Printf("Removing %s\n", val)
		err := c.ObjectDelete(ctx, container, val)
		if err != nil {
			fmt.Printf("Couldn't delete '%s' because: %s\n", val, err.Error())
			removed_all = false
		}
	}
	return removed_all
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
