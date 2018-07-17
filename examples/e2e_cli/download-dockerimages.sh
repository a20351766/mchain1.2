#!/bin/bash -eu
#
# Copyright IBM Corp. All Rights Reserved.
#
# SPDX-License-Identifier: Apache-2.0
#


##################################################
# This script pulls docker images from hyperledger
# docker hub repository and Tag it as
# hyperledger/mchain-<image> latest tag
##################################################

dockerFabricPull() {
  local FABRIC_TAG=$1
  for IMAGES in peer orderer couchdb ccenv javaenv kafka tools zookeeper; do
      echo "==> FABRIC IMAGE: $IMAGES"
      echo
      docker pull hyperledger/mchain-$IMAGES:$FABRIC_TAG
      docker tag hyperledger/mchain-$IMAGES:$FABRIC_TAG hyperledger/mchain-$IMAGES
  done
}

dockerCaPull() {
      local CA_TAG=$1
      echo "==> FABRIC CA IMAGE"
      echo
      docker pull hyperledger/fabric-ca:$CA_TAG
      docker tag hyperledger/fabric-ca:$CA_TAG hyperledger/fabric-ca
}
usage() {
      echo "Description "
      echo
      echo "Pulls docker images from hyperledger dockerhub repository"
      echo "tag as hyperledger/mchain-<image>:latest"
      echo
      echo "USAGE: "
      echo
      echo "./download-dockerimages.sh [-c <fabric-ca tag>] [-f <mchain tag>]"
      echo "      -c fabric-ca docker image tag"
      echo "      -f mchain docker image tag"
      echo
      echo
      echo "EXAMPLE:"
      echo "./download-dockerimages.sh -c 1.1.1 -f 1.1.0"
      echo
      echo "By default, pulls the 'latest' fabric-ca and mchain docker images"
      echo "from hyperledger dockerhub"
      exit 0
}

while getopts "\?hc:f:" opt; do
  case "$opt" in
     c) CA_TAG="$OPTARG"
        echo "Pull CA IMAGES"
        ;;

     f) FABRIC_TAG="$OPTARG"
        echo "Pull FABRIC TAG"
        ;;
     \?|h) usage
        echo "Print Usage"
        ;;
  esac
done

: ${CA_TAG:="latest"}
: ${FABRIC_TAG:="latest"}

echo "===> Pulling mchain Images"
dockerFabricPull ${FABRIC_TAG}

echo "===> Pulling mchain ca Image"
dockerCaPull ${CA_TAG}
echo
echo "===> List out hyperledger docker images"
docker images | grep hyperledger*
