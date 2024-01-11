# Start Minikube (if not already running)
minikube start --profile distilleries-of-scotland

# Set the Minikube context
kubectl config use-context distilleries-of-scotland

DOCKER_IMAGE_NAME="efuchsman/distilleries_of_scotland-distilleries_of_scotland"
DOCKER_IMAGE_TAG="latest"

docker build -t "${DOCKER_IMAGE_NAME}:${DOCKER_IMAGE_TAG}" .
docker push efuchsman/distilleries_of_scotland-distilleries_of_scotland:latest

kubectl set image deployment/distilleries-of-scotland-deployment distilleries-of-scotland-container=efuchsman/distilleries_of_scotland-distilleries_of_scotland:latest

kubectl apply -f kube/pg-pvc.yaml
kubectl apply -f kube/pg-deployment.yaml
kubectl apply -f kube/pg-service.yaml
kubectl apply -f kube/distilleries-of-scotland-deployment.yaml
kubectl apply -f kube/distilleries-of-scotland-service.yaml
