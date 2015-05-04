#项目简介

本项目基于 etcd、skydns项目，以etcd作为数据后端，skydns作为提供域名服务的部分进行服务发现

SDAgent是一个长久运行在服务节点对注册服务进行健康检查，同时根据检测结果将服务注册到数据后端etcd中，

再由SKYDNS读取etcd的数据对外暴露DNS域名进行服务定位。SDAgent默认通过给定的域名发现etcd后端的IP

所以需要其指向的DNS服务器上事先配置好etcd服务器的IP 和域名。默认zealove.xyz 是我自己的注册的

公共域名，做测试代码使用。

#项目架构说明

数据端：etcd

DNS服务端：SKYDNS

服务注册和检测端：SDAgent

SDAgent----> etcd <----> SKYDNS

1..*     group(3,6,9..)   1..*

#初次使用

安装GO语言配置GO语言编译环境

在工作目录下获取本代码 git clone ...

进入agent目录下，执行 go build 会生成可执行文件 

只要agent 路径下默认的sdaconfig.json文件配置完成。

直接 ./agent 执行即可

#运行参数介绍

./agent -h 可获取帮助信息，可供配置的参数主要是4个

需注册服务的配置文件来源         use -f=filapath  默认当前路径下的sdconfig.json

用于获取etcd集群IP地址的域名     use -d=DNS       默认 zealove.xyz 我的测试域名

写入Etcd的端口                use -p=port      默认2379

检测配置文件改动重新加载执行的周期 use -t=num       默认30 单位分钟


#配置文件说明

{"services":[

{"name":"mongo",

"port":2334,

"priority":10,

"node":"n3",

"weight":90,

"text":"for mongo",

"ttl":5,

"checks":[

{"name":"chk_1",

"id":"sc123",

"script":"ifconfig",},

{"name":"chk_2",

"id":"sc321",

"http":"http://baidu.com"}]}

]}

每个配置文件包含一个services字段，每个services里是具体的service配置，每个service又可以配置多个

健康检查，将会根据所有健康检查执行的结果进行综合判断，所有都Pass才算是Pass，只要有一个warn就会返回warn

只要一个error就返回error。每个service都会单独的协程处理健康检查和更新记录，都支持固定的超时检查（不可配）

具体配置文件说明如下

check ：对应多个健康检查，每个检查必须设置的字段仅仅 ttl\http\script 三者选一个

name和id可以根据需要设置，建议必须设置id作为标识，虽然目前并没有使用到

ttl类型是整形，这里数值并不重要，只要不为0即可，会根据service更新的制定频率，直接Pass检查，更新服务

http类型是字符串，根据给定的URL，发起http请求，根据返回码决定结果，只有200为pass，服务器繁忙为warn，其他都是error

script类型为字符串，Windows平台和Linux平台使用平台可识别的执行语句或者脚本，执行这些脚本后系统返回值将作为结果依据

0为pass，1为warn，2为error

service 是基本的单元，每个service有几个字段可以设置，主要用于更新服务记录的

包括 name/port/priority/node/weight/text 还有两个默认可以生成的但也可以手动给定的

key/host 下面将说明如何进行设置。

etcd本质上是一个提供KV服务的系统，需要将所有数据转换为键值对。而SKYDNS读取etcd的格式是

比如域名 www.google.com 在etcd中的key为 /skydns/com/google/www

value则为 {“host”：“192.198.11.1”，“weight”：100，“priority”：20，“”：..} 

如果给定key 和 host，则会配置好域名，比如key:"www.google.com" 既可以将改服务的域名直接设置成

www.google.com  否则会根据 node+name生成域名，比如 node = n3 name = mongo.sdp 则域名为n3.mongo.sdp

host默认如果没设置，会获取所在机器的本地IP。

最小化运行：name是必须设置的，否则服务的定义会比较模糊，不符合约定。check可以缺省，但会警告，其他值会生成默认值如下：

port = 8080

weight = 100

priority = 20

ttl = 10

text = "default text"



