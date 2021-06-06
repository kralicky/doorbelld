all: doorbelld

.PHONY: doorbelld
doorbelld:
	CGO_ENABLED=0 go build .

.PHONY: docker-build
docker-build:
	DOCKER_BUILDKIT=1 docker build -t joekralicky/doorbelld .
