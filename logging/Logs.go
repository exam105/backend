package logging

import (
	"log"
	"os"
)

var (
	warningLogger *log.Logger
	infoLogger    *log.Logger
	errorLogger   *log.Logger
)

func init() {
	file, err := os.OpenFile("/var/logs/exam105/exam105.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Println( "Problem with exam105.log ", err.Error())
	}

	infoLogger = log.New(file, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	warningLogger = log.New(file, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
	errorLogger = log.New(file, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
}

func LogInformation(msg string) {
	infoLogger.Println(msg)
}

func LogWarning(warningmsg string) {
	warningLogger.Println(warningmsg)
}

func LogError(errorMsg string) {
	errorLogger.Println(errorMsg)
}
