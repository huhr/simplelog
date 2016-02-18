### SimpleLog
简单的golang日志库，支持天级、小时级日期切割，支持多logger实例，配置简单。

##### 配置
```
{	
	"root": [
		{
			"Level": "info, debug, warn",
			"Output": "stdout",
			"Rotation":"",
			"Format": "detail"
		},
		{
			"Level": "error, fatal",
			"Output": "stderr",
			"Rotation":"",
			"Format": "detail"
		}
	],
	"logger_simpler": [
		{
			"Level": "info, debug, warn",
			"Output": "simple.debug",
			"Rotation":"daily",
			"Format": "detail"
		},
		{
			"Level": "error, fatal",
			"Output": "simple.err",
			"Rotation":"daily",
			"Format": "detail"
		}
	]
}
```
SimpleLog可配置多个logger对象，其中root为全局logger对象，配置文件中必须要配置root对象，每一个logger对象中可以包含多个节点，每一个节点代表一个输出对象，输出对象中包含日志的输出位置、输出格式、输出文件是否需要定时切割等信息。    
Level：指定该输出节点包含哪些日志级别，可选为关键字为info，debug，warn，error，fatal。注意区分大小写。    
Output：指定该输出节点输出的目标文件，可选为stdout，stderr或输出到磁盘文本。    
Rotation：指定该输出节点当输出目前为文本文件时，是否需要对文件进行定时切割，目前支持hourly，daily，其他值代表不切割。    
Format：指节点输出的日志格式，目前支持两种，detail和brife，格式分别为   
``` 
[time] [log_level] [logger name] [filename]:[line] [msg]  

[msg]
```
如果日志级别没有出现在任何一个节点的Level中，则该级别的日志不会有输出。同一个logger对象的日志级别重复出现时，总是以最先出现的配置为准。当多个输出节点指向同一个文件时，文件的切割规则以最先出现的配置为准。

##### 使用
```
import (
	log "github.com/huhr/simplelog"
)
func main() {
	err := log.LoadConfiguration("config.json")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	log.Debug("Hello root logger")

	// use multi logger
	logger, err := log.GetLogger("another")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	logger.Debug("Hello another logger")
}
```
