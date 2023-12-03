GCP start https://developer.hashicorp.com/terraform/tutorials/kubernetes/gke

TF https://registry.terraform.io/providers/hashicorp/google/5.7.0
GIT https://github.com/hashicorp/terraform-provider-google

https://registry.terraform.io/modules/terraform-google-modules/kubernetes-engine/google/latest

https://learnk8s.io/terraform-gke

https://kubernetes.github.io/ingress-nginx/deploy/

```bash

brew install --cask google-cloud-sdk


gcloud init
gcloud components install gke-gcloud-auth-plugin
gcloud auth application-default login

gcloud config set project bikeareaui

## enable API services ??
# gcloud services enable pubsub.googleapis.com
gcloud services enable compute.googleapis.com
gcloud services enable container.googleapis.com

## docker images

gcloud artifacts repositories create btmt \
     --repository-format=docker \
     --location=europe-west4 \
     --immutable-tags \
     --async

gcloud auth configure-docker europe-west4-docker.pkg.dev

gcloud artifacts repositories describe  btmt --location europe-west4

##  create public IP for our ingress
gcloud compute addresses create ingress --global
gcloud compute addresses list



docker build \
     -t europe-west4-docker.pkg.dev/bikeareaui/btmt/auth:latest \
     -f .deploy/Dockerfile \
     --build-arg MAIN_FILE=svc-auth/cmd/main.go \
     ./

docker build \
     -t europe-west4-docker.pkg.dev/bikeareaui/btmt/ui/web:latest \
     -f .deploy/Dockerfile \
     ./


docker push europe-west4-docker.pkg.dev/bikeareaui/btmt/auth:latest
docker push europe-west4-docker.pkg.dev/bikeareaui/btmt/ui/web:latest

helm install auth btmt-app/ -f auth.yaml

helm uninstall auth

terraform init
terraform plan -var-file=gcp.tfvars
terraform apply -var-file=gcp.tfvars

## connect to K8S
gcloud container clusters list
gcloud container clusters get-credentials  btmt --zone europe-west4

kubectl create clusterrolebinding cluster-admin-binding \
  --clusterrole cluster-admin \
  --user $(gcloud config get-value account)


kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/controller-v1.8.2/deploy/static/provider/cloud/deploy.yaml

```
