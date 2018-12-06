module github.com/italolelis/reception

require (
	contrib.go.opencensus.io/integrations/ocsql v0.1.2
	github.com/go-chi/chi v3.3.3+incompatible
	github.com/go-chi/render v1.0.1
	github.com/go-sql-driver/mysql v1.4.0 // indirect
	github.com/gogo/protobuf v1.1.1 // indirect
	github.com/golang/protobuf v1.2.0
	github.com/italolelis/kit v0.0.0
	github.com/jmoiron/sqlx v0.0.0-20180614180643-0dae4fefe7c0
	github.com/kelseyhightower/envconfig v1.3.0
	github.com/kr/pretty v0.1.0 // indirect
	github.com/lib/pq v1.0.0
	github.com/mattn/go-sqlite3 v1.10.0 // indirect
	github.com/rafaeljesus/rabbus v2.2.1+incompatible
	github.com/satori/go.uuid v1.2.0
	go.opencensus.io v0.18.0
)

replace github.com/italolelis/kit => ../kit
