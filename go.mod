module github.com/jwnpoh/abstractify

go 1.16

replace github.com/jwnpoh/abstractify/app => ./app

replace github.com/jwnpoh/abstractify/server => ./server

replace github.com/jwnpoh/abstractify/storage => ./storage

replace github.com/jwnpoh/abstractify/logger => ./logger

require (
	github.com/jwnpoh/abstractify/app v0.0.0-00010101000000-000000000000
	github.com/jwnpoh/abstractify/logger v0.0.0-00010101000000-000000000000
	github.com/jwnpoh/abstractify/storage v0.0.0-20210826170114-49c1d8fc7aea
)
