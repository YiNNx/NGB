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

通知包括私信，@，关注人发帖，被评论

notification:

- nid
- uid
- time
- type
- pid/cid
- talker
- content







    {
        "id": 126665964778,
        "users": [
            {
                "mid": 32895860,
                "fans": 0,
                "nickname": "进击的胡同学",
                "avatar": "http://i2.hdslb.com/bfs/face/8aede0d719ef1e7a44e2e1a103a46cb569a81366.jpg",
                "mid_link": "",
                "follow": false
            }
        ],
        "item": {
            "item_id": 111079049360,
            "pid": 0,
            "type": "reply",
            "business": "评论",
            "business_id": 0,
            "reply_business_id": 0,
            "like_business_id": 0,
            "title": "[酸了]",
            "desc": "",
            "image": "",
            "uri": "https://www.bilibili.com/video/BV1ir4y1H74w",
            "detail_name": "",
            "native_uri": "bilibili://comment/detail/1/768280609/111079049360/?subType=0\u0026anchor=111079049360\u0026showEnter=1\u0026extraIntentId=0\u0026scene=1\u0026enterUri=bilibili://video/768280609",
            "ctime": 1651239068
        },
        "counts": 1,
        "like_time": 1651310164,
        "notice_state": 0
    },
