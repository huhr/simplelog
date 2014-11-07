# Simple Log
简单的日志实现
支持配置日志等级
分级别单独配置
支持日志文件每天和每小时切割
支持配置format和log.file文件
# TODO
文档支持
更多支持和代码的完善

配置文件log.cfg:

		[basic]
		level=DEBUG

		[debug]
		file=/home/huhaoran/sos/go_ui/logs/debug.log
		hourly=true
		daily=true
		format=TEXT

		[info]
		file=/home/huhaoran/sos/go_ui/logs/info.log
		hourly=true
		daily=true
		format=BRIFE

		[warn]
		file=/home/huhaoran/sos/go_ui/logs/warn.log
		hourly=true
		daily=true
		format=TEXT

		[error]
		file=/home/huhaoran/sos/go_ui/logs/error.log
		hourly=true
		daily=true
		format=TEXT

加载配置文件:

		var log_cfg log.Config
		gcfg.ReadFileInto(&log_cfg, *exeDir+"/conf/log.cfg")
		log.LoadConfiguration(log_cfg)
		log.Debug("load all the cfg file")

