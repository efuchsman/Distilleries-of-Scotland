DOCKER_IMAGE_NAME="efuchsman/distilleries_of_scotland-distilleries_of_scotland"
DOCKER_IMAGE_TAG="latest"

docker build -t "${DOCKER_IMAGE_NAME}:${DOCKER_IMAGE_TAG}" .

kubectl apply -f kube/pg-pvc.yaml
kubectl apply -f kube/pg-deployment.yaml
kubectl apply -f kube/pg-service.yaml
kubectl apply -f kube/distilleries-of-scotland-deployment.yaml
kubectl apply -f kube/distilleries-of-scotland-service.yaml
