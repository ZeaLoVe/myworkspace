#Current Version TODO List

-

#Next Version TODO List

1.通过增加http服务线程，对外公开 Agent上运行的Service的情况

可以提供查询和删除Job等，也可以新增。。但主要是暴露运行情况，其他待考虑

2.通过在etcd上设置存储路径 /v2/sdagent/ 在存储路径下，设置一定的服务反馈格式。

应用APP可以通过向etcd写入反馈数据，Agent根据反馈数据来决定当前服务的权重和优先级

实现动态轮询和降级等

格式：类似DNS记录格式，Key就是DNS记录的Key，value包含 可供参考的运行时QPS（比如平均QPS，当前主机QPS）只是比如
