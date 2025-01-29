## Cache Server
in real case, cache server would be installed right onto the vm itself, without docker
this repository meant to initialize local development with redis

## Cache Client
cacheclient here is the codebase of cache client. It can be easily copied to any service that needs distributed caching

on .env.example and .env, add the following variables

### Cache Client using Redis
```
# REDIS
REDIS_HOST=
REDIS_PORT=6379
```

and add the following dependencies
```
go get github.com/redis/go-redis/v9
```

Ref: https://github.com/redis/go-redis

### Dependency Injection of Cache Client
```
do.Provide[serviceCache.CacheClient](Injector, serviceCache.NewRedisCacheClientInject)
if os.Getenv("MODE") == "DEBUG" {
    do.Provide[serviceCache.CacheClient](Injector, serviceCache.NewMockCacheClientInject)
}
do.Provide[serviceCache.CacheService](Injector, serviceCache.NewCacheServiceInject)
```