# NGB

National Geography of Bingyan！

目标是一个方便配置开箱即用的论坛系统

- 用户：自己设计用户基本信息，实现基本的用户操作，比如修改自己的信息以及修改密码。用户能发帖、收藏帖子、关注用户、发布动态。可以查看其他用户的信息，发帖、动态等等。

- 帖子：自己设计基本信息，要求帖子有标签、节点（相当于贴吧的XX吧，每个帖子都应该属于某一个），帖子可以被别人收藏、点赞、评论。评论可以回复，同时可以知道是不是贴主的回复or评论。

- 节点：用户可以选择节点来查看里面的帖子，同样的，也可以通过标签。

- 配置：论坛后端跑起来的基本的配置都需要写在配置文件里，开箱即用可不能让人家自己再改代码编译一遍

- [Optional] 指目前可以不做，但反正之后的阶段会让你写

  - 管理部分。超级管理员，可以管理节点，帖子，用户。节点有自己的管理员。管理员可以管理节点和帖子。

  - 搜索功能

  - 通知系统，用户可以查看自己的回复和对应的帖子，同时被回复时或被私信会收到通知。用户当然可以在自己的首页看到关注的人的动态，同时关注的人发帖时也会被提醒。管理员当然可以群发通知。

  - 用户可以向超级管理员申请创建节点

# API

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


- #### 查看板块

  `GET /board/:bid`

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
          "rgb"
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

  `GET /post/all`

- #### 查看Tag

  `GET /post?tag=xxxxx`


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
  	"statu":true
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
              "title": "安利一款rgb游戏",
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

  （字段均不能为空，且邮箱须有效格式）

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

# 数据库

```

                                       数据表 "public.users"
    栏位     |           类型           | Collation | Nullable |              Default
-------------+--------------------------+-----------+----------+------------------------------------
 uid         | bigint                   |           | not null | nextval('users_uid_seq'::regclass)
 email       | text                     |           | not null |
 username    | text                     |           | not null |
 phone       | text                     |           |          |
 pwd_hash    | text                     |           | not null |
 role        | boolean                  |           |          | false
 create_time | timestamp with time zone |           |          | now()
 avatar      | text                     |           |          |
 nickname    | text                     |           |          |
 gender      | bigint                   |           |          |
 intro       | text                     |           |          |
索引：
    "users_pkey" PRIMARY KEY, btree (uid)
    "users_email_key" UNIQUE CONSTRAINT, btree (email)
    "users_email_username_phone_key" UNIQUE CONSTRAINT, btree (email, username, phone)
    "users_phone_key" UNIQUE CONSTRAINT, btree (phone)
    "users_username_key" UNIQUE CONSTRAINT, btree (username)



                                     数据表 "public.boards"
  栏位  |           类型           | Collation | Nullable |               Default
--------+--------------------------+-----------+----------+-------------------------------------
 bid    | bigint                   |           | not null | nextval('boards_bid_seq'::regclass)
 name   | text                     |           | not null |
 avatar | text                     |           |          |
 time   | timestamp with time zone |           |          | now()
 intro  | text                     |           |          |
索引：
    "boards_pkey" PRIMARY KEY, btree (bid)



                                     数据表 "public.posts"
  栏位   |           类型           | Collation | Nullable |              Default
---------+--------------------------+-----------+----------+------------------------------------
 pid     | bigint                   |           | not null | nextval('posts_pid_seq'::regclass)
 board   | bigint                   |           | not null |
 time    | timestamp with time zone |           |          | now()
 author  | bigint                   |           | not null |
 tags    | jsonb                    |           |          |
 title   | text                     |           | not null |
 content | text                     |           | not null |
索引：
    "posts_pkey" PRIMARY KEY, btree (pid)



                                       数据表 "public.comments"
    栏位    |           类型           | Collation | Nullable |                Default
------------+--------------------------+-----------+----------+---------------------------------------
 cid        | bigint                   |           | not null | nextval('comments_cid_seq'::regclass)
 post       | bigint                   |           | not null |
 is_author  | boolean                  |           |          |
 parent_cid | bigint                   |           |          |
 time       | timestamp with time zone |           |          | now()
 from       | bigint                   |           | not null |
 to         | bigint                   |           |          |
 content    | text                     |           | not null |
索引：
    "comments_pkey" PRIMARY KEY, btree (cid)



               数据表 "public.likes"
   栏位   |  类型  | Collation | Nullable | Default
----------+--------+-----------+----------+---------
 user_uid | bigint |           |          |
 post_pid | bigint |           |          |



            数据表 "public.collections"
   栏位   |  类型  | Collation | Nullable | Default
----------+--------+-----------+----------+---------
 user_uid | bigint |           |          |
 post_pid | bigint |           |          |



           数据表 "public.join_ships"
 栏位 |  类型  | Collation | Nullable | Default
------+--------+-----------+----------+---------
 uid  | bigint |           |          |
 bid  | bigint |           |          |



          数据表 "public.manage_ships"
 栏位 |  类型  | Collation | Nullable | Default
------+--------+-----------+----------+---------
 bid  | bigint |           |          |
 uid  | bigint |           |          |



            数据表 "public.follow_ships"
   栏位   |  类型  | Collation | Nullable | Default
----------+--------+-----------+----------+---------
 followee | bigint |           |          |
 follower | bigint |           |          |

```

