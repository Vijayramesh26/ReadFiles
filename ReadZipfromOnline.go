package readfiles

import (
	"archive/zip"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/xuri/excelize/v2"
	// util "github.com/mrlakshmanan/fcsutility"
)

//----------------------------------------------------------- Read ZIP --------------------------------------------------------------

// ReadZip is a function that downloads a ZIP file from a given URL, extracts its contents, and processes supported file types (CSV, TXT, XLSX).

// Step 1: Initialize variables and data structures
// Step 2: Prepare and send an HTTP GET request to the specified URL
// Step 3: Receive the HTTP response
// Step 4: Check for errors in the HTTP request
// Step 5: Create a local ZIP file and write the response body to it
// Step 6: Check for errors in local file creation
// Step 7: Copy the response body to the local ZIP file
// Step 8: Read and process the contents of the ZIP file
// Step 9: Check for errors in ZIP file opening
// Step 10: Loop through files within the ZIP
// Step 11: Read and process supported file types (CSV, TXT, XLSX) within the ZIP
// Step 12: Log and return the result

// Step 1: Initialize variables and data structures
func ReadZip(pUrl string, pFilename string) ([][]string, error) {
	log.Println("ReadZip(+)")

	// Initialize a slice to store the extracted data
	var lFileData [][]string

	// Define the ZIP file name
	lZipFileName := pFilename

	// Create an HTTP client and prepare a GET request to the provided URL
	lClient := http.DefaultClient
	lRequest, lErr := http.NewRequest(http.MethodGet, pUrl, nil)

	// Step 2: Check for errors in request creation
	if lErr != nil {
		return lFileData, fmt.Errorf("ReadZip:001" + lErr.Error())
	}

	// Set HTTP headers for the request
	lRequest.Header.Set("User-Agent", "PostmanRuntime/7.26.10")
	lRequest.Header.Set("Accept", "*/*")
	lRequest.Header.Set("Accept-Encoding", "gzip, deflate, br")
	lRequest.Header.Set("Connection", "keep-alive")

	// Step 3: Send the HTTP request and get the response
	lResponse, lErr := lClient.Do(lRequest)

	// Step 4: Check for errors in the HTTP request
	if lErr != nil {
		return lFileData, fmt.Errorf("ReadZip:002" + lErr.Error())
	}

	// Ensure the response body is closed when done
	defer lResponse.Body.Close()

	// Step 5: Create a local ZIP file and write the response body to it
	out, lErr := os.Create(lZipFileName)

	// Step 6: Check for errors when creating the local ZIP file
	if lErr != nil {
		return lFileData, fmt.Errorf("ReadZip:003" + lErr.Error())
	}

	// Ensure the local ZIP file is closed when done
	defer out.Close()

	// Write the response body to the local ZIP file
	_, lErr = io.Copy(out, lResponse.Body)

	// Step 7: Check for errors when copying the ZIP file
	if lErr != nil {
		return lFileData, fmt.Errorf("ReadZip:004" + lErr.Error())
	}

	// Step 8: Read and process the contents of the ZIP file
	lZipFile, lErr := zip.OpenReader(lZipFileName)

	// Step 9: Check for errors when opening the ZIP file
	if lErr != nil {
		return lFileData, fmt.Errorf("ReadZip:005" + lErr.Error())
	}

	// Ensure the ZIP file is closed when done
	defer lZipFile.Close()

	// Step 10: Loop through files within the ZIP
	for _, lFile := range lZipFile.File {
		if filepath.Ext(lFile.Name) == ".csv" {
			// Read and process the CSV file
			lData, lErr := ReadCsvFromZip(lFile)

			// Step 11: Check for errors when reading CSV files
			if lErr != nil {
				return lFileData, fmt.Errorf("ReadZip:006" + lErr.Error())
			} else {
				// Use conditions to filter your records if needed
				log.Println(lData)
			}
		} else if filepath.Ext(lFile.Name) == ".txt" {
			// Read and process the Text file
			lData, lErr := ReadTextFromZip(lFile)

			// Step 11: Check for errors when reading TXT files
			if lErr != nil {
				return lFileData, fmt.Errorf("ReadZip:006" + lErr.Error())
			} else {
				// Use conditions to filter your records if needed
				log.Println(lData)
			}
		} else if filepath.Ext(lFile.Name) == ".xlsx" {
			// Read and process the XLSX file
			lData, lErr := ReadXlsxFromZip(lFile)

			// Step 11: Check for errors when reading XLSX files
			if lErr != nil {
				return lFileData, fmt.Errorf("ReadZip:006" + lErr.Error())
			} else {
				// Use conditions to filter your records if needed
				log.Println(lData)
			}
		}
	}

	// Step 12: Log and return the result
	log.Println("ReadZip(-)")
	return lFileData, lErr
}

//----------------------------------------------------------------------------------------------------------------------------------

//----------------------------------------------------------- Read CSV --------------------------------------------------------------
// ReadCsvFromZip reads a CSV file from a zip archive and returns the data as a 2D slice of strings.

// Initialize a 2D slice to store the CSV data
// Open the file from the zip archive
// Check for errors while opening the file
// Create a CSV reader for the opened file
// Read the CSV file row by row
// Append each row to the 2D slice
// Return the 2D slice containing the CSV data and no error

func ReadCsvFromZip(file *zip.File) ([][]string, error) {
	// Initialize a 2D slice to store the CSV data
	var lRecord [][]string

	// Open the file from the zip archive
	lFile, lErr := file.Open()
	defer lFile.Close() // Ensure the file is closed when done

	// Check if there was an error opening the file
	if lErr != nil {
		// If there's an error opening the file, return an error with a custom message
		return lRecord, fmt.Errorf("ReadCsvFromZip:001" + lErr.Error())
	} else {
		// Create a CSV reader for the opened file
		lRows := csv.NewReader(lFile)

		// Read the CSV file row by row
		for {
			// Read a row from the CSV
			lRecordRow, lErr := lRows.Read()

			// Check for the end of the file
			if lErr == io.EOF {
				break // Exit the loop when we reach the end of the file
			} else {
				// Append the read row to the 2D slice
				lRecord = append(lRecord, lRecordRow)
			}
		}
	}

	// Return the 2D slice containing the CSV data and no error
	return lRecord, nil
}

//-----------------------------------------------------------------------------------------------------------------------------------

//--------------------------------------------------------- Read Text ---------------------------------------------------------------

// ReadTextFromZip is a function that reads the contents of a Text file stored within a zip archive and returns the data as a 2D slice of strings.

// The function takes a `*zip.File` as input, representing the Text file within the zip archive.
// It returns a 2D slice of strings ([][]string) containing the Text data and an error if any.

func ReadTextFromZip(file *zip.File) ([][]string, error) {
	// Initialize a 2D slice to store the Text data
	var lRecord [][]string

	// Open the file from the zip archive
	lFile, lErr := file.Open()
	defer lFile.Close() // Ensure the file is closed when done

	// Check if there was an error opening the file
	if lErr != nil {
		// If there's an error opening the file, return an error with a custom message
		return lRecord, fmt.Errorf("ReadTextFromZip:001" + lErr.Error())
	} else {
		// Create a CSV reader for the opened file
		lRows := csv.NewReader(lFile)

		// Read the CSV file row by row
		for {
			// Read a row from the CSV
			lRecordRow, lErr := lRows.Read()
			lRows.Comma = '|'
			// Check for the end of the file
			if lErr == io.EOF {
				break // Exit the loop when we reach the end of the file
			} else {
				// Append the read row to the 2D slice
				lRecord = append(lRecord, lRecordRow)
			}
		}
	}

	// Return the 2D slice containing the Text data and no error
	return lRecord, nil
}

//-----------------------------------------------------------------------------------------------------------------------------------

//----------------------------------------------------------- Read XLSX -------------------------------------------------------------
// ReadXlsxFromZip reads an XLSX file from a ZIP archive and extracts its data.
// It takes a zip.File (file) as input and returns a 2D string array containing the XLSX data
// and an error if any occurs.

// Step-by-Step Process:
// 1. Open the file from the ZIP archive.
// 2. Check for errors during file opening.
// 3. If there is an error, return an empty 2D string array and an error with an informative message.
// 4. If the file is opened successfully, initialize an empty 2D string array (lRecord) to store the XLSX data.
// 5. Create a new XLSX file from the opened file.
// 6. Check for errors during XLSX file creation.
// 7. If there is an error, return an empty 2D string array and an error with an informative message.
// 8. Specify the name of the tab in the XLSX file to read (tabName).
// 9. Get all the rows from the specified tab and store them in the 2D string array (lRecord).
// 10. Check for errors during row retrieval.
// 11. If there is an error, return an empty 2D string array and an error with an informative message.
// 12. Iterate through the retrieved rows and append each row to the 2D string array (lRecord).
// 13. Return the populated 2D string array (lRecord) containing the XLSX data and a nil error.

func ReadXlsxFromZip(file *zip.File) ([][]string, error) {
	// Step 1: Open the file from the ZIP archive
	lFile, err := file.Open()
	if err != nil {
		// Step 3: If there is an error, return an empty 2D string array and an error with an informative message
		return nil, fmt.Errorf("ReadXlsxFromZip: Failed to open file: %v", err)
	}
	defer lFile.Close()

	// Step 4: Initialize an empty 2D string array to store the XLSX data
	var lRecord [][]string

	// Step 5: Create a new XLSX file from the opened file
	xlsxFile, err := excelize.OpenReader(lFile)
	if err != nil {
		// Step 7: If there is an error, return an empty 2D string array and an error with an informative message
		return nil, fmt.Errorf("ReadXlsxFromZip: Failed to open XLSX file: %v", err)
	}

	// Step 8: Specify the tab name in the XLSX file you want to read
	tabName := "Sheet1" // Replace with your desired tab name

	// Step 9: Get all the rows from the specified tab and append them to the 2D string array
	rows, err := xlsxFile.GetRows(tabName)
	if err != nil {
		// Step 11: If there is an error, return an empty 2D string array and an error with an informative message
		return nil, fmt.Errorf("ReadXlsxFromZip: Failed to get rows from XLSX file: %v", err)
	}

	// Step 12: Iterate through the retrieved rows and append each row to the 2D string array
	for _, row := range rows {
		lRecord = append(lRecord, row)
	}

	// Step 13: Return the populated 2D string array (lRecord) containing the XLSX data and a nil error
	return lRecord, nil
}

//-----------------------------------------------------------------------------------------------------------------------------------

// func main() {
// 	url := `https://nsearchives.nseindia.com/content/historical/EQUITIES/2023/OCT/cm06OCT2023bhav.csv.zip`
// 	fmt.Println(ReadZip(url))
// 	fmt.Println("Program End")
// }
