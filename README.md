# Regent.dev

Regent.dev ia software which bootstraps a non-IoT device into an IoT enabled one, or enabled rapid development of a new
IoT device.

A program on the device (the agent) acts as a supervisor program, running the other processes that you require.
The agent program monitors yours processes, restarting them if required.
The agent securely connects to the regent.io platform via MQTT/tls & HTTPS.
The regent.dev platform allows your devices to be remotely administered.

## Supported Platforms

Regent.io agent operates on any platform supported by Golang including Intel, Arm, Mips & RiscV hardware and
Linux & Windows OS.

# Platform deployment

Before deploying the platform, you need to generate CA and MQTT TLS certificates and keys. 
Run generate_certificates script located in platform/scripts:

```
./generate_certificates.sh domain_name
```

where domain_name is the domain of the platform. For local deployment use "localhost". generate_certificates
required openssl to be installed. After generating scripts you don't have to do anything else, the script will 
copy them into correct locations.

The platform consists of multiple docker containers tied together by docker-compose.yml file 
(https://docs.docker.com/compose/). To build and deploy type:

```
docker-compose up -d --build
```

This will build and deploy all containers and volumes. The containers require multiple environment variables to be set.
Variables for each container are set in docker-compose environment variables files such as: ".db.env.dev" located
in the same directory as docker-compose.yml.

After deploying the platform you can:

* See the dashboard at [localhost](http://localhost)
* See Django admin panel at [localhost/admin](http://localhost/admin)
* See the browsable version of the public API at [localhost/api](http://localhost/api)

In order to tear down the deployment type:

```
docker-compose down
```

This will stop all containers, delete the network and remove images. Note that this command will not remove volumes, 
in particular it will not remove the volume in which postgres container stores data. To tear down deployment 
and remove volumes type:

```
docker-compose down -v
```

Few notes about docker-compose.

When docker-compose deploys containers, it creates a docker network (https://docs.docker.com/network/) in a bridge mode. 
Each container has hostname identical to its service name.
