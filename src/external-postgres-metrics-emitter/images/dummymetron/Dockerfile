FROM golang

ADD src dummymetron
RUN cd dummymetron && GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o /bin/dummymetron main.go 

FROM alpine:latest  
COPY --from=0 /bin/dummymetron /bin/dummymetron

RUN mkdir /certs
ADD loggregator_agent.crt loggregator_agent.key loggregator_ca.crt /certs/
ADD entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

ENTRYPOINT /entrypoint.sh