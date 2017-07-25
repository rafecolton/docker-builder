go-dockerclient-sort
====================

[![Build Status](https://drone.io/github.com/rafecolton/go-dockerclient-sort/status.png)](https://drone.io/github.com/rafecolton/go-dockerclient-sort/latest)
[![Build Status](https://travis-ci.org/rafecolton/go-dockerclient-sort.svg?branch=master)](https://travis-ci.org/rafecolton/go-dockerclient-sort)
[![GoDoc](https://godoc.org/github.com/rafecolton/go-dockerclient-sort?status.png)](https://godoc.org/github.com/rafecolton/go-dockerclient-sort)
[![Coverage Status](https://img.shields.io/coveralls/rafecolton/go-dockerclient-sort.svg)](https://coveralls.io/r/rafecolton/go-dockerclient-sort?branch=master)

For sorting the results of Docker API calls made using https://github.com/fsouza/go-dockerclient

Example usage:

```go
package main

import(
	"fmt"

	"github.com/fsouza/go-dockerclient"
	"github.com/rafecolton/go-dockerclient-sort"
)

func main() {
	var endpoint = "unix:///var/run/docker.sock"
	var client, _ = docker.NewClient(endpoint)

	images, _ := client.ListImages(false)
	sort.Sort(dockersort.ByCreatedDescending(images))
	fmt.Printf("sorted images: %+v\n", images)
}
```
