module github.com/atomix/kubernetes-tests

go 1.13

require (
	github.com/atomix/go-client v0.1.1
	github.com/onosproject/helmit v0.6.0
	github.com/spf13/cobra v0.0.6
	github.com/stretchr/testify v1.5.1
	k8s.io/api v0.17.3
	k8s.io/apimachinery v0.17.3
)

replace github.com/atomix/api => ../atomix-api

replace github.com/atomix/go-client => ../atomix-go-client

replace github.com/atomix/go-framework => ../atomix-go-node

replace github.com/atomix/go-local => ../atomix-go-local
