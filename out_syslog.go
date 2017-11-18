package main

import "github.com/fluent/fluent-bit-go/output"
import (
	"os"
	"strings"
	"fmt"
	"unsafe"
	"C"
	"log/syslog"
)

var (
	severitiesMap = map[string]syslog.Priority{
		"EMERG":   syslog.LOG_EMERG,
		"ALERT":   syslog.LOG_ALERT,
		"CRIT":    syslog.LOG_CRIT,
		"ERR":     syslog.LOG_ERR,
		"WARNING": syslog.LOG_WARNING,
		"NOTICE":  syslog.LOG_NOTICE,
		"INFO":    syslog.LOG_INFO,
		"DEBUG":   syslog.LOG_DEBUG,
	}

	facilitiesMap = map[string]syslog.Priority{
		"KERN":     syslog.LOG_KERN,
		"USER":     syslog.LOG_USER,
		"MAIL":     syslog.LOG_MAIL,
		"DAEMON":   syslog.LOG_DAEMON,
		"AUTH":     syslog.LOG_AUTH,
		"SYSLOG":   syslog.LOG_SYSLOG,
		"LPR":      syslog.LOG_LPR,
		"NEWS":     syslog.LOG_NEWS,
		"UUCP":     syslog.LOG_UUCP,
		"CRON":     syslog.LOG_CRON,
		"AUTHPRIV": syslog.LOG_AUTHPRIV,
		"FTP":      syslog.LOG_FTP,
		"LOCAL0":   syslog.LOG_LOCAL0,
		"LOCAL1":   syslog.LOG_LOCAL1,
		"LOCAL2":   syslog.LOG_LOCAL2,
		"LOCAL3":   syslog.LOG_LOCAL3,
		"LOCAL4":   syslog.LOG_LOCAL4,
		"LOCAL5":   syslog.LOG_LOCAL5,
		"LOCAL6":   syslog.LOG_LOCAL6,
		"LOCAL7":   syslog.LOG_LOCAL7,
	}
)

//export FLBPluginRegister
func FLBPluginRegister(ctx unsafe.Pointer) int {
	return output.FLBPluginRegister(ctx, "syslog", "Syslog")
}

var Config struct {
	network string
	address string
	priority syslog.Priority
	tag string
}

//export FLBPluginInit
// (fluentbit will call this)
// ctx (context) pointer to fluentbit context (state/ c code)
func FLBPluginInit(ctx unsafe.Pointer) int {
	// Example to retrieve an optional configuration parameter
	Config.network = output.FLBPluginConfigKey(ctx, "network")

	Config.address = output.FLBPluginConfigKey(ctx, "address")

	severity := output.FLBPluginConfigKey(ctx, "severity")
	if severity == "" {
		severity = "INFO"
	}

	facility := output.FLBPluginConfigKey(ctx, "facility")
	if facility == "" {
		facility = "LOCAL0"
	}

	priority, err := getSyslogPriority(severity, facility)
	if err != nil {
		fmt.Fprintln(os.Stderr, "[out_syslog] failed to initialize: ", err);
		return output.FLB_ERROR
	}
	Config.priority = priority

	Config.tag = output.FLBPluginConfigKey(ctx, "tag")

	return output.FLB_OK
}

//export FLBPluginFlush
func FLBPluginFlush(data unsafe.Pointer, length C.int, tag *C.char) int {
	var ret int
	var record map[interface{}]interface{}
	var syslogTag string

	// Create Fluent Bit decoder
	dec := output.NewDecoder(data, int(length))
	if Config.tag == "" {
		syslogTag = C.GoString(tag)
	}else{
		syslogTag = Config.tag
	}

	sysLog, err := syslog.Dial(Config.network, Config.address, Config.priority, syslogTag)
	if err != nil {
		fmt.Fprintln(os.Stderr, "[out_syslog] failed to send logs: ", err);
		return output.FLB_ERROR
	}

	// Iterate Records
	for {
		// Extract Record
		ret, _, record = output.GetRecord(dec)
		if ret != 0 {
			break
		}

		str := ""
		for k, v := range record {
			str += fmt.Sprintf("%s=%v ", k, v)
		}

		// Print record keys and values
		fmt.Fprint(sysLog, str)

	}

	sysLog.Close()

	// Return options:
	//
	// output.FLB_OK    = data have been processed.
	// output.FLB_ERROR = unrecoverable error, do not try this again.
	// output.FLB_RETRY = retry to flush later.
	return output.FLB_OK
}

//export FLBPluginExit
func FLBPluginExit() int {
	return output.FLB_OK
}

func getSyslogPriority(severity string, facility string) (syslog.Priority, error) {
	severityPriority, ok := severitiesMap[strings.ToUpper(severity)]
	if !ok {
		return 0, fmt.Errorf("Unknown syslog severity '%s'", severity)
	}

	facilityPriority, ok := facilitiesMap[strings.ToUpper(facility)]
	if !ok {
		return 0, fmt.Errorf("Unknown syslog facility '%s'", facility)
	}

	return severityPriority | facilityPriority, nil
}

func main() {
}
