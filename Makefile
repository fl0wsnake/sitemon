export SITEMON_CONFIG=$(shell cat ./config.json)
run:
	go run ./main.go
