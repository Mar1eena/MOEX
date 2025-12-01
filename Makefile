name = trb


# Строка подключения

build:
	docker-compose --project-name=${name} build

up:
	docker-compose --project-name=${name} up -d

upd:
	docker-compose --project-name=${name} up --build -d
	
down:
	docker-compose --project-name=${name} down