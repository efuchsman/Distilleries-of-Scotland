#!/bin/bash

# Build and push Docker image
docker build -t efuchsman/distilleries_of_scotland-distilleries_of_scotland:latest .
docker push efuchsman/distilleries_of_scotland-distilleries_of_scotland:latest

# Update the deployment
kubectl set image deployment/distilleries-of-scotland-deployment distilleries-of-scotland-container=efuchsman/distilleries_of_scotland-distilleries_of_scotland:latest

# Delete and apply the deployment
kubectl delete deployment distilleries-of-scotland-deployment
kubectl apply -f kube/distilleries-of-scotland-deployment.yaml

# Wait for the deployment to be fully rolled out
kubectl rollout status deployment/distilleries-of-scotland-deployment

# Set timeout for waiting
timeout_seconds=300
start_time=$(date +%s)

# Wait for any pod to be running with a timeout
while true; do
    current_time=$(date +%s)
    elapsed_time=$((current_time - start_time))

    if [ $elapsed_time -ge $timeout_seconds ]; then
        echo "Timeout reached. Pod did not start within $timeout_seconds seconds."
        exit 1
    fi

    pod_name=$(kubectl get pods --field-selector=status.phase=Running -o jsonpath='{.items[0].metadata.name}')

    if [ -n "$pod_name" ]; then
        echo "Pod is running. Performing port forwarding to 8000..."
        kubectl port-forward service/distilleries-of-scotland-service 8000:8000
        break
    else
        echo "Pod is not running. Waiting for it to start..."
        sleep 5
    fi
done

