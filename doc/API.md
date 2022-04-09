# API

#### 基本格式

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

需user权限: 标 `*`

需board_manager权限: 标 `**`

需super_admin权限: 标 `***`

#### **http header：**

```
Authorization: Bearer xxxxxxxxxx
```



- #### 注册用户

  `POST /user`

  Request:

  ```
  {
      "email": "xxxx",
  	"username":"xxxxx",
      "pwd": "xxxxxxxx"
  }
  ```

  （字段均不能为空，且邮箱须有效格式，邮箱用户名不可重复）

  Response：

  ```
  {
        "uid": 23
        "token":"xxxxxxxxxxxx"
  }
  ```

- #### 用户登录

  `GET /user/token?email=xxxxx&pwd=xxxxx`

  Response：

  ```
  {
  	"uid": 23
      "token":"xxxxxxxxxxx"
  }
  ```

- #### 查看板块

  `GET /board/:bid`

  Response:

  ```
  {
  	"name":"xxx",
  	"avatar":"xxxx",
  	"intro":"xxxx",
  	"posts":[
  		{
  			"pid":xx,//帖子id
  			"title":"xx",//标题
  			"uid":xx,//发帖人
  			"comments":xx//评论数
  		},
  		{
  			"pid":xx,
  			"title":"xx",
  			"uid":xx,
  			"comments":xx
  		},
  		....
  	]
  }
  ```

- #### 查看所有板块

  `GET /board/all`

- #### 查看贴子

  `GET /post/:pid`

  Respond:

  ```
  {
  	"board":"xxx",
  
  	"time"
  	"title":"xxx",//标题
  	"uid":xx,//发帖人
  	"time":xx,
  	"content":"xxxxx",
  	"comments":[
  		{
  			.....
  		},
  		....
  	]
  	"likes":
  	"collection"
  }
  ```

- #### 查看所有帖子

  `GET /post/all`

- #### 查看用户信息

  `GET /user/:uid`

  Response：

  ```
  {
  	"nickname":"xxxx",
  	"avatar"
  	"gender"
  	"post":[
  		{
  			......
  		},
  		......
  	],
  	"followers"
  	"following"
  	"likes"
  	"collections"
  	"boards_join"
  }
  ```


- #### 修改用户信息*

  `PUT /user/:uid`

  Request:

  ```
  {
  	"email"
  	"username"
  	"phone"
  	"avatar"
  	"nickname":"xxxx",
  	"gender"
  	"intro"
  }
  ```

  （字段均不能为空，且邮箱须有效格式）

- #### 修改密码*

  `PUT /user/:uid/pwd`

  Request

  ```
  {
      "email": "xxxx",
  
      "pwd_old": "xxxxxxxx",
      "pwd_new": "xxxxxxxx"
  }
  ```

- #### 发帖*

  `POST /post`

  Request:

  ```
  {
  	"bid":11,
  	"title":"xxx",
  	"content":"xxx",
  }
  ```

  Response:

  ```
  {
  	"pid":"xxxxx"
  	"time":"xxxxxx"
  }
  ```

- #### 收藏帖子*

  `PUT /post/:pid/collection`

  Request

  ```
  {
  	"collect":true
  }
  ```

- #### 点赞帖子*

  `PUT /post/:pid/like`

  Request

  ```
  {
  	"like":true
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

  Respond

  ```
  {
  	"cid":xxxxx
  }
  ```

- #### 关注用户*

  `PUT /user/:uid/follow`

  ```
  {
  	"follow":true
  }
  ```

  

#### 需要admin权限的操作：

http header：

```
Authorization: Bearer xxxxxxxxxx
```

token异常返回401，role不为admin返回403

- `GET /user/all`

  查看所有用户信息

  返回:

  ```
  [
          {
              "id": 1,
              "email": "xxxxxx",
              "username": "xxxxxx",
              "createTime": "xxxxxx",
              "role": false
          },
  		......
          {
              "id": 11,
              "email": "xxxxxx",
              "username": "xxxxx",
              "createTime": "xxxxxx",
              "role": false
          }
  ]
  ```

- `DELETE /user/:id`

  删除用户

  http header：

  ```
  Authorization: Bearer xxxxxxxxxx
  ```

  返回:

  ```
  {
  	"success": true,
  	"msg": "",
  	"data": null
  }
  ```
