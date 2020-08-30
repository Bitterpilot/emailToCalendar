build:
	docker build -t emailtocal:latest .
	docker image save emailtocal:latest --output emailtocal_latest.tar