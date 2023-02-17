#!/bin/bash

POD=$(kubectl -n spire get po -l app=spire-server -o jsonpath="{.items[0].metadata.name}")

kubectl -n spire exec $POD -- /opt/spire/bin/spire-server entry create -spiffeID spiffe://example.org/ns/spire/sa/spire-agent -parentID spiffe://example.org/spire/server -selector k8s_psat:agent_ns:spire -selector k8s_psat:agent_sa:spire-agent -selector k8s_psat:cluster:example-cluster

kubectl -n spire exec $POD -- /opt/spire/bin/spire-server entry create -spiffeID spiffe://example.org/webserver -parentID spiffe://example.org/ns/spire/sa/spire-agent -selector k8s:container-name:spire-https-server -selector k8s:container-image:docker.io/abesharphpe/go-spiffe-https-example:v0.8

kubectl -n spire exec $POD -- /opt/spire/bin/spire-server entry create -spiffeID spiffe://example.org/webclient -parentID spiffe://example.org/ns/spire/sa/spire-agent -selector k8s:container-name:spire-client -selector k8s:container-image:docker.io/abesharphpe/go-spiffe-https-example:v0.8
