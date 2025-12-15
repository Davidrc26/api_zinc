package services

import (
	"log"
	"os"
	"runtime"
	"runtime/pprof"

	"github.com/Davidrc26/api_zinc.git/models"
)

type ProfilingFiles struct {
	CPU *os.File
	Log *os.File
}

// StartCPUProfiling inicia el profiling de CPU y retorna el archivo o error
func StartCPUProfiling() (*os.File, error) {
	cpu, err := os.Create("profiling/cpu.prof")
	if err != nil {
		return nil, err
	}
	err = pprof.StartCPUProfile(cpu)
	if err != nil {
		cpu.Close()
		return nil, err
	}
	return cpu, nil
}

// StopCPUProfiling detiene el profiling de CPU
func StopCPUProfiling(cpu *os.File) {
	if cpu != nil {
		pprof.StopCPUProfile()
		cpu.Close()
	}
}

// CreateLogFile crea y abre el archivo de log
func CreateLogFile() (*os.File, error) {
	logFile, err := os.OpenFile("logsindexer/log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}
	log.SetOutput(logFile)
	return logFile, nil
}

// CreateMemoryProfile crea el perfil de memoria
func CreateMemoryProfile() models.Response {
	runtime.GC()
	mem, err := os.Create("profiling/memory.prof")
	if err != nil {
		return models.Response{
			Status:  500,
			Message: "Error al crear perfil de memoria: " + err.Error(),
			Result:  nil,
		}
	}
	defer mem.Close()

	if err := pprof.WriteHeapProfile(mem); err != nil {
		return models.Response{
			Status:  500,
			Message: "Error al escribir perfil de memoria: " + err.Error(),
			Result:  nil,
		}
	}

	return models.Response{
		Status:  200,
		Message: "Perfil de memoria creado exitosamente",
		Result:  nil,
	}
}
