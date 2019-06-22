RELEASE=0.1.6

.PHONY: mod
mod:
	go mod tidy
	go mod vendor

debug:
	go run -v main.go

profile:
	go tool pprof http://localhost:8080/perf/profile

run:
	docker run -p 8080:8080 mchmarny/gorun:$(RELEASE)

build: mod
	docker build -t mchmarny/gorun:latest \
		-t mchmarny/gorun:$(RELEASE) .
	docker push mchmarny/gorun:$(RELEASE)

image: mod
	gcloud builds submit \
		--project cloudylabs-public \
		--tag gcr.io/cloudylabs-public/gorun:$(RELEASE)

service:
	gcloud beta run deploy gorun \
		--image=gcr.io/cloudylabs-public/gorun:$(RELEASE) \
		--region=us-central1 \
		--platform=managed

deploy:
	gcloud beta run deploy gorun \
		--image gcr.io/cloudylabs-public/gorun:$(RELEASE) \
		--platform gke \
		--cluster cr \
		--cluster-location us-east1

apply:
	kubectl apply -f service.yaml -n demo

everything: build image service deploy apply

comp:
	# 1 core:
	curl -s https://gorun.demo.knative.tech/cores/1/concurrency/10/calcs/1000000000 \
		| jq -c ".duration"
	curl -s http://gorun.default.knative.tech/cores/1/concurrency/10/calcs/1000000000 \
		| jq -c ".duration"
	curl -s http://gorun-2gtouos2pq-uc.a.run.app/cores/1/concurrency/10/calcs/1000000000 \
		| jq -c ".duration"

	# All available cores:
	curl -s https://gorun.demo.knative.tech/cores/4/concurrency/10/calcs/1000000000 \
		| jq -c ".duration"
	curl -s http://gorun.default.knative.tech/cores/4/concurrency/10/calcs/1000000000 \
		| jq -c ".duration"
	curl -s http://gorun-2gtouos2pq-uc.a.run.app/cores/8/concurrency/10/calcs/1000000000 \
		| jq -c ".duration"
