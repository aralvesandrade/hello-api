wrk -t12 -c1000 -d120s http://localhost:5001/ping

docker build -t aralvesandrade/hello-api .
docker push aralvesandrade/hello-api

kind delete cluster

# Criando cluster

kind create cluster
ou
kind create cluster --config=kind/config.yaml
ou
kind create cluster --config=kind/config.yaml --name {name-kind}

kind get clusters
kubectl config get-clusters

# Aplicando deploy, services, hpa e metricas

kubectl apply -f k8s/deployment.yaml -f k8s/service.yaml -f k8s/hpa.yaml -f k8s/metrics-server.yaml
ou
kubectl apply -f k8s/

kubectl get services
kubectl port-forward svc/hello-api-server 5001:80

# Fazendo um teste de carga fortio 

kubectl run -it fortio --rm --image=fortio/fortio -- load -qps 800 -t 120s -c 70 "http://hello-api-server/ping"

watch -n1 kubectl get pods -l app=hello-api
watch -n1 kubectl get hpa

# Instalando o dashboard do K8s e aplicando permissões para ignorar o uso de token e expor na port 8002

kubectl apply -f k8s/dash/dashboard.yaml -f k8s/dash/cluster-role.yaml
kubectl proxy --port=8002
http://localhost:8002/api/v1/namespaces/kubernetes-dashboard/services/https:kubernetes-dashboard:/proxy/

# Instalando o Kube-Prometheus

git clone https://github.com/prometheus-operator/kube-prometheus
cd kube-prometheus
kubectl create -f manifests/setup
kubectl apply -f manifests/

kubectl apply -f k8s/grafana && kubectl rollout restart -n monitoring deployment grafana
kubectl get pods -n monitoring
kubectl port-forward -n monitoring svc/grafana 3000

# Criar deploy, expor como LoadBalancer na porta 8080 e escalar 3 replicas

kubectl create deployment hello-server --image=gcr.io/google-samples/hello-app:1.0

kubectl expose deployment hello-server --type LoadBalancer --port 80 --target-port 8080

kubectl scale --replicas=5 deployment hello-server

# Instalando o MetalLB

kubectl apply -f https://raw.githubusercontent.com/metallb/metallb/v0.12.1/manifests/namespace.yaml

kubectl apply -f https://raw.githubusercontent.com/metallb/metallb/v0.12.1/manifests/metallb.yaml

docker network inspect -f '{{.IPAM.Config}}' kind

kubectl apply -f k8s/metallb/configmap.yaml

LB_IP=$(kubectl get svc/hello-api-server -o=jsonpath='{.status.loadBalancer.ingress[0].ip}')
for _ in {1..10}; do
  curl ${LB_IP}
done

LB_IP=$(kubectl get svc/hello-server -o=jsonpath='{.status.loadBalancer.ingress[0].ip}')
for _ in {1..10}; do
  curl ${LB_IP}
done

# Referencias

https://medium.com/groupon-eng/loadbalancer-services-using-kubernetes-in-docker-kind-694b4207575d
https://akyriako.medium.com/load-balancing-with-metallb-in-bare-metal-kubernetes-271aab751fb8
