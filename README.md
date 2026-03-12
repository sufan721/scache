### 项目结构：
```go
├─cache
│  ├─cache.go        (并发安全缓存)
│  ├─group.go        ⭐核心入口
│  └─lru
│      └─lru.go
├─consistenthash
│  └─consistenthash.go
├─httpserver
│  └─http.go
├─peer
│  └─peer.go
└─main.go
```
### 后面的优化(未完成)：

1. 缓存击穿保护
2. 虚拟节点(核心)
3. 节点通信优化
4. 淘汰策略升级(策略插件化)
5. 热点缓存
6. 监控指标
7. 缓存三件套
```
	缓存击穿
	缓存雪崩
	缓存穿透
```
