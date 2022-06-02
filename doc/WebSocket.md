## WebSocket

`ws://localhost:8080/chat?with=<uid>` 

和其他api一样需要jwt

```
Authorization: Bearer TOKEN
```

send:

```
{
	"content":"xxxxx"
}
```

对方若也建立ws连接，将实时Receive:

```
{
 "mid": ,
 "time": "xxxxx",
 "content": "xxxxx"
}
```
