wrk -t12 -c1000 -d120s http://localhost:5001/ping

docker build -t aralvesandrade/hello-api .
docker push aralvesandrade/hello-api

kind create cluster
kind get clusters
kubectl config get-clusters

kubectl apply -f k8s/deployment.yaml -f k8s/service.yaml -f k8s/hpa.yaml -f k8s/metrics-server.yaml
ou
kubectl apply -f k8s
kubectl get services
kubectl port-forward svc/hello-api-service 5001:80

kubectl run -it fortio --rm --image=fortio/fortio -- load -qps 800 -t 120s -c 70 "http://hello-api-service/ping"

watch -n1 kubectl get pods
watch -n1 kubectl get hpa

kubectl apply -f k8s/dashboard.yaml -f k8s/role-dash.yaml
kubectl proxy --port=8002
http://localhost:8002/api/v1/namespaces/kubernetes-dashboard/services/https:kubernetes-dashboard:/proxy/

kubectl port-forward svc/prometheus-server 9091:80


