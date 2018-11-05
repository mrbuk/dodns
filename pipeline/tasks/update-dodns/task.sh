#!/bin/bash

dodns -domain ${domain} -name ${subdomain} -token ${apitoken} -ip=$(cat myip/ip)
