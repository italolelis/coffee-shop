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
	github.com/prometheus/client_golang v0.9.1 // indirect
	github.com/prometheus/common v0.0.0-20181126121408-4724e9255275 // indirect
	github.com/prometheus/procfs v0.0.0-20181126161756-619930b0b471 // indirect
	github.com/rafaeljesus/rabbus v2.2.1+incompatible
	github.com/rafaeljesus/retry-go v0.0.0-20171214204623-5981a380a879 // indirect
	github.com/satori/go.uuid v1.2.0
	github.com/sony/gobreaker v0.0.0-20181109014844-d928aaea92e1 // indirect
	github.com/streadway/amqp v0.0.0-20181107104731-27835f1a64e9 // indirect
	go.opencensus.io v0.18.0
	golang.org/x/net v0.0.0-20181114220301-adae6a3d119a // indirect
	golang.org/x/sync v0.0.0-20181108010431-42b317875d0f // indirect
	google.golang.org/appengine v1.3.0 // indirect
	gopkg.in/check.v1 v1.0.0-20180628173108-788fd7840127 // indirect
)

replace github.com/italolelis/kit => ../kit
