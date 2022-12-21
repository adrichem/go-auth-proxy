# go-auth-proxy
[![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=adrichem_go-auth-proxy&metric=alert_status&token=00d5cfc6eb367ef9c92188c75bc6654dac8890cb)](https://sonarcloud.io/summary/new_code?id=adrichem_go-auth-proxy)

JWT token validation against AzureAD in a web-api.

It achieves 5.200 token validations per second while consuming <50% CPU and <30 MB of RAM.
1. Intel Core i7 @ 3Ghz.
1. Load generator is running on same machine 
1. 200 concurrent connections.