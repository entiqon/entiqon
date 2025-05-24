# Makefile for Entiqon Documentation (Dockerized)

IMAGE_NAME=entiqon-docs

.PHONY: docker serve build clean

docker:
	docker build -t $(IMAGE_NAME) .

serve: docker
	docker run -p 8000:8000 $(IMAGE_NAME)

build: docker
	docker run --rm -v $(PWD)/site:/docs/site $(IMAGE_NAME) mkdocs build --clean

clean:
	rm -rf site/