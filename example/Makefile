
build:
	docker build -t e2e-runner bases/e2e-runner

deploy-%:
	kubectl delete jobs e2e-runner -n $* --ignore-not-found=true
	kubectl apply -k ./overlays/$*/ --prune -l environment=$*
