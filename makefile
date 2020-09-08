IMAGE_NAME = docker.pkg.github.com/bitterpilot/emailtocalendar/emailtocalendar
IMAGE_TAG = :test
DISK_VOLUME = '/Users/nathan/Documents/Server Side/emailToCalendar/test_data:/config'
DOCKERFILE = ./dockerfile


no-cache:
	docker build -f ${DOCKERFILE} --no-cache -t ${IMAGE_NAME}${IMAGE_TAG} .
build:
	docker build -f ${DOCKERFILE} -t ${IMAGE_NAME}${IMAGE_TAG} .
run: build
	docker run --volume ${DISK_VOLUME} -it ${IMAGE_NAME}${IMAGE_TAG}