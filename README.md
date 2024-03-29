wrk -t12 -c1000 -d120s http://localhost:5001/ping

docker build -t aralvesandrade/hello-api .
docker push aralvesandrade/hello-api

kind delete cluster
ou
kind delete clusters $(kind get clusters)

### Criando cluster

kind create cluster
ou
kind create cluster --config=kind/2nodes.yaml
ou
kind create cluster --config=kind/2nodes.yaml --name {name-kind}
ou
kind create cluster --config=k8s/nginx-ingress/kind-config-2nodes.yaml

kind get clusters
kubectl config get-clusters
watch -n1 kubectl top nodes

kubectl get pods
kubectl get pods -A
kubectl get pods -A -o wide

### Aplicando deploy, services, hpa e metricas

kubectl apply -f k8s/deployment.yaml -f k8s/service.yaml -f k8s/hpa.yaml
kubectl apply -f k8s/metrics-server.yaml

kubectl get svc
watch -n1 kubectl get hpa
kubectl port-forward svc/hello-api-server 5001:80

### Fazendo um teste de carga fortio 

kubectl run -it fortio --rm --image=fortio/fortio -- load -qps 800 -t 120s -c 70 "http://hello-api-server/ping"

watch -n1 kubectl get pods -l app=hello-api
watch -n1 kubectl get hpa
watch -n1 kubectl top pods

### Instalando o dashboard do K8s e aplicando permissões para ignorar o uso de token e expor na port 8002

kubectl apply -f k8s/dash/dashboard.yaml -f k8s/dash/cluster-role.yaml
watch -n1 kubectl get pods -n kubernetes-dashboard
kubectl proxy --port=8002
http://localhost:8002/api/v1/namespaces/kubernetes-dashboard/services/https:kubernetes-dashboard:/proxy/

### Instalando o Kube-Prometheus

git clone https://github.com/prometheus-operator/kube-prometheus
cd kube-prometheus
kubectl create -f manifests/setup
kubectl apply -f manifests/

kubectl apply -f k8s/grafana && kubectl rollout restart -n monitoring deployment grafana
watch -n1 kubectl get pods -n monitoring
kubectl port-forward -n monitoring svc/grafana 3000

kubectl apply -f k8s/grafana/ingress.yaml
kubectl get ingress -n monitoring
kubectl describe ingress nginx-ingress -n monitoring
kubectl get ingress nginx-ingress -n monitoring -o jsonpath='{.status.loadBalancer.ingress[0].hostname}'
curl localhost/grafana-k8s

### Criar deploy, expor como LoadBalancer na porta 8080 e escalar 3 replicas

kubectl create deployment hello-server --image=gcr.io/google-samples/hello-app:1.0

kubectl expose deployment hello-server --type LoadBalancer --port 80 --target-port 8080

kubectl scale --replicas=5 deployment hello-server

### Instalando o MetalLB

kubectl apply -f https://raw.githubusercontent.com/metallb/metallb/v0.12.1/manifests/namespace.yaml

kubectl apply -f https://raw.githubusercontent.com/metallb/metallb/v0.12.1/manifests/metallb.yaml
kubectl apply -f https://raw.githubusercontent.com/metallb/metallb/v0.13.7/config/manifests/metallb-native.yaml

docker network inspect -f '{{.IPAM.Config}}' kind

kubectl apply -f k8s/metallb/configmap.yaml
ou
kubectl apply -f k8s/metallb/configmap2.yaml

LB_IP=$(kubectl get svc/hello-api-server -o=jsonpath='{.status.loadBalancer.ingress[0].ip}')
for _ in {1..10}; do
  curl ${LB_IP}
done

LB_IP=$(kubectl get svc/hello-server -o=jsonpath='{.status.loadBalancer.ingress[0].ip}')
for _ in {1..10}; do
  curl ${LB_IP}
done

### Nginx

kubectl create deployment nginx --image=nginx
kubectl expose deployment nginx --name nginx-server --type LoadBalancer --port 80 --target-port 80
ou
kubectl apply -f k8s/nginx/

kubectl delete deployment nginx && kubectl delete svc nginx-server

kubectl port-forward svc/nginx-server 8081:80

### Instalando o Nginx Ingress Controller

kubectl version

kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/controller-v1.10.0/deploy/static/provider/cloud/deploy.yaml
kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/main/deploy/static/provider/kind/deploy.yaml

kubectl get pods -n ingress-nginx
ou
kubectl wait --namespace ingress-nginx --for=condition=ready pod --selector=app.kubernetes.io/component=controller --timeout=90s

kubectl apply -f k8s/ingress.yaml
kubectl get ingress
kubectl describe ingress nginx-ingress
kubectl get ingress nginx-ingress -o jsonpath='{.status.loadBalancer.ingress[0].hostname}'
curl localhost

### Referencias

https://medium.com/groupon-eng/loadbalancer-services-using-kubernetes-in-docker-kind-694b4207575d
https://akyriako.medium.com/load-balancing-with-metallb-in-bare-metal-kubernetes-271aab751fb8
https://medium.com/@JohnxLe/kubernetes-nginx-deployment-using-cli-and-yaml-c517b90af0dc
https://github.com/badtuxx/DescomplicandoKubernetes/tree/main/pt/day-9
