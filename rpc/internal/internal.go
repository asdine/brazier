package internal

//go:generate protoc --go_out=plugins=grpc:. saver.proto getter.proto deleter.proto lister.proto
