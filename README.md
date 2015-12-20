# Digital Ocean Dns Updater (DynDNS style)

Application (or more a script) used to update a Digital Ocean domain name entry. Uses https://www.ipify.org/ to determine the external IP address and the the Digital Ocean Go API for the dom

    dodns -domain=DOMAIN -name=SUBDOMAIN -token=DO_TOKEN

So the following call
 
    dodns -domain=example.com -name=somesubdomain -token=110d11ca11dac1161191a18a011ed3d111f11861130117b7fd1168ba112f119c

would result in keeping
    
    somesubdomain.example.com

up-to-date with the current external IP address.

## Installing

Simply use `go install`

    go install github.com/mrbuk/dodns

Make sure that $GOPATH is set before.

## Running periodically

To execute the script periodically you can use cron or any other scheduling mechanism that is able to execute a shell command. E.g. to run the application every 5 min. use

    */5 *  *   *   *     $GOBIN/dodns -domain=example.com -name=somesubdomain -token=110d11ca11dac1161191a18a011ed3d111f11861130117b7fd1168ba112f119c >> /var/log/dodns.log 2>&1

