module jobs

go 1.20

replace utils => ./utils

replace pojo => ./pojo

require (
	github.com/PuerkitoBio/goquery v1.8.1
	github.com/tebeka/selenium v0.9.9
	pojo v0.0.0-00010101000000-000000000000
	utils v0.0.0-00010101000000-000000000000
)

require (
	github.com/andybalholm/cascadia v1.3.1 // indirect
	github.com/blang/semver v3.5.1+incompatible // indirect
	golang.org/x/net v0.7.0 // indirect
)
