# power

[![GoDoc](https://godoc.org/github.com/scjalliance/power?status.svg)](https://godoc.org/github.com/scjalliance/power)

Power infrastructure health sampling and monitoring.

## Example Docker Invocation

```
docker run -d --name=power-monitor --restart=always -e SOURCE=lcy-rack2n-ups,lcy-rack2s-ups -e COMMUNITY=tripplite -e INTERVAL=1m -e RECIPIENT=stathat:STATHATKEY scjalliance/power
```
