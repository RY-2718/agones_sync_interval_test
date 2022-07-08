MINIKUBE_PROFILE := sync-interval-test
AGONES_VERSION := 1.23.0
KUBERNETES_VERSION := 1.22.10

up:
	minikube start -p $(MINIKUBE_PROFILE) --driver virtualbox --kubernetes-version $(KUBERNETES_VERSION)
	minikube profile $(MINIKUBE_PROFILE)
	helm repo add agones https://agones.dev/chart/stable
	helm upgrade --install agones --version v$(AGONES_VERSION) \
		--namespace agones-system \
		--create-namespace agones/agones \
		--set agones.ping.install=false \
		--set agones.featureGates="CustomFasSyncInterval=true" \
		--set agones.allocator.replicas=1

delete:
	minikube -p $(MINIKUBE_PROFILE) delete

test:
	go test -count=1 .