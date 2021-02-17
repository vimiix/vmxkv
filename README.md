# vmxkv

> **[WIP Project]**

### 项目介绍

[[文档]](https://vimiix.com/post/83/) | [[视频]](https://www.bilibili.com/video/bv1vy4y1Y71X)

## Run 

```
make build
./dist/vmxkv -c config-dev.json
```

## Benchmark

```
go test -bench=. ./server
goos: darwin
goarch: amd64
pkg: github.com/vimiix/vmxkv/server
BenchmarkBPTree_InsertWith3Degree-8               	 1000000	      1319 ns/op
BenchmarkBPTree_InsertWith4Degree-8               	 2153112	       568 ns/op
BenchmarkBPTree_InsertWith6Degree-8               	 2647177	       457 ns/op
BenchmarkBPTree_FindWith3Degree1000Elements-8     	15752443	        71.3 ns/op
BenchmarkBPTree_FindWith4Degree1000Elements-8     	35257952	        33.0 ns/op
BenchmarkBPTree_FindWith6Degree1000Elements-8     	30432619	        38.1 ns/op
BenchmarkBPTree_FindWith3Degree10000Elements-8    	11167942	       103 ns/op
BenchmarkBPTree_FindWith4Degree10000Elements-8    	15770590	        74.5 ns/op
BenchmarkBPTree_FindWith6Degree10000Elements-8    	18535426	        62.1 ns/op
BenchmarkBPTree_FindWith3Degree100000Elements-8   	 6328099	       192 ns/op
BenchmarkBPTree_FindWith4Degree100000Elements-8   	10498346	        98.7 ns/op
BenchmarkBPTree_FindWith6Degree100000Elements-8   	17056522	        71.3 ns/op
PASS
ok  	github.com/vimiix/vmxkv/server	17.367s
```

## TODO 

- HTTP RESTFul server
- vmxkv-cli
- docker image