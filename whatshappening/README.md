# 项目简介
提供某一时间内发生的热点新闻关键词聚合绘制词云图。


创建了kratos项目，基础功能并没有实现
// TODO:
停用词表: https://blog.csdn.net/dilifish/article/details/117885706


测试；
```
curl -G   -d "plantforms=weibo"   -d "plantforms=baidu"   -d "is_exclude=false"   -d "limit=10"   "http://localhost:8000/stats/wordcount"
```

# 部署
## Docker
```bash
# build
docker build -t <your-docker-image-name> .

# run
docker run --rm -p 8000:8000 -p 9000:9000 -v /home/xxx/whatshappening/configs/:/data/config hotword
# 如果copy了一份配置，则可以这样启动：
docker run -d --rm -p 8000:8000 -p 9000:9000  hotword
```


