package main

import (
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/digitalocean/godo"
	"golang.org/x/oauth2"
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

func getIp() (string, error) {
	res, err := http.Get("https://api.ipify.org")
	if err != nil {
		return "", err
	}
	ip, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	return string(ip), nil
}

func main() {

	// Flags with no serious DEFAULT value
	tokenPtr := flag.String("token", "", "Digital Ocean Token")
	domainPtr := flag.String("domain", "", "Domain that will be maintained with current IP")
	recordNamePtr := flag.String("name", "", "Record name that should be used to idenfity the IP")

	flag.Parse()

	ip, err := getIp()

	if err != nil {
		log.Fatal("Couldn't get my IP address:", err)
		os.Exit(1)
	}

	log.Printf("Current IP: %s\n", ip)

	tokenSource := &TokenSource{
		AccessToken: *tokenPtr,
	}

	oauthClient := oauth2.NewClient(oauth2.NoContext, tokenSource)
	client := godo.NewClient(oauthClient)

	domain, _, err := client.Domains.Get(*domainPtr)

	if err != nil {
		log.Fatalf("Domain %s doesn't exist. Please create the domain first.\n", *domainPtr)
		os.Exit(1)
	}
	log.Printf("Domain %s exists.\n", domain)

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
