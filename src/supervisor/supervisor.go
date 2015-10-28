 package fast2.se

import (
    //"fmt"
    "net/url"
    "log"
    "io/ioutil"
    "github.com/atsaki/golang-cloudstack-library"
    "os"
)

const (
    WORKER_NAME = "wpau-worker"
)

func main() {
    endpoint, _ := url.Parse("https://api.rbcloud.net/client/api")
    apikey := os.Getenv("RBC_API_KEY")
    secretkey := os.Getenv("RBC_SECRET")
    username := ""
    password := ""
    client, err := cloudstack.NewClient(endpoint, apikey, secretkey, username, password)

    if (err != nil) {
        panic(err)
    }

    startJob(client)

    //var (
    //    customer = os.Args[1]
    //)
}

func startJob(client* cloudstack.Client) {
    log.Println("Starting job..")
    userdata, _ := getUserdata()
    serviceOfferingMicro := "1bd74b58-ac1e-46e4-86cb-a9542064b8a4"
    defaultZone := "806945e8-2431-4526-9d1c-70748f287439"
    //networkId := "19313259-68af-4d65-9e28-1249ee60887a"
    //defaultZone := "19313259-68af-4d65-9e28-1249ee60887a"
    //ubuntuTemplate := "643ccc7d-87e8-4c65-9c1d-2df68a23e82d"
    //ubuntuTemplate := "497cacef-6492-4130-bd15-45748c0a4864"
    //ubuntuTemplate := "rbc/ubuntu-14.04-server-cloudimg-amd64-20GB-20153214"
    ubuntuTemplate := "497cacef-6492-4130-bd15-45748c0a4864" // Det är den som ligger ovan..
    //serviceOfferingMini := "84d98576-17c7-4bc4-831b-27ceec3e35bc"
    params := cloudstack.NewDeployVirtualMachineParameter(serviceOfferingMicro, ubuntuTemplate, defaultZone)
    params.KeyPair.Set("ubuntu")
    params.UserData.Set(userdata)
    params.Group.Set("wpau-worker")

    _, err := client.DeployVirtualMachine(params)
    if (err != nil) {
        log.Printf("Couldn't create/deploy new instance, error from API: %s", err.Error())
    }
}

func getUserdata() (string, error) {
    dat, err := ioutil.ReadFile("./cloud-config.txt")
    if err == nil {
        return string(dat), nil
    } else {
        return "", err
    }
}