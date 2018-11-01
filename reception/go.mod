module github.com/italolelis/reception

require (
	github.com/go-chi/chi v3.3.3+incompatible
	github.com/go-chi/render v1.0.1
	github.com/go-sql-driver/mysql v1.4.0 // indirect
	github.com/italolelis/coffee-shop v0.0.0-20181025085226-81da13d22b95
	github.com/italolelis/log v0.0.0
	github.com/jmoiron/sqlx v0.0.0-20180614180643-0dae4fefe7c0
	github.com/kelseyhightower/envconfig v1.3.0
	github.com/kr/pretty v0.1.0 // indirect
	github.com/lib/pq v1.0.0
	github.com/mattn/go-sqlite3 v1.10.0 // indirect
	github.com/satori/go.uuid v1.2.0
	golang.org/x/net v0.0.0-20181029044818-c44066c5c816 // indirect
	google.golang.org/appengine v1.2.0 // indirect
	gopkg.in/check.v1 v1.0.0-20180628173108-788fd7840127 // indirect
)

replace github.com/italolelis/log => ../log
