

go:
	go run ./cmd/web/

css:
	npx tailwindcss -o ./ui/static/css/main.css --watch
	
opendb:
	docker compose exec db psql -U acis
	
mc:
	migrate create -seq -ext=.sql -dir=./migrations ${name} 
	
mu: 
	migrate -path=./migrations -database=postgres://acis:acis@localhost/acis?sslmode=disable up

mdrop:
	migrate -path=./migrations -database=postgres://acis:acis@localhost/acis?sslmode=disable drop force