build:
	docker build --no-cache -t gcr.io/adindopustaka/backend .

push:
	gcloud docker -- push gcr.io/adindopustaka/backend