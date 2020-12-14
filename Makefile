clean:
	docker-compose down --remove-orphans

setup:
	docker-compose up --build
