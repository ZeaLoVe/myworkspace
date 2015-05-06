rm -f /etc/init.d/sdagent
rm -f /usr/bin/agent
rm -f /etc/sdconfig.json

cd $GOPATH/src/sdagent/agent
go clean
go build

cp sdagent /etc/init.d/
cp agent /usr/bin/
cp sdconfig.json /etc/

f [ -e /etc/monit.conf ];then
        echo "monit exist in this machine"
        cp monit-file /etc/monit.d/sdagent
fi

if [ ! -e /etc/init.d/sdagent ];then
        echo "sdagent install fail"
        return 1
fi

if [ ! -e /etc/sdconfig.json ];then
        echo "sdagent install fail"
        return 1
fi

if [ ! -e /usr/bin/agent ];then
        echo "sdagent install fail"
        return 1
fi
