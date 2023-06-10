include .env

prepare-zip:
	rm -f go-lambda*
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 \
		go build -o go-lambda \
		-ldflags="-X 'main.awsRegion=$(AWS_REGION)' -X 'main.awsAccessKeyId=$(AWS_ACCESS_KEY)' -X 'main.awsSecretAccessKeyId=$(AWS_SECRET_ACCESS_KEY)'" \
		main.go
	zip go-lambda.zip go-lambda

awsSamCliSha256Sum := 420413bf2e399d0c3e27fab60bfd6af7c09550c654a4c99799c3acd8bf8820c1
install-aws-sam: uninstall-aws-sam
	wget -O AWS-SAM-CLI.pkg https://github.com/aws/aws-sam-cli/releases/latest/download/aws-sam-cli-macos-arm64.pkg
	$(eval result=$(shell shasum -a 256 AWS-SAM-CLI.pkg |\
		awk '$$1=="$(awsSamCliSha256Sum)" {print "Checksum match.";} $$1!="$(awsSamCliSha256Sum)" {print "Checksum mismatch!";}'))
	@echo ${result}
	@if [ "${result}" = "Checksum mismatch!" ]; then\
    	exit;\
    fi;\
    rm AWS-SAM-CLI.pkg
	sudo -S installer -pkg AWS-SAM-CLI.pkg -target /
	sam --version

uninstall-aws-sam:
	rm -f AWS-SAM-CLI.pkg
	which sam
	ls -l /usr/local/bin/sam
	sudo -S rm /usr/local/bin/sam
	sudo -S rm -rf /usr/local/aws-sam-cli

install-aws-cli: uninstall-aws-cli
	curl "https://awscli.amazonaws.com/AWSCLIV2.pkg" -o "AWSCLIV2.pkg"
	sudo -S installer -pkg AWSCLIV2.pkg -target /
	rm AWSCLIV2.pkg
	aws --version

uninstall-aws-cli:
	rm -f AWSCLIV2.pkg
	which aws
	ls -l /usr/local/bin/aws
	sudo -S rm /usr/local/bin/aws
	sudo -S rm /usr/local/bin/aws_completer
	sudo -S rm -rf /usr/local/aws-cli
