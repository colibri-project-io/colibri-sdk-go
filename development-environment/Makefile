STACK_NAME=dev

start:
	docker-compose -p ${STACK_NAME} up -d --remove-orphans

stop:
	docker-compose -p ${STACK_NAME} stop

restart: stop start

clean:
	docker-compose -p ${STACK_NAME} down -v

logs:
	docker-compose -p ${STACK_NAME} logs -f
