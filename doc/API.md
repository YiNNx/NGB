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



- #### 注册

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

- #### 登录

  `GET /user/token?email=xxxxx&pwd=xxxxx`

  Response：

  ```
  {
  	"uid": 23
      "token":"xxxxxxxxxxx"
  }
  ```

- #### 查看用户公开资料

  `GET /user/:uid`

  Response：

  ```
  {
  	"username":"xxxxx",
  	"nickname":"xxxx",
  	"avatar"
  	"gender"
  	"posts":[
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


- #### 查看账户信息*

  `GET /user/account`

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

- #### 修改用户信息*

  `PUT /user/account`

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

  `PUT /user/password`

  Request

  ```
  {
      "email": "xxxx",
  
      "pwd_old": "xxxxxxxx",
      "pwd_new": "xxxxxxxx"
  }
  ```

- #### 关注用户*

  `PUT /user/:uid/follow`

  ```
  {
  	"status":true
  }
  ```

  取消关注则status为false

  后面同理

- #### 查看板块

  `GET /board/:bid`

  Response:

  ```
  {
  	"bid":
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

- #### 加入板块

  `PUT /board/:bid/join`

  ```
  {
  	"status":true
  }
  ```
  
- #### 查看贴子

  `GET /post/:pid`

  ```
  {
       "pid": "帖子ID",
       "title": "标题"
       "author": {
           "uid": "发布者ID",
           "avatar": "发布者头像URL",
           "nickname": "发布者昵称",
       },
       "time": "发布时间",
   	"board":{
   		"bid":"板块ID"，
   		"avatar":"板块头像"
   		"name":"板块名称"
   	}
       "tags": ["xxxx","xxxxx"],
       "content": "文本内容",
       "likes_count": 1,
       "is_like": true, 
       "likes":[
       {
           "uid": "发布者用户ID",
           "avatar": "发布者头像URL",
           "nickname": "发布者昵称",
       },
       ...
       ]
       "collections_count": 1,
       "is_collect": true, 
       "comments_count": 1, 
       "comments": [
       {
       	"cid":"xx",
       	"parent_cid":"xx",
       	"from":{
           	"uid": "xxx",
           	"avatar": "xxxx",
           	"nickname": "xxx",
       	},
       	"time":"",
       	"is_author":false,
       	"content":"xxxxx"
       },
       ...
       ] // 评论
  }
  ```
  
- ### 查看所有帖子

  `GET /post/all`

- #### 查看Tag

  `GET /post?tag=xxxxx`

  Response:

  ```
  {
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
  	"statu":true
  }
  ```

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
