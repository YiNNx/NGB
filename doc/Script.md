```
ngb=# \d users
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
 posts       | jsonb                    |           |          |
 followers   | jsonb                    |           |          |
 following   | jsonb                    |           |          |
 likes       | jsonb                    |           |          |
 collections | jsonb                    |           |          |
 boards_join | jsonb                    |           |          |
 boards_mng  | jsonb                    |           |          |
 comments    | jsonb                    |           |          |
索引：
    "users_pk" PRIMARY KEY, btree (uid)
    "users_email_key" UNIQUE CONSTRAINT, btree (email)
    "users_email_username_phone_key" UNIQUE CONSTRAINT, btree (email, username, phone)
    "users_phone_key" UNIQUE CONSTRAINT, btree (phone)
    "users_username_key" UNIQUE CONSTRAINT, btree (username)


ngb=# \d comments
                                      数据表 "public.comments"
   栏位    |           类型           | Collation | Nullable |                Default
-----------+--------------------------+-----------+----------+---------------------------------------
 cid       | bigint                   |           | not null | nextval('comments_cid_seq'::regclass)
 sub_cid   | bigint                   |           |          |
 post      | bigint                   |           | not null |
 time      | timestamp with time zone |           |          | now()
 from      | bigint                   |           | not null |
 to        | bigint                   |           | not null |
 content   | text                     |           | not null |
 is_author | boolean                  |           |          | false
索引：
    "comments_pk" PRIMARY KEY, btree (cid)


ngb=# \d posts
                                       数据表 "public.posts"
    栏位     |           类型           | Collation | Nullable |              Default
-------------+--------------------------+-----------+----------+------------------------------------
 pid         | bigint                   |           | not null | nextval('posts_pid_seq'::regclass)
 board       | bigint                   |           | not null |
 time        | timestamp with time zone |           |          | now()
 author      | bigint                   |           | not null |
 title       | text                     |           | not null |
 content     | text                     |           | not null |
 
 likes       | jsonb                    |           |          |
 collections | jsonb                    |           |          |
 tags        | jsonb                    |           |          |
索引：
    "posts_pk" PRIMARY KEY, btree (pid)


ngb=# \d boards
                                      数据表 "public.boards"
   栏位   |           类型           | Collation | Nullable |               Default
----------+--------------------------+-----------+----------+-------------------------------------
 bid      | bigint                   |           | not null | nextval('boards_bid_seq'::regclass)
 name     | text                     |           | not null |
 avatar   | text                     |           |          |
 time     | timestamp with time zone |           |          | now()
 intro    | text                     |           |          |
 managers | jsonb                    |           |          |
 members  | jsonb                    |           |          |
 posts    | jsonb                    |           |          |
索引：
    "boards_pk" PRIMARY KEY, btree (bid)
```

