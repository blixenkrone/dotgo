.PHONY: docker_build,docker_run,docker_run_it

tag=latest

build:
	go build -o ./bin/dotgo ./cmd/dotgo/dotgo.go

reset:
	rm ./bin/dotgo

all: reset build

docker_build:
	echo "building dotgo:${tag}" \
	&& docker build --pull --rm -f "Dockerfile.test" -t dotgo:${tag} "."

docker_run:
	echo "running container dotgo:${tag}" \
	&& docker run --rm -d dotgo:${tag}

docker_run_it:
	echo "running it container dotgo:${tag}" \
	&& docker run --rm -it dotgo:${tag}
