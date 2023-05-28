# gocurl
## A simple curl written in golang
can be used to obtain the tlsKey during the https handshake


## usage
```shell
gocurl -X POST -H "Content-Type: application/json" -d '{"name":"test"}' -p "tlsKey.txt" "https://localhost:8080
```