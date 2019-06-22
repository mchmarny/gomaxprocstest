RELEASE=0.1.5

.PHONY: mod
mod:
	go mod tidy
	go mod vendor

debug:
	go run -v main.go

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
