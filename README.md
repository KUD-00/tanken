[system architecture](/resources/architecture.png)

[micro-service architecture](/resources/architecture-micro-service.png)

mapbox的价格大概是google map的1/2，可能要考虑从mapbox迁移到google map

db和redis里面出来的数据还是不一样，所以type还是需要写两套，能否复用？好像意义不大。

现在还是不够微服务，比如说user啊post啊都应该分开？拥有独立的数据库且并不互相强关联，这种情况可能需要不选用sql数据库而是nosql, 感觉sql和微服务的理念不合