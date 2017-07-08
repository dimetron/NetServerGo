**Installation**

```
go get gopkg.in/mgo.v2

```


**Start mongo**

```
docker run --restart=always --name mongo -m512m -p 27017:27017 -d mongo
```


**Build**

```
go build NetServer.go
go build InjectData.go
```


**Test**

```
./NetServer
./InjectData

echo -n "test out the server" | nc localhost 3333

```