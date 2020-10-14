build_dir=docdb-docker-compose
tar_file=${build_dir}.tar.gz

export addr=${DOC_DB_EC2_ADDR}

PHONY: install
install:
	@echo "building docker compose artifact ${build_dir}"
	./poc/build.sh ${build_dir}

	@echo "scp artifact on to EC2 instance: ${addr}"
	scp -i ~/.ssh/ons-web-development.pem ${tar_file} ${addr}:

	@echo "scp install script to EC2 instance:  ${addr}"
	scp -i ~/.ssh/ons-web-development.pem poc/install-poc.sh ${addr}:

	@echo "installation complete"

PHONY: ssh
ssh:
	@echo "SSHing on top EC2 instance: ${addr}"
	ssh -i ~/.ssh/ons-web-development.pem ${addr}