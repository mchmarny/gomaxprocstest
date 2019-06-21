# gomaxprocstest

Simple golang concurrency test app. It spins up CPU-intensive calculation for defined number of `goroutines`. This app is helpful in assessing runtime performance in virtualized envirnments where often the number of host vCPUs are exposed to app runtime as `runtime.NumCPU()` but the actual number of CPUs available to the application is actually capped at lower number (pften 1 vCPU) which may lead to perception of low performance.

REST endpoints:

* `/` - shows total number of CPUs "visible" to go runtime
* `/cores/:core/concurrency/:count/calcs/:calc` - where
  * `:core` represents the number of max cores to set for this request
  * `:count` represents the number of concurrent goroutines to execute
  * `:calc` represents the number of mathematical operations to perform (each op includes both `+` and `-`)

As an example, this request:

`/cores/4/concurrency/4/calcs/1000000000`

Will result in response looking something like this:

```json
{
    "total_cores": 4,
    "max_cores": 4,
    "duration": "962.592169ms",
    "messages": [
        {
            "goroutine": 2,
            "message": "Done: 903.902341ms"
        },
        {
            "goroutine": 4,
            "message": "Done: 956.269309ms"
        },
        {
            "goroutine": 1,
            "message": "Done: 958.983235ms"
        },
        {
            "goroutine": 3,
            "message": "Done: 962.547894ms"
        }
    ]
}
```

## Deploy

### Local

```shell
go run main.go
```

http://localhost:8080/cores

### Cloud Run on GKE

```shell
gcloud beta run deploy gomaxprocs \
    --image gcr.io/cloudylabs-public/gomaxprocs:0.1.3 \
    --platform gke \
    --cluster cr \
    --cluster-location us-east1
```


### Cloud Run on GKE (w/ 1 vCPU resource limit)

```shell
kubectl apply -f service.yaml
```


## Cloud Run (managed)

```shell
gcloud beta run deploy gomaxprocs \
	--image=gcr.io/cloudylabs-public/gomaxprocs:0.1.3 \
	--region=us-central1
```


