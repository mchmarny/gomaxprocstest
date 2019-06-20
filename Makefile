RELEASE=0.1.2

.PHONY: mod
mod:
	go mod tidy
	go mod vendor

.PHONY: image
image: mod
	gcloud builds submit \
		--project cloudylabs-public \
		--tag gcr.io/cloudylabs-public/gomaxprocs:$(RELEASE)

.PHONY: service
service:
	gcloud beta run deploy gomaxprocs \
		--image=gcr.io/cloudylabs-public/gomaxprocs:$(RELEASE) \
		--region=us-central1

.PHONY: deploy
deploy:
	gcloud beta run deploy gomaxprocs \
		--image gcr.io/cloudylabs-public/gomaxprocs:$(RELEASE) \
		--platform gke \
		--cluster cr \
		--cluster-location us-east1

.PHONY: undeploy
undeploy:
	gcloud beta run services delete kadvice
