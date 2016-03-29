package main

import (
    "fmt"
    "net/url"
    "log"
    "github.com/fsouza/go-dockerclient"
    "github.com/atsaki/golang-cloudstack-library"
    "os"
    "net/http"
    "bytes"
    "strconv"
    "time"
)

const (
    WORKER_NAME = "wpau-worker"
    WARM_UP_MINUTES = 15
    WARNING_NO_HOURS = 6
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
            id := res[i].Id.String()
            ipadress := res[i].Nic[0].IpAddress.String()
            duration_running, _ := getDurationRunning(res[i].Created.String())
            log.Printf("Found worker with id: %s, ip: %s, checking status of docker containers..\n", id, ipadress)
            isWarmedUp, err := hasWarmUp(res[i].Created.String())

            if (err != nil) {
                println(err.Error())
                continue
            }

            if (!isWarmedUp) {
                log.Printf("Container with id %s is not yet warmed up, doing nothing", id)
                continue
            }

            hasRunningContainers, err := hasRunningContainers(ipadress)
            if (err != nil) {
                println(err.Error()) 
                continue
            }

            if (!hasRunningContainers) {
                log.Printf("No running containers for id: %s, destroying...\n", res[i].Id.String())
                destroyInstance(res[i].Id.String(), client)
                sendStatus("Destroyed instance with id " + id + " ip " + ipadress)
            } else {
                log.Printf("Has some running docker containers, wont destroy")
                if (duration_running.Hours() > WARNING_NO_HOURS) {
                    sendStatus("Instance with id " + id + " ip " + ipadress + " has been running over 6 hours, normal?")
                }
            }
        }
    }
}

func hasWarmUp(datetime string) (bool, error) {
    const layout = "2006-01-02T15:04:05Z0700"
    t, err := time.Parse(layout, datetime)
    if (err != nil) {
        return false, err
    } else {
        return time.Now().After( t.Add(time.Duration(WARM_UP_MINUTES) * time.Minute)), nil
    }
}

func getDurationRunning(datetime string) (time.Duration, error) {
    const layout = "2006-01-02T15:04:05Z0700"
    t, err := time.Parse(layout, datetime)
    if (err != nil) {
        return time.Second, err
    } else {
        return time.Since(t), nil
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

func sendStatus(msg string) {
    apiUrl := os.Getenv("WPAU_SLACK_HOOK_URL")
    data := url.Values{}
    data.Set("payload", "{\"username\":\"WPAU-robot\", \"icon_emoji\":\":speaking_head_in_silhouette:\", \"text\":\"" + msg + "\"}")

    client := &http.Client{}
    r, _ := http.NewRequest("POST", apiUrl, bytes.NewBufferString(data.Encode()))
    r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
    r.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

    client.Do(r)
}

