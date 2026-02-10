export BUCKET='some bucket'
export ACCESS_TOKEN='ya29.a0AcM612wm6zU_LeDFrkmdYbKYVXqYyYqxxrugoFwpsdRtQS4cTZi9wAUnhhOpLUv8sojqiKskXWTheeNRj4HkYdvk8d6wnwi5vglWRwyeLatS4Qpsrsqad8Wumkcw1-8gl1snky-YTJ1Cg_ox_6gvsRkZ7R62evlwLkeA-AH6gDDoj9oaCgYKAYISARASFQHGX2MifyS_aXg1xhwXNfsry7svFQ0182'
# GOOGLE_APPLICATION_CREDENTIALS

setup:
	gcloud auth login
	gcloud config set project personal-1
	gcloud projects list
	gcloud storage buckets list
	gcloud storage buckets create gs://sitemon --location=us-east1
	gcloud iam service-accounts list
	gcloud iam service-accounts keys create editor_keys.json --iam-account='personal-437021@appspot.gserviceaccount.com'
	gcloud beta monitoring channels create \
    --display-name="My Email Alert" \
    --type=email \
    --channel-labels=email_address="lastthursdayist@gmail.com" \
# Created notification channel [projects/personal-437021/notificationChannels/9989632443993401243].
	gcloud alpha monitoring policies create --policy-from-file="alert-policy.json"
	# TODO add trigger


run:
	go run ./main.go

curl:
	bash ./curl.sh

