docker-build:
	docker build -t arena-service:latest .

docker-push:
	docker tag arena-service itzyaboirey/arena-service:latest && docker push itzyaboirey/arena-service

gen-proto:
	protoc --go_out=. --go-grpc_out=. --go_opt=paths=source_relative --go-grpc_opt=paths=source_relative proto/arena.proto

gen-migration:
	migrate create -ext sql -dir db/migrations -seq ${title}

deploy:
	kubectl apply -f k8s/configmap.yaml && kubectl apply -f k8s/secret.yaml && kubectl apply -f k8s/deployment.yaml
