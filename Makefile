include .env

prepare-zip:
	rm -f go-lambda*
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 \
		go build -o go-lambda \
		-ldflags="-X 'main.awsRegion=$(AWS_REGION)' -X 'main.awsAccessKeyId=$(AWS_ACCESS_KEY)' -X 'main.awsSecretAccessKeyId=$(AWS_SECRET_ACCESS_KEY)'" \
		main.go
	zip go-lambda.zip go-lambda

awsSamCliSum := 1a7e99bfcf898a8dfff7032a729ac52c3482461936901fae215347087bf9000e
install-aws-sam: uninstall-aws-sam
	rm AWS-SAM-CLI.pkg
	wget -O AWS-SAM-CLI.pkg https://github.com/aws/aws-sam-cli/releases/latest/download/aws-sam-cli-macos-arm64.pkg
	@shasum -a 256 AWS-SAM-CLI.pkg | awk '$$1=="1a7e99bfcf898a8dfff7032a729ac52c3482461936901fae215347087bf9000e"{print"Checksum matches!"}'
	sudo -S installer -pkg AWS-SAM-CLI.pkg -target /
	sam --version

uninstall-aws-sam:
	which sam
	ls -l /usr/local/bin/sam
	sudo -S rm /usr/local/bin/sam
	sudo -S rm -rf /usr/local/aws-sam-cli
