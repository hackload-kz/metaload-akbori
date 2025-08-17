.PHONY: build run test clean dev-deps dev-deps-down build-prod rebuild-prod up-prod down-prod restart-prod logs-prod

build: 
	./mvnw clean package -DskipTests

run: 
	./mvnw spring-boot:run -Dspring-boot.run.profiles=dev

test: 
	./mvnw test

clean: 
	./mvnw clean
	docker system prune -f

dev-deps: 
	docker-compose up -d

dev-deps-down: 
	docker-compose down


build-prod: 
	docker compose -f docker-compose-prod.yml build

up-prod: 
	docker compose -f docker-compose-prod.yml up -d

down-prod: 
	docker compose -f docker-compose-prod.yml down

restart-prod: 
	docker compose -f docker-compose-prod.yml restart

logs-prod: 
	docker compose -f docker-compose-prod.yml logs -f --tail=100