# NGB论坛

## 数据库

```
                                       数据表 "public.users"
    栏位     |           类型           | Nullable |              Default
-------------+--------------------------+----------+------------------------------------
 uid         | bigint                   | not null | nextval('users_uid_seq'::regclass)
 email       | text                     | not null |
 username    | text                     | not null |
 phone       | text                     |          |
 pwd_hash    | text                     | not null |
 role        | boolean                  |          | false
 create_time | timestamp with time zone |          | now()
 avatar      | text                     |          |
 nickname    | text                     |          |
 gender      | bigint                   |          |
 intro       | text                     |          |
 posts       | jsonb                    |          |
 followers   | jsonb                    |          |
 following   | jsonb                    |          |
 likes       | jsonb                    |          |
 collections | jsonb                    |          |
 boards_join | jsonb                    |          |
 boards_mng  | jsonb                    |          |
索引：
    "users_pk" PRIMARY KEY, btree (uid)
    "users_email_key" UNIQUE CONSTRAINT, btree (email)
    "users_email_username_phone_key" UNIQUE CONSTRAINT, btree (email, username, phone)
    "users_phone_key" UNIQUE CONSTRAINT, btree (phone)
    "users_username_key" UNIQUE CONSTRAINT, btree (username)
```

```
                                      数据表 "public.boards"
   栏位   |           类型           | Nullable |               Default
----------+--------------------------+----------+-------------------------------------
 bid      | bigint                   | not null | nextval('boards_bid_seq'::regclass)
 name     | text                     | not null |
 avatar   | text                     |          |
 time     | timestamp with time zone |          | now()
 intro    | text                     |          |
 managers | jsonb                    |          |
 members  | jsonb                    |          |
 posts    | jsonb                    |          |
索引：
    "boards_pk" PRIMARY KEY, btree (bid)
```

```
                                       数据表 "public.posts"
    栏位     |           类型           | Nullable |              Default
-------------+--------------------------+----------+------------------------------------
 pid         | bigint                   | not null | nextval('posts_pid_seq'::regclass)
 board       | bigint                   | not null |
 time        | timestamp with time zone |          | now()
 author      | bigint                   | not null |
 tags  | jsonb                    |          |
 title       | text                     | not null |
 content     | text                     | not null |
 comments    | jsonb                    |          |
 likes       | jsonb                    |          |
 collections | jsonb                    |          |
索引：
    "posts_pk" PRIMARY KEY, btree (pid)
```

```
                                       数据表 "public.comments"
    栏位    |           类型           | Nullable |                Default
------------+--------------------------+----------+---------------------------------------
 cid        | bigint                   | not null | nextval('comments_cid_seq'::regclass)
 parent_cid | bigint                   |          |
 post       | bigint                   | not null |
 time       | timestamp with time zone |          | now()
 from       | bigint                   | not null |
 to         | bigint                   | not null |
 content    | text                     | not null |
索引：
    "comments_pk" PRIMARY KEY, btree (cid)
```