# gomaxprocstest

Simple golang concurrency and parallelism test app. It spins CPU-intensive calculation in a goroutine, per each CPU.

REST endpoints:

`/` - shows total number of CPUs "visible" to go runtime
`/cores` - runs with max `GOMAXPROCS` set to total number of CPUs
`/cores/4` - runs with max `GOMAXPROCS` set to passed value (e.g. `4`)

## MacBook Pro (Intel Core i7, 3.5 GHz)

```shell
go run main.go
```

http://localhost:8080/cores

```shell
CPU cores: 4
Core 0 start
Core 1 start
Core 2 start
Core 3 start
Core 3 done in 12.618280296s
Core 2 done in 12.669301014s
Core 1 done in 12.721737848s
Core 0 done in 12.72762304s
Total duration: 12.727656509s
```

## Cloud Run on GKE

```shell
gcloud beta run deploy gomaxprocs \
    --image gcr.io/cloudylabs-public/gomaxprocs:0.1.2 \
    --platform gke \
    --cluster cr \
    --cluster-location us-east1
```

http://gomaxprocs.default.knative.tech/run/cores

```shell
CPU cores: 4
Core 0 start
Core 1 start
Core 2 start
Core 3 start
Core 0 done in 12.002515815s
Core 2 done in 12.157338676s
Core 1 done in 12.960494389s
Core 3 done in 13.051117131s
Total duration: 13.051166992s
```

## Cloud Run on GKE (w/ 1 vCPU resource limit)

```shell
kubectl apply -f service.yaml
```

https://gomaxprocs.demo.knative.tech/run/cores

```shell
CPU cores: 4
Core 0 start
Core 1 start
Core 2 start
Core 3 start
Core 0 done in 34.793659559s
Core 1 done in 36.10635037s
Core 2 done in 36.327062228s
Core 3 done in 36.439878516s
Total duration: 36.439928693s
```


## Cloud Run

```shell
gcloud beta run deploy gomaxprocs \
	--image=gcr.io/cloudylabs-public/gomaxprocs:0.1.2 \
	--region=us-central1
```

https://gomaxprocs-2gtouos2pq-uc.a.run.app/run/cores

```shell
CPU cores: 8
Core 0 start
Core 1 start
Core 2 start
Core 3 start
Core 4 start
Core 5 start
Core 6 start
Core 7 start
Core 1 done in 1m47.657496904s
Core 3 done in 1m47.847700716s
Core 4 done in 1m48.744847984s
Core 2 done in 1m49.467559325s
Core 7 done in 1m49.766156001s
Core 0 done in 1m49.860066382s
Core 5 done in 1m50.293746955s
Core 6 done in 1m50.381939829s
Total duration: 1m50.382130567s
```
