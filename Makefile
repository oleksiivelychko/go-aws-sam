# https://github.com/aws/aws-sam-cli/releases/
awsSamCliSha256Sum := 985c673d317f773cdecfb92fce9c658c4ddfc50ee0ad7927b9654ee6985962f0

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
