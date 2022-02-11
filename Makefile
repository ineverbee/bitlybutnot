in-memory:
	docker-compose build
	docker-compose up

postgres:
	docker-compose -f docker-compose.postgres.yml build 
	docker-compose -f docker-compose.postgres.yml up 
