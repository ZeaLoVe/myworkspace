# sdagent 

This is a project using skydns,etcd for service discoury.using DNS records

The Major part is a agent called sdagent whitch run in every service machines and update DNS records

on SKYDNS every period. it offer a health check ,all that learn from consul by hashcorp.inc

eliminating the complex parts .

# Usage

1.install go 

2.get this code,run "go get" to the libs using by sdagent

3.run go build in the path of myworkspace/agent

4.prepare etcd,skydns  env

5.prepare config file whitch define services and healthcheck

6.run agent using -f point to config file , -e point to etcd machines

# Config file define

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
"script":"ifconfig",
},
{"name":"chk_2",
"id":"sc321",
"http":"http://baidu.com"}]
},
{
"name":"mysql",
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

# config detail

every machines can run more than one service

every service has more than one health check ,only when all pass,it will update DNS records

"services" is the service define running in that machine. it has:

"key":string , SKYDNS need key-value to update DNS record ,if Key not set ,will gen auto: node+name

"name":string , means service name ,will be use to get Key if Key is not set.

"node":string , means machine name ,will be use to get Key if Key is not set.

"host":string , DNS record ,can be CNAME or A/AAAA .if not set ,will use machine's IP  

"port":int , not use now , keep it

"priority":int , not use now  

"weight":int ,not use now , can use for Weight RR 

"text":string , not use now ,for infomation add to 

"ttl":int , live time of DNS record , agent will update it by ttl/2

"checks", array of healt check define ,Currently support three way of check, ttl ,script , http:

"name" , string ,HealthCheck name

"id", string ,unique id of HealthCheck

"ttl" , int ,when ttl is set ,check nothing ,return pass .

"script",string , run script and get return value ,0-pass ,1-warm ,2-fail

"http" , string , get connection of http url, get status code ,<200-pass ,429-warm , others-fail

