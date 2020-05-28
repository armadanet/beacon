IMAGE_NAME := beacon
CONTAINER_NAME := beacon
OVERLAY_NAME := overlay_beacon

run: build
	docker run -it --name $(CONTAINER_NAME) -p 9898:9898 -v /var/run/docker.sock:/var/run/docker.sock $(IMAGE_NAME)
	#docker run -it --name $(CONTAINER_NAME) -p 9898:9898 -p 8787:8787 -v /var/run/docker.sock:/var/run/docker.sock $(IMAGE_NAME)


build:
	cd .. && docker build -t $(IMAGE_NAME) -f build/Dockerfile .

clean:
	docker rm $(CONTAINER_NAME)
	docker rmi $(IMAGE_NAME)
	docker swarm leave -f
