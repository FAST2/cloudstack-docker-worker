package main

import (
    //"fmt"
    "net/url"
    "log"
    //"github.com/fsouza/go-dockerclient"
    "github.com/atsaki/golang-cloudstack-library"
    "os"
    "net/http"
    "bytes"
    "strconv"
    "time"
)

const (
    WORKER_NAME = "wpau-worker"
    WARNING_NO_HOURS = 6
)

func main() {
    // https://github.com/aerofs/gockerize

    log.Println("Starting cleanup")

    endpoint, _ := url.Parse("https://api.rbcloud.net/client/api")
    apikey := os.Getenv("RBC_API_KEY")
    secretkey := os.Getenv("RBC_SECRET")
    username := ""
    password := ""
    client, err := cloudstack.NewClient(endpoint, apikey, secretkey, username, password)

    if (err != nil) {
        panic(err)
    }

    workerCheckup(client)

}

func workerCheckup(client* cloudstack.Client) {
    log.Println("Worker checkup")
    params := cloudstack.NewListVirtualMachinesParameter()
    res, _ := client.ListVirtualMachines(params)

    for i := range res {
        if (res[i].Group.String() == WORKER_NAME) {
            id := res[i].Id.String()
            ipadress := res[i].Nic[0].IpAddress.String()
            duration_running, _ := getDurationRunning(res[i].Created.String())
            log.Printf("Found worker with id: %s, ip: %s, checking uptime..\n", id, ipadress)

            if (duration_running.Hours() > WARNING_NO_HOURS) {
                sendStatus("Instance with id " + id + " ip " + ipadress + " has been running over 6 hours, normal?")
            }
        }
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

