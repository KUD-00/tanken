# Architecture

Database: RDS(DynamoDB/Postgres)
- on-demand

Cache: ElastiCache(为什么不ec2跑redis => 为什么要用AWS)
- 通过缓存命中、聚合写入来降低对DynamoDB的读取写入操作
- 设置缓存失效策略、缓存淘汰机制
- 缓存一致性：事件/消息队列，发布/订阅模型？

ElasticSearch:
- 

gateway(in eks?)

read/write data service(EC2 in eks)

frontend service(EC2 in eks)
- nuxt.js
- Server-Side Rendering

client browser
- grpc-web to connect read/write data service?
- 似乎得绕过frontend service?因为nuxt.js不支持在server环境执行函数？

---

## Data Cache
redis cluster使用，但是不知道怎么增加master node的数量，怎么搞都是1master，3replicas

设计理念：
- 建立“查询模式 -> id”的永久性cache
- data-fetcher通过redis获取到相应id然后尝试对应的缓存命中
- 同时向后端数据库请求部分未命中的数据
- 前端收到两次数据，一次是缓存命中直接返回，一次是未命中的数据库数据

大体有三种类型，db.Post, cache.Post, pb(rpc).Post

数据存储，基本只存cache, 统一从cache写回DB，特点：redis新建和更新是一样的
- cacheX(xid, details)

数据的读取，先读cache, 读不到读不满或者有需要从数据库随机加
- 大的getX(xid, details)，里面包括两部分，getXFromCache(), getXFromDB(), 这是保证一定能读到的
- getXFromDB()里面调cacheX()，所以getX()是可以保证X被cache的？
    - getXFromDB()，就把整个都检索出来，瓶颈在DB，所以处理其实问题不大？
- getX

### geodata-postid mapping redis cache

`(116.405285, 39.904989): "1784605228588990595"`

gpt估计每个地理位置元素可能需要大约 40 字节的存储空间，1G可容纳26,843,546 元素（可能遇到性能瓶颈吗？）

无失效时间，当数据库来用。不需要与任何数据存储同步

### postid-postcontent mapping redis cache

手动性失效的缓存，将postid映射到post数据

失效应该使一系列数据失效，同时将这一系列数据写入数据库，应该有一个后端监控redis的内存使用率然后进行写入备份

post:postId hash
- timestamp
- userId
- content
- likes
- bookmarks
- cacheScore

sets:
- likedBy
- tags
- pictureLinks
- commentIds



## rpc
connectrpc

## Postgres Database

post table
- PostId
- Timestamp
- AuthorId
- Content
- Likes
- Bookmarks
- PictureLinks
- Location
- CommentIds
- DeleteFlag
- Reports
- Rating

user table
- UserId
- OAuthId
- Name
- PostIds
- LikedPostIds
- CheckedPostIds
- BookmarkedPostIds
- DeleteFlag

comment table
- CommentId
- UserId
- Content
- Timestamp
- Likes
- DeleteFlag

report table
- ReportId
- UserId
- Content
- Timestamp


关于地理位置的存储，要考虑几个问题：
- 检索必须快速，意味着计算资源不能消耗太多，同时尽可能的利用redis缓存
- 如果说两个不同的东西在同一地理位置上，那么实际上没有必要将其分开，是可以当作属于一个聚类的，只不过
    - 或许还可以在聚类里面再次进行聚类，不通过地理位置而是通过别的
- 或许还可以再建立一个redis缓存， 其内容为
    - 聚类id - 聚类中心
    - 聚类id - 所属postid


关于身份验证
- 