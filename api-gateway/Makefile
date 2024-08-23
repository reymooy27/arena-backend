docker-build:
	docker build -t api-gateway:latest .

docker-push:
	docker tag api-gateway itzyaboirey/api-gateway:latest && docker push itzyaboirey/api-gateway

gen-proto:
	protoc --go_out=. --go-grpc_out=. --go_opt=paths=source_relative --go-grpc_opt=paths=source_relative proto/${file}/${file}.proto

deploy:
	kubectl apply -f k8s/configmap.yaml && kubectl apply -f k8s/deployment.yaml
