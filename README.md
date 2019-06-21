# gomaxprocstest

Simple golang concurrency and parallelism test app. It spins CPU-intensive calculation in a goroutine, per each CPU.

REST endpoints:

`/` - shows total number of CPUs "visible" to go runtime
`/cores` - runs with max `GOMAXPROCS` set to total number of CPUs
`/cores/2` - runs with max `GOMAXPROCS` set to passed value (e.g. `4`)

The response will look something like this

```json
{
    "total_cores": 4,
    "max_cores": 2,
    "duration": "59.606265ms",
    "messages": [
        {
            "core": 2,
            "message": "Done: 58.945942ms"
        },
        {
            "core": 1,
            "message": "Done: 59.462734ms"
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


