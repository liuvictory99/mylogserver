package logutil

import (
	"encoding/json"
	"github.com/cihub/seelog"
)

type logcontent struct { 
	Access string  `json:"access"`	  
	Path  string  `json:"path"` 
}



var Logger seelog.LoggerInterface

func init() {
    DisableLog()
}

// DisableLog disables all library log output
func DisableLog() {
    Logger = seelog.Disabled
}

// UseLogger uses a specified seelog.LoggerInterface to output library log.
// Use this func if you are using Seelog logging system in your app.
func UseLogger(newLogger seelog.LoggerInterface) {
    Logger = newLogger
}


func loadAppConfig(path string) {
    appConfig := `
<seelog>
    <outputs formatid="main">
        <filter levels="debug">
            <console />
        </filter>
        <filter levels="info">

				<file path="logs/`+path+`"/>
		
        </filter>
    </outputs>

    <formats>
        <format id="main" format="%Date/%Time  %Msg%n"/>
    </formats>
</seelog>
`
    logger, err := seelog.LoggerFromConfigAsBytes([]byte(appConfig))
    if err != nil {
        seelog.Debug("err parsing config log file", err)
        return
    }
    UseLogger(logger)
}


func Write(logbyte []byte) {
	defer seelog.Flush()
	var lc logcontent
	err := json.Unmarshal(logbyte,&lc)
	
	if err != nil {
	   return
	}
	//实例化配置文件
	loadAppConfig(lc.Path)


	seelog.ReplaceLogger(Logger)	
	seelog.Info(lc.Access)

}