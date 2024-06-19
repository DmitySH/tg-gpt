PUBLIC_IP=dimasik@178.154.220.4
APP_DIR=/home/dimasik/app
RSA_FILE=~/.ssh_yc/id_rsa

.PHONY: up
up:
	docker-compose up --build

.PHONY: tf-apply
tf-apply:
	cd ./deployments/terraform && bash tf.sh apply

.PHONY: tf-destroy
tf-destroy:
	cd ./deployments/terraform && bash tf.sh destroy

.PHONY: push
push:
	scp -r -i $(RSA_FILE) ./* $(PUBLIC_IP):$(APP_DIR)

.PHONY: format
format:
	smartimports .
