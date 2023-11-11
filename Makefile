include .env

# https://github.com/aws/aws-sam-cli/releases/
awsSamCliSha256Sum := 604dca05d13dac09e8343b8099c19fe844d74cda6ff065e622717f09f7cd59de

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

build-put-message:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 GO111MODULE=auto \
		go build -C lambda/put-message \
		-ldflags "-X main.awsRegion=$(AWS_REGION) -X main.awsAccessKeyID=$(AWS_ACCESS_KEY_ID) -X main.awsSecretAccessKey=$(AWS_SECRET_ACCESS_KEY) -X main.awsEndpoint=$(AWS_ENDPOINT)" \
		-o handler-bin
	zip -jrm lambda/put-message/put-message.zip lambda/put-message/handler-bin

build-pop-message:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 GO111MODULE=auto \
		go build -C lambda/pop-message \
		-o handler-bin
	zip -jrm lambda/pop-message/pop-message.zip lambda/pop-message/handler-bin
