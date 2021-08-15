module github.com/jwnpoh/abstractify/app

go 1.16

replace github.com/jwnpoh/abstractify/storage => ../storage

require (
	github.com/fogleman/gg v1.3.0
	github.com/golang/freetype v0.0.0-20170609003504-e2365dfdc4a0 // indirect
	github.com/jwnpoh/abstractify/storage v0.0.0-00010101000000-000000000000
	golang.org/x/image v0.0.0-20210628002857-a66eb6448b8d // indirect
)
