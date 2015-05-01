# sdagent Summary

This is a project using skydns, etcd for service discoury by DNS records.

This is called sdagent whitch run on every service machines and update DNS records

on SKYDNS every period. it offer health check ,all that learn from consul by hashcorp.inc

eliminating the complex parts, keep it simple.

# Usage

1.Install go 

2.Get this code, run "go get" to get the libs using by sdagent

3.Run "go build"" in path of "myworkspace/agent"

4.Prepare etcd,skydns  env. Go and see etcd doc and skydns doc

5.Prepare config file whitch define services and healthchecks

6.Run agent -h to see Usage

# Config File Define

it is JSON file like these:

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

"http":"http://baidu.com"}]},

{"name":"mysql",

"port":8080,

"priority":10,

"node":"n10",

"weight":100,

"text":"for mysql",

"ttl":10,

"checks":[

{"name":"chk_1",

"id":"sc8877",

"ttl":10}

]}

]}

# Config File Details

Every machines can run multiple services

Every service has more than one health check ,only when all pass,it will update DNS records,

if health check is missing , will pass the check but warm by agent.

For More:

"services" means the services running in that machine. it has muti service

Each service has the these:

"key":string , SKYDNS need key-value to update DNS record ,if Key not set ,will gen auto by node+name

"name":string , means service name ,will be use to get Key when Key is not set.

"node":string , means machine name ,will be use to get Key when Key is not set.

"host":string , DNS record ,can be CNAME or A/AAAA .if not set ,will use machine's IP  

"port":int , not use now , keep it

"priority":int , not use now  

"weight":int ,not use now , can use for Weight RR 

"text":string , not use now ,for infomation add to 

"ttl":int , live time of DNS record , agent will update it by ttl-1, this need be re-consider

"checks": array of healt check define, Currently support three way of check, ttl ,script , http:

"name": string ,HealthCheck name

"id": string ,unique id of HealthCheck

"ttl": int ,when ttl is set ,check nothing ,return pass .

"script": string , run script and get return value ,0-pass ,1-warm ,2-fail

"http": string , get connection of http url, get status code ,<200-pass ,429-warm , others-fail

