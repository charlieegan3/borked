PROJECT := borked
TAG := $(shell tar -cf - . | md5sum | cut -f 1 -d " ")

build:
	docker build -t charlieegan3/$(PROJECT):latest -t charlieegan3/$(PROJECT):${TAG} .
	docker build -f Dockerfile.arm -t charlieegan3/$(PROJECT):arm-latest -t charlieegan3/$(PROJECT):arm-${TAG} .

push: build
	docker push charlieegan3/$(PROJECT):latest
	docker push charlieegan3/$(PROJECT):arm-latest
	docker push charlieegan3/$(PROJECT):${TAG}
	docker push charlieegan3/$(PROJECT):arm-${TAG}
