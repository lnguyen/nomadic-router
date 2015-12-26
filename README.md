# Nomadic Router

Nomadic router is used with in conjunction with Nomadic service broker to handle deploying of clustered services on Cloud Foundry.

The router job is to listen to consul for new services and register them to [LVS](http://www.linuxvirtualserver.org/). Which allows services to be advertised via a static ip. The LVS used is [Gorb](https://github.com/kobolog/gorb).


# Running tests

```
go get github.com/onsi/ginkgo/ginkgo
ginkgo -r
```

# Usage
```
nomadic-router -c 127.0.0.0:8500 -g http://127.0.0.1:4672 -ip 192.168.1.1
```