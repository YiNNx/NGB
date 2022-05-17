#  API

- #### 基本格式


```
{
    "success": true,
    "msg": "",
    "data": {
    }
}
```

```
{
	"success": false,
	"msg": "xxxxx",
	"data": null
}
```

- #### 注册

  `POST /user`

  Request:

  ```
  {
      "email": "2333@moe.com",
  	"username": "233333",
      "pwd": "123456"
  }
  ```

  Response：

  ```
  {
          "uid": 5,
          "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6NSwicm9sZSI6ZmFsc2UsImV4cCI6MTY1MDE1NzYxNX0.UwSnB0dCwvusKYvKgfFLKBqzJLt1ZU-KDhDoS2n8r2o"
      }
  ```
  
- #### 登录

  `GET /user/token?email=2333@moe.com&pwd=123456`

  Response：

  ```
  {
          "uid": 5,
          "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6NSwicm9sZSI6ZmFsc2UsImV4cCI6MTY1MDE1Nzc2N30.OlEbOnKDl81Aws6TDapjQaMzXZTCdj5s8Kg4Wy-_UYY"
      }
  ```


- #### 查看板块内容

  `GET /board/:bid?amount=3&page=1`

  分页默认为amount=10，page=1

  Response:

  ```
  {
      "bid": 1,
      "name": "综合版",
      "avatar": "",
      "intro": "综合板块",
      "posts": [
          {
              "pid": 2,
              "title": "今天是星期四",
              "author": 1,
              "time": "2022-04-14T13:23:14.739842+08:00",
              "board": 1,
              "likes_count": 1
          },
          {
              "pid": 4,
              "title": "安利一款rpg游戏",
              "author": 9,
              "time": "2022-04-14T15:34:32.929336+08:00",
              "board": 1,
              "likes_count": 1
          }
      ]
  }
  ```

- #### 查看所有板块

  `GET /board/all`

  Response:

  ```
  [
          {
              "bid": 1,
              "name": "综合版",
              "intro": "综合板块"
          },
          {
              "bid": 2,
              "name": "游戏版",
              "intro": "游戏交流"
          },
          {
              "bid": 3,
              "name": "版务",
              "intro": "论坛管理板块"
          },
          {
              "bid": 4,
              "name": "技术版",
              "intro": "技术交流"
          }
      ]
  ```

- #### 查看贴子

  `GET /post/:pid`

  ```
  {
      "uid": 6,
      "title": "安利一款rpg游戏",
      "author": {
          "uid": 9,
          "username": "meeseeeks",
      },
      "time": "2022-04-14T15:46:34.518722+08:00",
      "board": {
          "bid": 2,
          "name": "游戏版",
          "intro": "游戏交流"
      },
      "tags": [
          "rpg"
      ],
      "content": "安利一款rpg游戏是怎么回事呢？安利相信大家都很熟悉， 但是安利一款rpg游戏是怎么回事呢？下面就让小编带大家一起了解吧。 安利一款rpg游戏，其实就是安利一款rpg游戏了。 大家可能会感到很惊讶，安利怎么会一款rpg游戏呢？... 但事实就是这样，小编也感到非常惊讶。 那么这就是关于安利一款rpg游戏的事情了，大家有什么想法呢？欢迎在评论区告诉小编一起讨论哦",
      "likes_count": 2,
      "comments_count": 1,
      "comments": [
          {
              "cid": 8,
              "parent_cid": 0,
              "is_author": false,
              "from": {
                  "uid": 1,
                  "username": "xxxxx",
              },
              "time": "2022-04-14T15:51:05.970339+08:00",
              "content": "好！"
          }
      ]
  }
  ```
  
- ### 查看所有帖子

  `GET /post/all?amount=3&page=1`

- #### 查看Tag

  `GET /post?tag=xxxxx&amount=3&page=1`


- #### 发帖*

  `POST /post`

  Request:

  ```
  {
  	"title":"今天是星期四",
  	"content":"今天是星期四是怎么回事呢？今天是相信大家都很熟悉， 但是今天是星期四是怎么回事呢？下面就让小编带大家一起了解吧。 今天是星期四，其实就是今天是星期四了。 大家可能会感到很惊讶，今天是怎么会星期四呢？... 但事实就是这样，小编也感到非常惊讶。 那么这就是关于今天是星期四的事情了，大家有什么想法呢？欢迎在评论区告诉小编一起讨论哦",
      "tag":"灌水"
  }
  ```

  Response:

  ```
  {
  	"pid": 2,
  	"time": "2022-04-14T13:23:14.739842+08:00"
  }
  ```

- #### 收藏帖子*

  `PUT /post/:pid/collection`

  Request

  ```
  {
  	"status":true
  }
  ```

  取消收藏则status为false，后面同理

- #### 点赞帖子*

  `PUT /post/:pid/like`

  Request

  ```
  {
  	"status":true
  }
  ```

- #### 评论帖子*

  `POST /post/:pid/comment`

  Request

  ```
  {
  	"content":"xxxxxx"
  }
  ```

- #### 发表子评论*

  `POST /post/:pid/comment/:cid/subcomment`

  Request

  ```
  {
  	"to":12,
  	"content":"xxxxxx"
  }
  ```

- #### 查看用户公开资料

  `GET /user/:uid`

  Response：

  ```
  {
      "username": "233333",
      "nickname": "moe",
      "avatar": "",
      "gender": 0,
      "posts": [
          {
              "pid": 3,
              "title": "今天是星期四",
              "author": 5,
              "time": "2022-04-14T15:29:22.460602+08:00",
              "board": 1,
              "likes_count": 1
          }
      ],
      "followers": null,
      "following": [
          {
              "uid": 9,
              "username": "meeseeeks",
          }
      ],
      "likes": [
          {
              "pid": 1,
              "title": "test",
              "author": 1,
              "time": "2022-04-12T13:20:32.426289+08:00",
              "board": 2,
              "likes_count": 2
          }
      ],
      "collections": [
          {
              "pid": 4,
              "title": "安利一款rpg游戏",
              "author": 9,
              "time": "2022-04-14T15:34:32.929336+08:00",
              "board": 1,
              "likes_count": 1
          }
      ],
      "boards_join": null
  }
  ```


- #### 关注用户*

  `PUT /user/follow/:uid`

  ```
  {
  	"status":true
  }
  ```

- #### 查看账户信息*

  `GET /user/account`

  Response:

  ```
  {
      "email": "2333@moe.com",
      "username": "233333",
      "phone": "1582333",
      "nickname": "moe",
      "gender": 0,
      "intro": "a boring person"
  }
  ```

- #### 修改用户信息*

  `PUT /user/account`

  Request:

  ```
  {
      "email": "2333@moe.com",
      "username": "233333",
      "phone": "1582333",
      "nickname": "moe",
      "gender": 0,
      "intro": "a boring person"
  }
  ```

- #### 修改密码*

  `PUT /user/password`

  Request

  ```
  {
      "email": "2333@moe.com",
      "pwd_old": "123456",
      "pwd_new": "654321"
  }
  ```



## LEVEL 2


### super_admin

- #### 创建板块 

`POST /board`

```
  {
   "name": "",
      "avatar": "",
      "intro": "",
  }
```

- #### 编辑板块信息 

`POST /board/:bid?`

```
  {
   "name": "",
      "avatar": "",
      "intro": "",
  }
```

- #### 删除板块 

`DELETE /board/:bid?`

- #### 删帖 

`DELETE /post/:pid`

- #### 删用户 

`DELETE /user/:uid`

- #### 查看所有用户 

`GET /user/all`

- #### 查看所有板块管理员 

`GET /user/admin`

- #### 查看管理员 & 创建板块申请 

`GET /apply/admin`

`GET /apply/board`

- #### 通过申请 

`POST /apply/:apid`

```
  {
   "status":true,
  }
```

- #### 邮件

  `POST /admin/email`

  ```
  {
  	"to":[
  	"xxxxx@xx.com",
  	"xxxxx@xx.com",
  	...
  	],
  	"subject":"xxxxx",
  	"content":
  }
  ```

  

### admin

- #### 编辑板块信息✔

`POST /board/:bid?`

```
  {
   "name": "",
      "avatar": "",
      "intro": "",
  }
```

- #### 删除帖子✔

`DELETE /post/:pid`

- #### 申请成为板块管理员✔

`POST /apply/admin`

```
  {
   "bid": 123,
   "reason": ""
  }
```

- #### 申请创建板块✔

`POST /apply/board`

```
  {
   "name": "",
   "reason": ""
  }
```

### notification

- #### 查看通知✔

`GET /notification`

- #### 查看未读通知✔

`GET /notification/new`

### message

- #### 私信✔

`POST /message?receiver=<uid>`

```
  {
   "content":""
  }
```



### 邮件

`POST /admin/email`

```
{
	"to":[
	"xxxx@xx.com",
	"xxxx@xx.com",
	...
	],
	c
}
```





