# go-auth-proxy
[![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=adrichem_go-auth-proxy&metric=alert_status&token=00d5cfc6eb367ef9c92188c75bc6654dac8890cb)](https://sonarcloud.io/summary/new_code?id=adrichem_go-auth-proxy)

## Objective
JWT token validation against AzureAD to secure some upstream web-api. Created as Proof of Concept
to compare against a NodeJs version.

## Rough performance measurement
It achieves 12.000 token validations per second while consuming <50% CPU and <30 MB of RAM.
1. Intel Core i7 3Ghz.
1. Load generator  running on same machine configured with 200 concurrent connections.
1. Proxy and load generator are competing for CPU time. Together they max out the CPU.