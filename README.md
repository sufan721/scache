
---

# gocache

一个使用 **Go** 实现的轻量级分布式缓存系统，用于学习缓存架构设计与高并发优化。

项目模拟真实生产环境中的缓存系统，实现了缓存淘汰策略、分布式节点、缓存保护机制以及数据库后端等核心功能。

---

# 项目亮点

- 支持 **LRU / LFU 淘汰策略（策略插件化）**
    
- 使用 **一致性哈希** 实现分布式节点路由
    
- 使用 **BloomFilter** 防止缓存穿透
    
- 使用 **Singleflight** 防止缓存击穿
    
- 使用 **TTL + 随机过期** 防止缓存雪崩
    
- 支持 **HotCache 热点数据缓存**
    
- 支持 **MySQL 后端数据源**
    
- 支持 **缓存统计指标**
    
- 支持 **高并发访问**
    

---

# 系统架构

```
        +--------+
        | Client |
        +--------+
             |
             v
      +-------------+
      | HTTP Server |
      +-------------+
             |
             v
      +-------------+
      |   Cache     |
      |  LRU / LFU  |
      +-------------+
             |
      +--------------+
      | Bloom Filter |
      +--------------+
             |
      +--------------+
      | Singleflight |
      +--------------+
             |
             v
          +-------+
          | MySQL |
          +-------+
```

数据库：

- MySQL
---

# 请求流程

```
Client Request
      |
      v
HTTP Server
      |
      v
BloomFilter 判断 key 是否存在
      |
      ├── 不存在 → 直接返回
      |
      v
本地缓存查询 (LRU / LFU)
      |
      ├── 命中 → 返回数据
      |
      v
singleflight 控制并发请求
      |
      v
查询 MySQL
      |
      v
写入缓存
      |
      v
返回数据
```

---

# 解决缓存三大问题

缓存系统通常会遇到三个经典问题：

|问题|描述|解决方案|
|---|---|---|
|缓存穿透|查询不存在的数据|BloomFilter|
|缓存击穿|热点 key 同时失效|Singleflight|
|缓存雪崩|大量缓存同时过期|随机 TTL|

核心算法：

- Bloom Filter
    
- Consistent Hashing
    

---

# 项目结构

```
gocache
│
├── cache
│   ├── lru          LRU缓存策略
│   ├── lfu          LFU缓存策略
│   ├── bloom        布隆过滤器
│   └── group.go     缓存核心逻辑
│
├── consistenthash   一致性哈希实现
│
├── httpserver       HTTP服务
│
├── peer             分布式节点通信
│
└── main.go          程序入口
```

---

# 启动项目

安装依赖：

```
go mod tidy
```

启动缓存集群：

```
go run main.go -port=8001
go run main.go -port=8002
go run main.go -port=8003
```

启动成功后：

```
cache node running at http://localhost:8001
```

---

# API 示例

查询缓存：

```
http://localhost:8001/_gocache/scores/Tom
```

返回：

```
630
```

---

# 性能测试

使用压测工具：

- hey
    

测试命令：

```
hey -n 10000 -c 100 http://localhost:8001/_gocache/scores/Tom
```

测试结果：

```
Requests/sec: ~60000
Average latency: ~1.5ms
P99 latency: ~17ms
```

说明系统在高并发情况下仍然保持低延迟。

---

# 技术栈

- Go
    
- MySQL
    
- HTTP
    
- 一致性哈希
    
- BloomFilter
    
- Singleflight
    

---

# 后续优化

未来可以继续优化：

- 缓存分片（Shard Cache）
    
- gRPC 节点通信
    
- Prometheus 监控
    
- 自动节点发现
    
- 持久化缓存
    

---

# 项目总结

本项目实现了一个完整的分布式缓存系统，涵盖：

- 缓存淘汰策略
    
- 分布式节点
    
- 高并发控制
    
- 缓存保护机制
    
- 数据库后端
    

通过该项目可以深入理解缓存系统设计以及高并发系统的核心原理。