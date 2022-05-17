# NGB

**National Geography of Bingyan！**

目标是一个方便配置开箱即用的论坛系统

- 实现基本的用户操作

  - 发帖、关注、申请管理

- 帖子
  - 有标签、节点
  - 可收藏、点赞、评论 & 二级评论

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

- logrus日志

- 提供一个websocket接口实现实时聊天

  > 接口： `ws://localhost:8080/chat?with=<uid>` 
  >
  > 和其他api一样需要身份验证
  >
  > ```
  > Authorization: Bearer TOKEN
  > ```
  >
  > send:
  >
  > ```
  > {
  > 	"content":"xxxxx"
  > }
  > ```
  >
  > 对方若也建立ws连接，将实时Receive:
  >
  > ```
  > {
  >     "Mid": ,
  >     "Time": "xxxxx",
  >     "Content": "xxxxx"
  > }
  > ```
  >

- 异步邮件
  - 开启10个发送协程
  - 记录日志，将发送结果反馈给管理员

- SEO:

  - 使用 [rendora](https://github.com/rendora/rendora)

- elasticsearch

