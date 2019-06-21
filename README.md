# gorun

Simple golang concurrency test app. It spins up CPU-intensive calculation for defined number of `goroutines`. This app is helpful in assessing runtime performance in virtualized environments where often the number of host vCPUs are exposed to app runtime as `runtime.NumCPU()` but the actual number of CPUs available to the application is actually capped at lower number (often 1 vCPU) which may lead to perception of low performance.

REST endpoints:

* `/` - shows total number of CPUs "visible" to go runtime
* `/cores/:core/concurrency/:count/calcs/:calc` - where
  * `:core` represents the number of max cores to set for this request
  * `:count` represents the number of concurrent goroutines to execute
  * `:calc` represents the number of mathematical operations to perform (each op includes both `+` and `-`)

As an example, this request:

`/cores/4/concurrency/4/calcs/1000000000`

Will run `1000000000` mathematical calculations in `4` separate `goroutines` with `runtime.GOMAXPROCS` set to `4` and result in response looking something like this:

```json
{
    "available_cores": 4,
    "max_cores": 2,
    "concurrency": 4,
    "calculations": 1000000000,
    "duration": "1.706786234s",
    "details": [
        {
            "goroutine": 1,
            "duration": "1.390726667s"
        },
        {
            "goroutine": 2,
            "duration": "1.673043201s"
        },
        {
            "goroutine": 4,
            "duration": "1.673078606s"
        },
        {
            "goroutine": 3,
            "duration": "1.673091242s"
        }
    ]
}
```

## Deploy

### Local

```shell
go run main.go
```

### Cloud Run on GKE

```shell
gcloud beta run deploy gorun \
    --image gcr.io/cloudylabs-public/gorun:0.1.4 \
    --platform gke \
    --cluster cr \
    --cluster-location us-east1
```


### Cloud Run on GKE (w/ 1 vCPU resource limit)

```shell
kubectl apply -f service.yaml
```


### Cloud Run (managed)


```shell
gcloud beta run deploy gorun \
	--image=gcr.io/cloudylabs-public/gorun:0.1.4 \
	--region=us-central1
```


