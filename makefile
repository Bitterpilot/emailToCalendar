IMAGE_NAME = docker.pkg.github.com/bitterpilot/emailtocalendar/emailtocalendar
DISK_VOLUME = '/Users/nathan/Documents/Server Side/emailToCalendar/docker_vol:/config'

no-cache:
	docker build --no-cache -t ${IMAGE_NAME}:test .
build:
	docker build -t ${IMAGE_NAME}:test .
run:
	docker run --volume ${DISK_VOLUME} -it ${IMAGE_NAME}