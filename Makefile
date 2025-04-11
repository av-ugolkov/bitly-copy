.PHONY: run-back
run-back:
	docker compose -p bitly-copy -f backend/deploy/develop.yaml up --build --force-recreate    

.PHONY: run-redis
run-redis:
	docker compose -p bitly-copy -f backend/deploy/develop.yaml up redis --build --force-recreate    

.PHONY: run-front
run-front:
	npm --prefix ./frontend run dev

.PHONY: run-test
run-test:
	cd backend && go test ./... -count=1 && cd ..