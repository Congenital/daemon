# daemon

## run
> 以守护进程方式启动一个程序。
>> 默认为当前路径寻找程序。

> Example : 
>> run ./server -ip=127.0.0.1 -port=2000 > server.log
>
>> run server -ip=127.0.0.1 -port=2000 > server.log


## daemon
> 进程监察者。监察程序是否运行，不在运行列表则运行程序。
>> 可以使用 run 提升为守护进程，内部使用了run 命令进行守护进程的执行。
>
>> 需要将程序放入PATH路径中。才可以使监察者找到 run 执行路径。
>
>> -root 参数为程序根目录，可以为空
>
>> -filter 为程序过滤器，用以查找到关注的应用，并重启。需要与server参数中服务器参数保持一致。（逗号分割，参数需要加双引号， 例如: -filter="./account,home/admin/pushserver 10"）
>
>> -server 想要启动的真实路径。(逗号分割，参数需要加双引号,程序可以接收参数 例如: -server="./account,home/admin/pushserver 10")
>
>> -time 查询时间。一般就是服务器检查的时间间隔，也是重启间隔。
>
>> 默认会在所执行程序输出与错误输出重定向到程序所在目录下 （程序名.log）

> Example
>> daemon -root="/home/user/server" -server="./account,pushserver/push-server"-filter="account,pushserver/push-server" -time=10 > daemon.log &
>
>> run ~/bin/daemon -root="/home/user/server" -server="./account,pushserver/push-server"-filter="account,pushserver/push-server" -time=10 > daemon.log


