package main

import (
    "fmt"
    "net/url"
    "log"
    "github.com/fsouza/go-dockerclient"
    "github.com/atsaki/golang-cloudstack-library"
    "os"
)

const (
    WORKER_NAME = "wpau-worker"
)

func main() {
    // https://github.com/aerofs/gockerize

    log.Println("Starting cleanup")

    endpoint, _ := url.Parse("https://api.rbcloud.net/client/api")
    // TODO Lägg i env-vars..
    apikey := os.Getenv("RBC_API_KEY")
    secretkey := os.Getenv("RBC_SECRET")
    username := ""
    password := ""
    client, err := cloudstack.NewClient(endpoint, apikey, secretkey, username, password)

    if (err != nil) {
        panic(err)
    }

    workerCleanup(client)

}

func workerCleanup(client* cloudstack.Client) {
    log.Println("Worker cleanup sweep")
    params := cloudstack.NewListVirtualMachinesParameter()
    res, _ := client.ListVirtualMachines(params)

    for i := range res {
        if (res[i].Group.String() == WORKER_NAME) {
            ipadress := res[i].Nic[0].IpAddress.String()
            log.Printf("Found worker with id: %s, ip: %s, checking status of docker containers..\n", res[i].Id.String(), ipadress)
            hasRunningContainers, err := hasRunningContainers(ipadress)
            if (err != nil) {
                println(err.Error())
            } else {
                if (!hasRunningContainers) {
                    log.Printf("No running containers for id: %s, destroying...\n", res[i].Id.String())
                    destroyInstance(res[i].Id.String(), client)
                } else {
                    log.Printf("Has some running docker containers")
                }
            }
        }
    }
}

func destroyInstance(id string, client* cloudstack.Client) {
    params := cloudstack.NewDestroyVirtualMachineParameter(id)
    client.DestroyVirtualMachine(params)
}

func hasRunningContainers(ip string) (bool, error) {
    endpoint := fmt.Sprint("tcp://", ip, ":", 2375)
    client, err := docker.NewClient(endpoint)
    if (err != nil) {
        return false, err
    } else {
        opts := map[string][]string{ "status": []string{"running"}}
        containers, err := client.ListContainers(docker.ListContainersOptions{Filters: opts})
        if (err != nil) {
            return false, err
        } else {
            return len(containers) > 0, nil    
        }
    }
}

