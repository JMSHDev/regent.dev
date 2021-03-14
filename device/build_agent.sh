#!/bin/bash
# this is a script to automate building the agent and ExampleApp

# REQUIREMENTS: You must have Golang 1.15+ installed

SCRIPTPATH="$( cd "$(dirname "$0")" >/dev/null 2>&1 ; pwd -P )"
cd $SCRIPTPATH

# build the agent and place it in the build directory
cd agent || exit
go build -ldflags="-linkmode internal" -o ../bin/agent
cd ..

# build the example application and place it in the build directory
cd example_app || exit
go build -ldflags="-linkmode internal" -o ../bin/ExampleApp -i ExampleApp.go
cd ..

# copy an example configuration in place if there isn't one already present
[ -f ./bin/config.json ] && echo "agent config already exist." || { echo "copying example agent config"; cp ./agent/example/config.json ./bin; }
