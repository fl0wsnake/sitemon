export CONFIG=$(shell cat ./config.json)
export BUCKET='some bucket'
export ACCESS_TOKEN='ya29.a0AcM612wm6zU_LeDFrkmdYbKYVXqYyYqxxrugoFwpsdRtQS4cTZi9wAUnhhOpLUv8sojqiKskXWTheeNRj4HkYdvk8d6wnwi5vglWRwyeLatS4Qpsrsqad8Wumkcw1-8gl1snky-YTJ1Cg_ox_6gvsRkZ7R62evlwLkeA-AH6gDDoj9oaCgYKAYISARASFQHGX2MifyS_aXg1xhwXNfsry7svFQ0182'

run:
	go run ./main.go

curl:
	bash ./curl.sh

