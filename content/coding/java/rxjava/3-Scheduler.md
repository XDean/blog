---
title: "理解RxJava与响应式(3) - 线程调度"
date: 2018-05-05
source: "https://github.com/XDean/Share/blob/master/doc/rxjava/3-Scheduler.md"
tags: 
  - "Knowledge-Sharing"
  - "RxJava"
categories:
  - "coding"
  - "java"
---

# Scheduler

RxJava 为我们提供了两个调度操作符`subscribeOn`和`observeOn`.

## How they work

![Scheduler.png](Scheduler.png)

## Delay error

![SchedulerDelayError.png](SchedulerDelayError.png)

## Schedulers

- `computation`
- `io`
- `single`
- `newThread`
- `trampoline`


[Sample Code](https://github.com/XDean/Share/blob/master/src/main/java/xdean/share/rx/ReactiveChapter3.java)

| Previous | Next |
| --- | --- |
| [Operator](2-Operator.html) |  [BackPressure](4-BackPressure.html)|