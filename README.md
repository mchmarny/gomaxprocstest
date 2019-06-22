# gorun

In virtualized environments the number of virtual CPUs (vCPUs) exposed to the application runtime (e.g. `runtime.NumCPU()` in Go) is based on the total number of vCPUs available on the underlined host node. The actual number of vCPUs available to your application is actually capped at much lower number to prevent one application from consuming all of the compute resources of the host node.

That means that if you set the number of threads in your code (e.g `GOMAXPROCS` in Go) to anything above the capped number (e.g. `runtime.GOMAXPROCS(runtime.NumCPU())`), your app performance will be degraded with some correlation to the number of goroutines you run in your application. This happens because the Go scheduler will try to distribute each one of the goroutines over multiple worker threads on every available processor.

> Note, starting with Go 1.5+ and 1.6+, `GOMAXPROCS` is set to runtime.NumCPU() by default!

This simple golang concurrency test app helps asses performance impact under different combinations of cores, goroutines, and number of CPU-intensive calculations.

More on scheduling in Go [here](https://www.ardanlabs.com/blog/2018/08/scheduling-in-go-part1.html)

## REST endpoints

* `GET /` - shows total number of CPUs "visible" to go runtime
* `GET /cores/:core/concurrency/:count/calcs/:calc` - where
  * `:core` represents the number of max cores to set for this request
  * `:count` represents the number of concurrent goroutines to execute
  * `:calc` represents the number of mathematical operations to perform (each op includes both `+` and `-`)

As an example, this request:

`/cores/4/concurrency/4/calcs/1000000000`

Will run `1000000000` mathematical calculations in `4` separate goroutines with `runtime.GOMAXPROCS` set to `4` and result in response looking something like this:

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

### Local (using docker)

```shell
docker run -p 8080:8080 mchmarny/gorun:0.1.5
```

### Cloud Run on GKE

```shell
gcloud beta run deploy gorun \
    --image gcr.io/cloudylabs-public/gorun:0.1.5 \
    --platform gke \
    --cluster cr \
    --cluster-location us-east1
```

### Cloud Run on GKE (w/ 1 vCPU resource limit)

```shell
kubectl -n demo apply -f \
    https://raw.githubusercontent.com/mchmarny/gorun/master/service.yaml
```

### Cloud Run (managed)

```shell
gcloud beta run deploy gorun \
    --image=gcr.io/cloudylabs-public/gorun:0.1.5 \
    --region=us-central1 \
    --platform=managed
```


