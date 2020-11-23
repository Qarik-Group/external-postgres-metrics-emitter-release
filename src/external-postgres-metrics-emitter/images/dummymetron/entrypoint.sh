#!/bin/sh

/bin/dummymetron \
   --cert /certs/loggregator_agent.crt \
   --key /certs/loggregator_agent.key \
   --ca /certs/loggregator_ca.crt \
   --grpc-port 3459 \
   "${@}"
