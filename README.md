# go-auth-proxy
[![SonarCloud](https://sonarcloud.io/images/project_badges/sonarcloud-black.svg)](https://sonarcloud.io/summary/new_code?id=adrichem_go-auth-proxy)

[![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=adrichem_go-auth-proxy&metric=alert_status&token=00d5cfc6eb367ef9c92188c75bc6654dac8890cb)](https://sonarcloud.io/summary/new_code?id=adrichem_go-auth-proxy)
[![Bugs](https://sonarcloud.io/api/project_badges/measure?project=adrichem_go-auth-proxy&metric=bugs&token=00d5cfc6eb367ef9c92188c75bc6654dac8890cb)](https://sonarcloud.io/summary/new_code?id=adrichem_go-auth-proxy)
[![Code Smells](https://sonarcloud.io/api/project_badges/measure?project=adrichem_go-auth-proxy&metric=code_smells&token=00d5cfc6eb367ef9c92188c75bc6654dac8890cb)](https://sonarcloud.io/summary/new_code?id=adrichem_go-auth-proxy)
[![Vulnerabilities](https://sonarcloud.io/api/project_badges/measure?project=adrichem_go-auth-proxy&metric=vulnerabilities&token=00d5cfc6eb367ef9c92188c75bc6654dac8890cb)](https://sonarcloud.io/summary/new_code?id=adrichem_go-auth-proxy)


## Objective
Secure some upstream web-api with AzureAd authentication. Created as Proof of Concept to compare against a NodeJs version.

It validates the JWT token is: 
1. Valid untampered AzureAD token.
2. Valid in its lifetime (exp and nbf claims)
3. Issued by your AzureAD tenant (iss claim)
4. Issued to your expected audience (aud claim)

## Running in docker

```bash
docker run -ti --rm -p80:80 adrichem/go-auth-proxy `
    --upstream https://my-web-api.com `
    --header-value my-secret-api-key `
    --header-name Apikey `
    --aud expected-value-for-aud-claim `
    --iss expected-value-for-iss-claim
```

## Very rough performance measurement
It achieves 12.000 token validations per second while consuming <50% CPU and <30 MB of RAM. Docker image size is <15 MB.
1. Intel Core i7 3Ghz.
1. 200 concurrent connections.
1. Load generator running on same machine and competing for CPU time. Together they max out the CPU.
