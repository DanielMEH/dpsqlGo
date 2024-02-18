package logger

// define las importaciones
import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"time"

	"github.com/DanielMEH/database/settings"
)

// define las constantes que voy a utilizar en logs
const (
	logNone = iota
	logInfo
	logWarning
	logError
	logVerbose
	logDebug
)

var (
	logIdFunc   string
	levelStr    string
	colorAir    string
	colorblue   = "\033[34m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	_           = "\033[35m"
	colorCyan   = "\033[36m"
	colorWhite  = "\033[37m"
	colorRed    = "\033[31m"
	colorReset  = "\033[0m"
)

type Loggers interface {
	Log(level int, message string) error
	Stop() error
}

type fileLogger struct {
	logger   *log.Logger
	logFile  *os.File
	logLevel int
}

type HeadersRequest struct {
	Authorization string
	Actions       string
	LogSecId      string
	DateInfo      string
	Ip            string
	Username      string
	_             int
}

func newFileLogger() *fileLogger {
	return &fileLogger{
		logger:   nil,
		logFile:  nil,
		logLevel: logNone,
	}

}

func (myLogger *fileLogger) startLog(level int, file string) error {

	err := os.MkdirAll(filepath.Dir(file), os.ModePerm) // Asegurarse de que el directorio exista
	if err != nil {
		return err
	}

	f, err := os.OpenFile(file, os.O_WRONLY|os.O_CREATE|os.O_APPEND, os.ModePerm)
	if err != nil {
		return err
	}

	// Utilizar go-colorable para cambiar los colores en la consola
	multiWriter := io.MultiWriter(io.Discard, f)

	myLogger.logger = log.New(multiWriter, "", 0)
	myLogger.logLevel = level
	myLogger.logFile = f

	return nil
}

func (myLogger *fileLogger) stopLog() error {
	if myLogger.logFile != nil {
		return myLogger.logFile.Close()
	}
	return nil
}
func PrintColor(colorAir string, text string) string {
	return fmt.Sprint(string(colorAir), text, string(colorReset))
}
func (myLogger *fileLogger) log(level int, message string, callerSkip int) error {
	if myLogger.logger == nil {
		return errors.New("myFileLogger is not initialized correctly")
	}

	switch level {
	case logDebug:
		levelStr = "Debug"
		colorAir = colorCyan
	case logInfo:
		levelStr = "Info"
		colorAir = colorGreen
	case logWarning:
		levelStr = "Warning"
		colorAir = colorYellow
	case logError:
		levelStr = "Error"
		colorAir = colorRed
	default:
		return fmt.Errorf("invalid log level: %d", level)
	}

	if level >= myLogger.logLevel {
		// Obtener información sobre el lugar donde se produjo el error
		_, file, line, _ := runtime.Caller(2 + callerSkip)
		fileParts := strings.Split(file, "/")
		fileName := fileParts[len(fileParts)-1]

		logMessage := fmt.Sprintf("[%v] | logid:[%s] | level: [%v] | file: (%s:%v) | %s",
			PrintColor(colorWhite, dateLog("")),
			PrintColor(colorYellow, logId(logIdFunc)),
			PrintColor(colorAir, levelStr),
			PrintColor(colorblue, fileName),
			PrintColor(colorblue, fmt.Sprintf("%v", line)),
			PrintColor(colorAir, message))

		coloredLogMessage := logMessage
		fmt.Println(coloredLogMessage)
		plainMessage := regexp.MustCompile("\x1b\\[[0-9;]*m").ReplaceAllString(coloredLogMessage, "")
		myLogger.logger.Print(plainMessage)
	}
	return nil
}

func Logger(level int, ms interface{}, callerSkip int) {
	fileName, err := settings.LoadConfig()

	if err != nil {
		panic(err)
	}
	packageDir, err := filepath.Abs(".")
	if err != nil {
		panic(err)
	}

	unionPath := filepath.Join(packageDir, fileName.Log.FileName)

	// imprimir mi config
	logger := newFileLogger()
	if err := logger.startLog(level, unionPath); err != nil {
		panic(err.Error())
	}

	defer func() {
		logger.stopLog()
	}()

	// Generar un identificador de log único (puedes implementar tu propia lógica aquí)

	// Formatear el mensaje según el tipo
	var formattedMessage string
	switch v := ms.(type) {
	case string:
		formattedMessage = v
	case int, int32, int64, float32, float64, bool:
		formattedMessage = fmt.Sprintf("%v", v)
	default:
		// Si el tipo no es manejado, intenta convertir a cadena
		formattedMessage = fmt.Sprintf("%v", ms)
	}

	// Llamar a log con el parámetro adicional para indicar la cantidad de niveles a retroceder
	logger.log(level, formattedMessage, callerSkip)
}

func Logsdebug(ms interface{}) {
	Logger(logDebug, ms, 1)
}
func LogsInfo(ms interface{}) {
	Logger(logInfo, ms, 1)
}
func LogsWarning(ms interface{}) {
	Logger(logWarning, ms, 1)
}
func LogsError(ms interface{}) {

	Logger(logError, ms, 1)
}

func dateLog(newDate string) string {
	if newDate == "" {
		return time.Now().Format("2006-01-02 15:04:05")
	} else {
		return time.Now().Format(newDate)
	}
}

func logId(logIdentity string) string {
	if logIdentity == "" {
		return fmt.Sprintf("%d", time.Now().UnixNano())
	} else {
		return fmt.Sprintf("%v", logIdentity)
	}

}

func BuildTransformData(infoReq *HeadersRequest) *HeadersRequest {
	logId(infoReq.LogSecId)
	logIdFunc = infoReq.LogSecId

	return infoReq
}
