#!/bin/bash

set -e

manifest="
cert: ((loggregator_tls_agent.certificate))
key: ((loggregator_tls_agent.private_key))
ca: ((loggregator_tls_agent.ca))
variables:
- name: loggregator_ca
  type: certificate
  options:
    common_name: loggregatorCA
    is_ca: true
    duration: 3600
- name: loggregator_tls_agent
  type: certificate
  update_mode: converge
  options:
    ca: loggregator_ca
    common_name: metron
    alternative_names:
    - metron
    extended_key_usage:
    - client_auth
    - server_auth
    duration: 3600
"

bosh int --vars-store /tmp/vars-store.yml \
     <(echo "${manifest}") --path /cert > loggregator_agent.crt

bosh int --vars-store /tmp/vars-store.yml \
     <(echo "${manifest}") --path /key > loggregator_agent.key

bosh int --vars-store /tmp/vars-store.yml \
     <(echo "${manifest}") --path /ca > loggregator_ca.crt

