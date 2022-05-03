# NGB

National Geography of Bingyan！

目标是一个方便配置开箱即用的论坛系统

#### Level 1

- 实现基本的用户操作
  - 修改用户信息 & 密码
  - 发帖、关注

- 帖子
  - 有标签、节点
  - 可收藏、点赞、评论 & 二级评论

- 节点

#### Level 2

- 超级管理员
  - 管理节点，帖子，用户

- 节点管理员
  - 管理节点和帖子
  - 用户可以申请创建节点，或者申请成为节点管理员

- 通知功能
  - 私信
  - 评论
  - 用户在帖文中被@
  - 关注人发帖

- 用logrus实现简单的日志系统



### Level2新增API：

### super_admin

- #### 创建板块✔

  `POST /board`

  ```
  {
  	"name": "",
      "avatar": "",
      "intro": "",
  }
  ```

- #### 编辑板块信息✔

  `POST /board/:bid?`

  ```
  {
  	"name": "",
      "avatar": "",
      "intro": "",
  }
  ```

- #### 删除板块✔

  `DELETE /board/:bid?`

- #### 删帖✔

  `DELETE /post/:pid`

- #### 删用户✔

  `DELETE /user/:uid`

- #### 查看所有用户✔

  `GET /user/all`

- #### 查看所有板块管理员✔

  `GET /user/admin`

- #### 查看管理员 & 创建板块申请✔

  `GET /apply/admin`

  `GET /apply/board`

- #### 通过申请✔

  `POST /apply/:apid`

  ```
  {
  	"status":true,
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

  
