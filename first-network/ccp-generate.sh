#!/bin/bash

function one_line_pem {
    echo "`awk 'NF {sub(/\\n/, ""); printf "%s\\\\\\\n",$0;}' $1`"
}

function json_ccp {
    local PP=$(one_line_pem $5)
    local CP=$(one_line_pem $6)
    sed -e "s/\${ORG}/$1/" \
        -e "s/\${P0PORT}/$2/" \
        -e "s/\${P1PORT}/$3/" \
        -e "s/\${CAPORT}/$4/" \
        -e "s#\${PEERPEM}#$PP#" \
        -e "s#\${CAPEM}#$CP#" \
        -e "s#\${ORGC}#$7#" \
        ccp-template.json 
}

function yaml_ccp {
    local PP=$(one_line_pem $5)
    local CP=$(one_line_pem $6)
    sed -e "s/\${ORG}/$1/" \
        -e "s/\${P0PORT}/$2/" \
        -e "s/\${P1PORT}/$3/" \
        -e "s/\${CAPORT}/$4/" \
        -e "s#\${PEERPEM}#$PP#" \
        -e "s#\${ORGC}#$7#" \
        ccp-template.yaml | sed -e $'s/\\\\n/\\\n        /g'
}

ORG=manufacturer
ORGC=Manufacturer
P0PORT=7051
P1PORT=8051
CAPORT=7054
PEERPEM=crypto-config/peerOrganizations/manufacturer.pharma-network.com/tlsca/tlsca.manufacturer.pharma-network.com-cert.pem
CAPEM=crypto-config/peerOrganizations/manufacturer.pharma-network.com/ca/ca.manufacturer.pharma-network.com-cert.pem

echo "$(json_ccp $ORG $P0PORT $P1PORT $CAPORT $PEERPEM $CAPEM $ORGC)" > connection-manufacturer.json
echo "$(yaml_ccp $ORG $P0PORT $P1PORT $CAPORT $PEERPEM $CAPEM $ORGC)" > connection-manufacturer.yaml

ORG=distributor
ORGC=Distributor
P0PORT=9051
P1PORT=10051
CAPORT=8054
PEERPEM=crypto-config/peerOrganizations/distributor.pharma-network.com/tlsca/tlsca.distributor.pharma-network.com-cert.pem
CAPEM=crypto-config/peerOrganizations/distributor.pharma-network.com/ca/ca.distributor.pharma-network.com-cert.pem

echo "$(json_ccp $ORG $P0PORT $P1PORT $CAPORT $PEERPEM $CAPEM $ORGC)" > connection-distributor.json
echo "$(yaml_ccp $ORG $P0PORT $P1PORT $CAPORT $PEERPEM $CAPEM $ORGC)" > connection-distributor.yaml

ORG=transporter
ORGC=Transporter
P0PORT=11051
P1PORT=12051
CAPORT=9054
PEERPEM=crypto-config/peerOrganizations/transporter.pharma-network.com/tlsca/tlsca.transporter.pharma-network.com-cert.pem
CAPEM=crypto-config/peerOrganizations/transporter.pharma-network.com/ca/ca.transporter.pharma-network.com-cert.pem

echo "$(json_ccp $ORG $P0PORT $P1PORT $CAPORT $PEERPEM $CAPEM $ORGC)" > connection-transporter.json
echo "$(yaml_ccp $ORG $P0PORT $P1PORT $CAPORT $PEERPEM $CAPEM $ORGC)" > connection-transporter.yaml

ORG=retailer
ORGC=Retailer
P0PORT=13051
P1PORT=14051
CAPORT=10054
PEERPEM=crypto-config/peerOrganizations/retailer.pharma-network.com/tlsca/tlsca.retailer.pharma-network.com-cert.pem
CAPEM=crypto-config/peerOrganizations/retailer.pharma-network.com/ca/ca.retailer.pharma-network.com-cert.pem

echo "$(json_ccp $ORG $P0PORT $P1PORT $CAPORT $PEERPEM $CAPEM $ORGC)" > connection-retailer.json
echo "$(yaml_ccp $ORG $P0PORT $P1PORT $CAPORT $PEERPEM $CAPEM $ORGC)" > connection-retailer.yaml
