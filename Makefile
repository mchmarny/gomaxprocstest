RELEASE=0.1.4

.PHONY: mod
mod:
	go mod tidy
	go mod vendor

.PHONY: run
run:
	go run -v main.go

.PHONY: image
image: mod
	gcloud builds submit \
		--project cloudylabs-public \
		--tag gcr.io/cloudylabs-public/gorun:$(RELEASE)

.PHONY: service
service:
	gcloud beta run deploy gorun \
		--image=gcr.io/cloudylabs-public/gorun:$(RELEASE) \
		--region=us-central1 \
		--platform=managed

.PHONY: deploy
deploy:
	gcloud beta run deploy gorun \
		--image gcr.io/cloudylabs-public/gorun:$(RELEASE) \
		--platform gke \
		--cluster cr \
		--cluster-location us-east1

.PHONY: apply
apply:
	kubectl apply -f service.yaml

.PHONY: deployall
deployall: service deploy apply

