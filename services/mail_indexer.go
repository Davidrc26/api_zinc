package services

import (
	"bufio"
	"encoding/json"
	"log"
	_ "net/http/pprof"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/Davidrc26/api_zinc.git/models"
)

func StartIndexing() {
	maildir := "./data/maildir"
	files, err := os.ReadDir(maildir)

	if err != nil {
		log.Println(err)
		return
	}

	// Limitar concurrencia basado en CPUs disponibles
	// Usar 2x el número de CPUs para balancear CPU y I/O
	maxWorkers := runtime.NumCPU() * 2
	if maxWorkers < 1 {
		maxWorkers = 1
	}

	log.Printf("Iniciando indexación con %d workers concurrentes (CPUs disponibles: %d)", maxWorkers, runtime.NumCPU())

	// Semáforo para limitar goroutines concurrentes
	semaphore := make(chan struct{}, maxWorkers)
	var wg sync.WaitGroup

	for _, f := range files {
		wg.Add(1)
		// Adquirir slot del semáforo
		semaphore <- struct{}{}

		go func(f os.DirEntry) {
			defer wg.Done()
			defer func() { <-semaphore }() // Liberar slot del semáforo

			result := ReadFolder(maildir + "/" + f.Name())
			bulk := models.Data{Index: "maildir", Records: result}
			jsonData, err := json.MarshalIndent(bulk, "", "  ")
			if err != nil {
				log.Println(err)
				return
			}
			BulkIndex(jsonData)
		}(f)
	}
	wg.Wait()
	log.Println("Indexación completada")
}

func ReadFolder(folder_name string) []models.Email {
	var files, err = os.ReadDir(folder_name)
	var object = make([]models.Email, 0)
	if err != nil {
		log.Println("Error leyendo el directorio" + folder_name + "\nDetalles: " + err.Error())
		return object
	}
	for _, f := range files {
		if f.IsDir() {
			object = append(object, ReadFolder(folder_name+"/"+f.Name())...)
		} else {
			object = append(object, ProcessFile(folder_name+"/"+f.Name()))
		}
	}
	return object
}

func ProcessFile(file string) models.Email {
	f, err := os.Open(file)
	if err != nil {
		log.Println("Error procesando el archivo " + file + "\nDetalles: " + err.Error())
		return models.Email{}
	}

	defer f.Close()

	fileInfo, err := f.Stat()
	if err != nil {
		log.Println("Error obteniendo información del archivo " + file + "\nDetalles: " + err.Error())
		return models.Email{}
	}

	// Ajustar capacidad del buffer basado en el tamaño del archivo
	// Mínimo: 64KB, Máximo: 2MB, Default: tamaño del archivo
	const minCapacity = 64 * 1024
	const maxCapacity = 2 * 1024 * 1024

	fileSize := fileInfo.Size()
	bufferCapacity := int(fileSize)

	if bufferCapacity < minCapacity {
		bufferCapacity = minCapacity
	} else if bufferCapacity > maxCapacity {
		bufferCapacity = maxCapacity
	}

	scanner := bufio.NewScanner(f)
	data := models.Email{}
	buf := make([]byte, 0, bufferCapacity)
	scanner.Buffer(buf, bufferCapacity)
	inHeaders := true
	var bodyLines []string

	for scanner.Scan() {
		line := scanner.Text()

		if inHeaders && strings.TrimSpace(line) == "" {
			inHeaders = false
			continue
		}

		if inHeaders {
			// Procesar headers
			i := strings.Index(line, ":")
			if i >= 0 {
				key := strings.TrimSpace(line[:i])
				value := strings.TrimSpace(line[i+1:])
				switch key {
				case "Message-ID":
					data.Message_ID = value
				case "Date":
					data.Date = parseDate(value)
				case "From":
					data.From = value
				case "To":
					data.To = value
				case "Subject":
					data.Subject = value
				case "Mime-Version":
					data.Mime_Version = value
				case "Content-Type":
					data.Content_Type = value
				case "Content-Transfer-Encoding":
					data.Content_Transfer_Encoding = value
				case "X-From":
					data.X_From = value
				case "X-To":
					data.X_To = value
				case "X-cc":
					data.X_cc = value
				case "X-bcc":
					data.X_bcc = value
				case "X-Folder":
					data.X_Folder = value
				case "X-Origin":
					data.X_Origin = value
				case "X-FileName":
					data.X_FileName = value
				case "Cc":
					data.Cc = value
				}
			}
		} else {
			bodyLines = append(bodyLines, line)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Println("Error escaneando el archivo " + file + "\nDetalles: " + err.Error())
	}

	data.Body = strings.Join(bodyLines, "\n")

	return data
}

func parseDate(dateStr string) time.Time {
	layouts := []string{
		time.RFC1123Z,
		time.RFC1123,
		time.RFC822Z,
		time.RFC822,
		time.ANSIC,
		"Mon, 2 Jan 2006 15:04:05 -0700 (MST)",
		"2 Jan 2006 15:04:05 -0700",
	}
	for _, layout := range layouts {
		if t, err := time.Parse(layout, dateStr); err == nil {
			return t
		}
	}
	log.Println("No se pudo parsear la fecha: " + dateStr)
	return time.Time{}
}
