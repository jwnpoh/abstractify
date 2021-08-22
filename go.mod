module github.com/jwnpoh/abstractify

go 1.16

replace github.com/jwnpoh/abstractify/app => ./app

replace github.com/jwnpoh/abstractify/server => ./server

replace github.com/jwnpoh/abstractify/storage => ./storage

require (
	cloud.google.com/go/storage v1.16.0
	github.com/jwnpoh/abstractify/app v0.0.0-00010101000000-000000000000
	github.com/jwnpoh/abstractify/storage v0.0.0-00010101000000-000000000000
	google.golang.org/api v0.54.0
)
