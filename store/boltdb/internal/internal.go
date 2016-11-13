package internal

//go:generate protoc --go_out=. item.proto meta.proto
//go:generate protoc-go-inject-tag -input=./item.pb.go
//go:generate protoc-go-inject-tag -input=./meta.pb.go
