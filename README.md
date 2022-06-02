# NGB

**National Geography of Bingyan！**

目标是一个方便配置开箱即用的论坛系统

### Feature

- `main`分支使用jwt模式鉴权

  `cookie-session-version`分支使用cookie-session模式鉴权

- User：基本用户系统、发帖、关注、申请管理。

  Post：有标签、板块；可收藏、点赞、评论 & 二级评论。

  Super Admin：管理板块，帖子，用户。

  Board Admin：管理板块和帖子；用户可申请创建板块，或者申请板块管理员。

  Notification：私信&评论&用户在帖文中被提及&关注人发帖

- websocket实现简易IM功能

- 消息通知系统

  [NGB-notification ](https://github.com/YiNNx/NGB-notification/)

  通过RabbitMQ进行websocket/邮件通知

- 异步邮件

  - 开启10个发送协程
  - 记录日志，将发送结果反馈给管理员

- 搜索

  使用 elasticsearch 对帖子标题和内容进行分词搜索

- SEO

  使用 [rendora](https://github.com/rendora/rendora)

- 日志系统：logrus

  

