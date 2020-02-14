---
title: "Spring Boot 番外 (04) - 缓存 Cache"
date: 2020-02-10
tags: 
  - "Knowledge-Sharing"
categories:
  - "Coding"
  - "Java"
  - "Spring Boot"
src:
    - text: 'Custom Encrypt'
      url: https://github.com/XDean/Share/tree/master/src/spring-boot/encrypt/src/main/java/xdean/share/spring/encrypt/custom
    - text: 'Jasypt Encrypt'
      url: https://github.com/XDean/Share/tree/master/src/spring-boot/encrypt/src/main/java/xdean/share/spring/encrypt/jasypt
---

# 什么是缓存

对于程序员，缓存一定不是一个陌生的概念。
再我们编写的程序中，数据的处理和传递无时无刻不在发生着。
而其中总会存在着一些的重复的操作和数据。

- 重复的处理会浪费CPU
- 重复的读写会浪费磁盘
- 重复的传输会浪费网络

为了减少这种浪费，也为了加快响应时间，我们就需要缓存。
把处理好的数据保留下来(最常见的是放在内存)，下次再要做相同的事情，就可以直接返回结果。

# Spring Boot 缓存

Spring Boot核心库就为我们提供了缓存功能，相关代码在`org.springframework.cache`包下。

要开启缓存功能，只需要添加上`@EnableCaching`注解即可。
接着我们就要在需要的方法上加上注解来实现缓存。
Spring Boot提供了以下几种注解。

- `@Cacheable`, 该方法会使用缓存
- `@CachePut`, 该方法不使用缓存，但是结果会置入缓存
- `@CacheEvict`, 该方法不使用缓存，而是将匹配的缓存删除
- `@Caching`, 复合了以上三个注解，通常你不会用到
- `@CacheConfig`, 对类下的所有缓存进行默认配置

其中，前三个注解是核心，每个注解都有若干配置，所以配置的含义都会在后面一一讲解。

我们先从最简单的例子看起

## `@Cacheable`

我们现在有一个方法需要开启缓存

```java
public String hello(String who) {
    // 假设该方法消耗了大量资源，需要缓存
    return "Hello " + who;
}
```

我们只需要在这个方法上加上`@Cacheable`即可。
但是注意，同时必须要指定`cacheNames`。

### `cacheNames`

`cacheNames`的存在类似于作用域。
想象你有一个数据库里有两张表，都需要通过主键ID来缓存。
这时你就需要两个互不相关的缓存，即`cacheNames = "table_1"`和`cacheNames = "table_2"`。
这样即使同样是查询`ID = 1`，也不会相互干扰，只会命中自己这个域中的缓存。
这里我们仅作示例，随便取一个缓存名。

_*有心的读者应该注意到了复数形式，一个缓存可以定义多个缓存域，在示例代码中有相关的测试_

```java
@Cacheable(cacheNames = "test")
public String hello(String who) {
    System.out.println("Calculating Hello: " + who); // 打印日志以方便知道调用了几次
    return "Hello " + who;
}
```

现在让我们来测试代码

```java
System.out.println(app.hello("World"));
System.out.println(app.hello("World"));
System.out.println(app.hello("XDean"));
System.out.println(app.hello("XDean"));
```

输出结果为

    Calculating Hello: World
    Hello World
    Hello World
    Calculating Hello: XDean
    Hello XDean
    Hello XDean
    
可以看到，虽然连续调用了两次，但是只有第一次真正地调用了方法。
而如果用不同的参数去调用，第一次也会真正地调用方法。

那么`Cacheable`是如何区分不同参数调用的缓存命中呢？答案就是`key`参数

### `key`

对于每次缓存调用，都会计算出一个关键值，即`key`，通过这个值来寻找缓存。
默认情况下，`key`通过所有参数计算(详见`SimpleKeyGenerator`)。
你也可以通过`key`属性来配置自定义的值。
`key`的值是一个SpEL表达式，其中可用的参数有