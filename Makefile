run:
	docker-compose up

clean:
	docker-compose down -v

build:
	docker-compose -f docker-compose.build.yml build
	docker push canhlinh/cointracker_backend
	docker push canhlinh/cointracker_dashboard
	docker push canhlinh/cointracker_nginx

run_prod:
	sh "./scripts/deploy.sh"
