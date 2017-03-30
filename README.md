# ampel manages shared resources
ampel (German for traffic light) is a simple application managing shared services
in a collaborative environment. A user can reserve a resource for a given amount of
time and free it after he finished using it.

Ampel uses HTTP status codes in conjunction with a very primitive web UI to allow
machines as well as humans to easily retrieve the status of resources.

## License
[New BSD license](LICENSE)
## Usage

### Command line
Usage:
    ampel service [service...]

Provide as many services at the command line as you want.

If there are two files "server.crt" and "server.key", ampel will use https and use
port 8443. If no certificate is provided, it will use HTTP over port 8080.

### Docker
The provided Dockerfile allows building a minimal container running ampel inside
that container.

Usage (HTTP):
    docker run --rm --name ampel -p 8080:8080 ampel service [service...]
Usage (HTTPS):
    docker run --rm --name ampel -p 8443:8443 -v /path/to/certificate.crt:/server.crt -v /path/to/certificate.key:/server.key ampel service [service...]

## Building

### Prerequisites
You need a go environment set up (see https://golang.org/doc/install). The
environment variable GOPATH should be set. If you want to build the docker
container, you also need to install docker and be able to run docker.

### compiling
There is a simple build script build.sh, which builds a static binary "ampel" and
if docker is installed, a docker container tagged "ampel". You can provide an
alternative tag as commandline parameter
  ./build.sh registry.mycompany.com/ampel:latest


## Endpoints
This application supports the following endpoints:
GET / delivers a list of available services

GET /service delivers a list of available services

GET /service/$SERVICENAME returns the status of the service $SERVICENAME
The HTTP status codes reflect the service status
- 200/OK is returned, if the resource is free,
- 423/Locked, if it is locked,
- If the requested resource is not known, 404/Not found is returned.

POST /service/$SERVICENAME/stop sets the status of a single service to "Locked".
The HTTP status codes reflect the service status
- 403/Denied is returned, if the resource is already locked,
- 200/OK, if it is was locked successfully.
- If the requested resource is not known, 404/Not found is returned.

POST /service/$SERVICENAME/go sets the status of a single service to "Free".
The HTTP status codes reflect the service status
- 409/Conflict is returned, if the resource is already free,
- 200/OK, if it is was freed successfully.
-  If the requested resource is not known, 404/Not found is returned.

GET /list delivers a list of available services

GET /list/available or
GET /list/free returns a list of free services


GET /list/reserved or
GET /list/stopped returns a list of stopped services

If you add ?output=json to the GET requests, instead of a web page, data is returned as JSON
objects

## Credits
* Karl Peglau (original Ampelmann design); Matthew Gates (SVG version)
* [Wikipedia](https://en.wikipedia.org/wiki/Ampelm%C3%A4nnchen), pointing me to the file
