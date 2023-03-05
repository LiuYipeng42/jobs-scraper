module jobs

go 1.20

replace utils => ./utils

replace job => ./job

replace task => ./task
require (
	github.com/PuerkitoBio/goquery v1.8.1
	github.com/tebeka/selenium v0.9.9
	job v0.0.0-00010101000000-000000000000
	task v0.0.0-00010101000000-000000000000
	utils v0.0.0-00010101000000-000000000000
)

require (
	github.com/andybalholm/cascadia v1.3.1 // indirect
	github.com/blang/semver v3.5.1+incompatible // indirect
	github.com/go-sql-driver/mysql v1.7.0 // indirect
	github.com/golang/snappy v0.0.1 // indirect
	github.com/jmoiron/sqlx v1.3.5 // indirect
	github.com/nsqio/go-nsq v1.1.0 // indirect
	golang.org/x/net v0.7.0 // indirect
)
