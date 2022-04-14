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

