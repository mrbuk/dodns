package main

import (
    "flag"
    "log"
    "os"
    "golang.org/x/oauth2"
    "github.com/digitalocean/godo"
    "github.com/rdegges/go-ipify"
)

type TokenSource struct {
    AccessToken string
}

func (t *TokenSource) Token() (*oauth2.Token, error) {
    token := &oauth2.Token{
        AccessToken: t.AccessToken,
    }
    return token, nil
}


func main() {
    
    ip, err := ipify.GetIp()

    if err != nil {
        log.Fatal("Couldn't get my IP address:", err)
        os.Exit(1)
    }

    log.Println("Current IP: ", ip)

    tokenPtr := flag.String("token", "DEFAULT", "Digital Ocean Token")
    domainPtr := flag.String("domain", "DEFAULT", "Domain that will be maintained with current IP")
    recordNamePtr := flag.String("name", "DEFAULT", "Record name that should be used to idenfity the IP")
    
    flag.Parse()

    tokenSource := &TokenSource{
        AccessToken: *tokenPtr,
    }

    oauthClient := oauth2.NewClient(oauth2.NoContext, tokenSource)
    client := godo.NewClient(oauthClient)

    domain, _, err := client.Domains.Get(*domainPtr)

    if err != nil {
        log.Fatal("Domain ", *domainPtr, " doesn't exist. Please create the domain first.")
        os.Exit(1)
    }
    log.Println("Domain ", domain, " exists.")

    records, _, err := client.Domains.Records(*domainPtr, nil)
    if err != nil {
        log.Fatal(err)
        os.Exit(1)
    }

    // check if home exists and points the current ip
    var record godo.DomainRecord
    for key, value := range records {
        log.Println("Record: ", key, " - ", value)
        if value.Name == *recordNamePtr {
            record = value
        }
    }

    createRequest := &godo.DomainRecordEditRequest{
        Type:     "A",
        Name:     *recordNamePtr,
        Data:     ip,
        Priority: 0,
        Port:     0,
        Weight:   0,
    }

    // update home dns record
    if record.Name != "" {
        if record.Data != ip {
            log.Println("Old IP address ", record.Data, " found. Updating with current IP: ", ip)
            client.Domains.EditRecord(*domainPtr, record.ID, createRequest)
        } else {
            log.Println("IP address is up to date.")
        }
    } else {
        log.Println("Record with name '", *recordNamePtr, "' does not exist for domain '",
             *domainPtr, ". Creating new record.")
        client.Domains.CreateRecord(*domainPtr, createRequest)
    }

    
}
