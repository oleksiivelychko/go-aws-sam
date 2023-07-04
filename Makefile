download-aws-sam:
	wget -O AWS-SAM-CLI.pkg https://github.com/aws/aws-sam-cli/releases/latest/download/aws-sam-cli-macos-arm64.pkg

awsSamCliSha256Sum := 70a5583160398391cdf0dd5d946448bc36c078d72465ac7c095ad1f56190c707
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

install-aws-cli:
	curl "https://awscli.amazonaws.com/AWSCLIV2.pkg" -o "AWSCLIV2.pkg"
	sudo -S installer -pkg AWSCLIV2.pkg -target /
	rm AWSCLIV2.pkg
	aws --version

uninstall-aws-cli:
	which aws
	ls -l /usr/local/bin/aws
	sudo -S rm /usr/local/bin/aws
	sudo -S rm /usr/local/bin/aws_completer
	sudo -S rm -rf /usr/local/aws-cli
