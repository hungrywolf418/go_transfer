package handler

import (
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/sirupsen/logrus"
)

func Upload(response http.ResponseWriter, request *http.Request) {
	var err error
	file, fileHeader, err := request.FormFile(`file`)
	if err != nil {
		logrus.Error("Was not able to access the uploaded file: ", err)
		http.Error(response, err.Error(), http.StatusBadRequest)
		return
	} else {

		// Close the file afterwards:
		defer file.Close()

		// Get the original filename:
		sourceFilename := fileHeader.Filename

		// Read the entire file into memory:
		data, err := ioutil.ReadAll(file)
		if err != nil {
			logrus.Error("Error while reading file from client: ", err)
			http.Error(response, err.Error(), http.StatusInternalServerError)
			return
		} else {
			currDir, err := os.Getwd()
			if err != nil {
				logrus.Fatal(err)
			}

			// Check if the directory exists!  If not, then we need to create it now
			directory := filepath.Join(currDir, "data")
			_, err = os.Stat(directory + "/" + strings.Split(sourceFilename, "/")[0])
			if err != nil {
				logrus.Info(err)
			}
			if os.IsNotExist(err) {
				logrus.Info("Server: Unable to find directory, '", directory, "'.  Creating now...")
				err := os.MkdirAll(directory+"/"+strings.Split(sourceFilename, "/")[0], 0755)
				if err != nil {
					logrus.Error(err)
				}
			}

			// Write the data into a new file on server's side:
			logrus.Info(directory)
			err = ioutil.WriteFile(filepath.Join(directory, sourceFilename), data, 0600)
			if err != nil {
				logrus.Error("ERROR:", err)
				http.Error(response, err.Error(), http.StatusInternalServerError)
				return
			}
			logrus.Info("Server: File was read from client and written to disk.")
			response.Write([]byte("OK"))
			return
		}
	}
}