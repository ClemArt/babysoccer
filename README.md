# Local development

## Docker

Create machine

    ./create-env.sh

Add proxy

    <docker-machine start cmd> --engine-env HTTP_PROXY=$PROXY --engine-env HTTPS_PROXY=$PROXY --engine-env NO_PROXY=$NO_PROXY


## Utils

__WINDOWS__ : Link folder, to ease sharing

    cmd > mklink /J <Link> <Target Folder>
    ex :  mklink /J workspace E:\perso\workspace