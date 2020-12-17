package main


func initLog(logPath, logName, chanSize, splitType, splitSize string)  {
	config := make(map[string]string)
	config["log_path"] = logPath
	config["log_name"] = logName
	config["log_chan_size"] = chanSize
	config["log_split_type"] = splitType
	config["log_split_size"] = splitSize
	err := InitLog("file", config)
	if err != nil {
		return
	}
}

func Run()  {
	Warn("this is warn")
}

func main()  {
	initLog(".", "server", "50000", "size", "104857600")
	Run()
}
