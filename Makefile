# https://github.com/aws/aws-sam-cli/releases/
awsSamCliSha256Sum := 401480ad9ccf3bdf4298b7111e22babbb4fafaf336f60eabdd1735dd1a43fce7

download-aws-sam:
	wget -O AWS-SAM-CLI.pkg https://github.com/aws/aws-sam-cli/releases/latest/download/aws-sam-cli-macos-arm64.pkg

install-aws-sam:
	$(eval result=$(shell shasum -a 256 AWS-SAM-CLI.pkg |\
		awk '$$1=="$(awsSamCliSha256Sum)" {print "Checksum match.";} $$1!="$(awsSamCliSha256Sum)" {print "Checksum mismatch!";}'))
	@echo ${result}
	@if [ "${result}" = "Checksum mismatch!" ]; then\
    	exit;\
    fi;\
	sudo -S installer -pkg AWS-SAM-CLI.pkg -target /
	rm AWS-SAM-CLI.pkg
	sam --version

uninstall-aws-sam:
	which sam
	ls -l /usr/local/bin/sam
	sudo -S rm /usr/local/bin/sam
	sudo -S rm -rf /usr/local/aws-sam-cli

localstack:
	docker run --rm -p 4566:4566 -p 4510-4559:4510-4559 -v /var/run/docker.sock:/var/run/docker.sock --name localstack localstack/localstack

include .env
build-lambda:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 GO111MODULE=auto \
		go build -C lambda/put-message \
		-ldflags "-X main.awsRegion=$(AWS_REGION) -X main.awsAccessKeyID=$(AWS_ACCESS_KEY_ID) -X main.awsSecretAccessKey=$(AWS_SECRET_ACCESS_KEY) -X main.awsEndpoint=$(AWS_ENDPOINT)" \
		-o handler-bin
	zip lambda/put-message/put-message.zip lambda/put-message/handler-bin
