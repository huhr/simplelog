## SimpleLog
简单 高效

配置：
SimpleLog使用json格式的配置文件，可读性较高同时也比较容易
解析。

json配置文件中可以配置多个节点，一个节点表示一个日志输出目标，每个节点包含四个配置项：    
Level：该节点所处理的日志级别，可选关键字为info，debug，warn，error，fatal。注意区分大小写    
Out：节点对应日志级别输出目标，stdout，stderr，以及输出文件。    
Cut：节点输出文件时，文件切割选项，目前支持hourly，daily，以及不切割。    
Format：指节点输出的日志格式，目前支持两种，detail和brife，格式分别为   

``` 
[log_level] [time] [filename]:[line] [msg]  

[msg]
```
如果某日志级别没有出现在任何一个节点的Level中，name该级别的日志没有输出。不同的节点配置不能   
输出到同一个日志文件中。


```
[
	{
		"Level": "info, debug, warn",
		"Out": "logs/server.debug",
		"Cut":"daily",
		"Format": "detail"
	},
	{
		"Level": "error, fatal",
		"Out": "logs/server.err",
		"Cut":"daily",
		"Format": "detail"
	}
]
```

运行go test，查看输出
