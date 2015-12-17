 package main

import (
    "fmt"
    "net/url"
    "log"
    "io/ioutil"
    "github.com/atsaki/golang-cloudstack-library"
    "os"
    "strings"
    "net/http"
    "bytes"
    "strconv"
)

func main() {
    if (len(os.Args) < 4) {
        fmt.Printf("Usage: %s docker-repository image-name customer\n", os.Args[0])
        os.Exit(1)
    }

    var (
        repo string = os.Args[1]
        worker_name string = os.Args[2]
        customer string = os.Args[3]
    )


    endpoint, _ := url.Parse("https://api.rbcloud.net/client/api")
    apikey := os.Getenv("RBC_API_KEY")
    secretkey := os.Getenv("RBC_SECRET")
    username := ""
    password := ""
    client, err := cloudstack.NewClient(endpoint, apikey, secretkey, username, password)

    if (err != nil) {
        panic(err)
    }

    startJob(customer, repo, worker_name, client)
}

func startJob(customer string, repo string, worker_name string, client* cloudstack.Client) {
    log.Println("Starting job..")
    userdata, err := generateUserdata(repo, worker_name, customer)
    if (err != nil) {
        log.Println("Couldn't read user data template file, stopping..")
        return
    }
    serviceOfferingMedium := "8c7a1b14-19d3-4f91-a59b-08064e2b5692"
    defaultZone := "806945e8-2431-4526-9d1c-70748f287439"
    ubuntuTemplate := "497cacef-6492-4130-bd15-45748c0a4864"
    params := cloudstack.NewDeployVirtualMachineParameter(serviceOfferingMedium, ubuntuTemplate, defaultZone)
    params.KeyPair.Set("ubuntu")
    params.UserData.Set(userdata)
    params.Group.Set(worker_name)

    _, err = client.DeployVirtualMachine(params)
    if (err != nil) {
        log.Printf("Couldn't create/deploy new instance, error from API: %s", err.Error())
    } else {
        sendStatus("Created new instance in da cloud for customer " + customer)
        
    }
}

func getUserdataTemplate() (string, error) {
    dat, err := ioutil.ReadFile("./cloud-config-template.txt")
    if err == nil {
        return string(dat), nil
    } else {
        return "", err
    }
}

func generateUserdata(repo string, worker_name string, customer string) (string, error) {
    content, err := getUserdataTemplate()
    if (err != nil) {
        return "", err
    } else {
        content = strings.Replace(content, "__DOCKER_REPO__", repo, -1)
        content = strings.Replace(content, "__WORKER_NAME__", worker_name, -1)
        content = strings.Replace(content, "__CUSTOMER__", customer, -1)
        content = strings.Replace(content, "__SWIFT_API_USER__", os.Getenv("SWIFT_API_USER"), -1)
        content = strings.Replace(content, "__SWIFT_API_KEY__", os.Getenv("SWIFT_API_KEY"), -1)
        content = strings.Replace(content, "__SWIFT_AUTH_URL__", os.Getenv("SWIFT_AUTH_URL"), -1)
    }
    return content, nil
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


