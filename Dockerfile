FROM scratch
MAINTAINER Joern Ott <joern.ott@ott-consult.de>
COPY ampel /
COPY static /static

EXPOSE 8080
EXPOSE 8443

ENTRYPOINT ["/ampel"]
