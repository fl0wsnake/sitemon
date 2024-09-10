export SITEMON_CONFIG=$(shell cat ./sites.json)
export BUCKET='some bucket'

run:
	go run ./main.go

curl:
	bash ./curl.sh

