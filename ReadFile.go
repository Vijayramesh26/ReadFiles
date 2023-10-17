package readfiles

import (
	"encoding/csv"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"strings"

	"github.com/xuri/excelize/v2"
)

// GetEnhanceUploadedFile retrieves an uploaded file from an HTTP request.
// It takes an HTTP request (r) and the name of the form field containing the file (formName) as input.
// The function returns a strings.Reader containing the file's content, the filename, the multipart.FileHeader,
// and an error if any occurs.

// Step-by-Step Process:
// 1. Initialize variables: fileStr (to store the file's content) and file (a strings.Reader to hold the file content).
// 2. Log the start of the function.
// 3. Attempt to retrieve the file data and header using r.FormFile(formName).
// 4. If an error occurs during retrieval, log the error and return an empty file, empty fileStr, the header, and the error.
// 5. If the file data is successfully retrieved, read its content into fileStr.
// 6. Create a strings.Reader (file) from fileStr to facilitate further use of the file's content.
// 7. Log the end of the function.
// 8. Return the file, fileStr, the multipart.FileHeader, and nil as the error.
func GetFileDetails(r *http.Request, formName string) (*strings.Reader, string, *multipart.FileHeader, error) {
	log.Println("GetFileDetails(+)")

	fileStr := ""
	var file *strings.Reader

	// Attempt to retrieve the file data and header using r.FormFile(formName)
	fileBody, header, lErr := r.FormFile(formName)

	if lErr != nil {
		// If an error occurs during retrieval, return an empty file, empty fileStr, the header, and the error
		return file, fileStr, header, fmt.Errorf("GetFileDetails:001" + lErr.Error())
	} else {
		// If the file data is successfully retrieved, read its content into fileStr
		datas, _ := ioutil.ReadAll(fileBody)
		fileStr = string(datas)

		// Create a strings.Reader (file) from fileStr to facilitate further use of the file's content
		file = strings.NewReader(fileStr)

		log.Println("GetFileDetails(-)")
		// Log the end of the function
		return file, fileStr, header, nil
	}
}

// ReadCSV reads the contents of a CSV file from a strings.Reader and returns the data as a 2D slice of strings.

// Step 1: Initialize a 2D slice to store the CSV data
// Step 2: Create a CSV reader for the input string reader
// Step 3: Read the CSV file row by row
// Step 4: Check for the end of the file
// Step 5: Append each row to the 2D slice
// Step 6: Return the 2D slice containing the CSV data and no error
func ReadCSV(r *http.Request, pFile string) ([][]string, error) {
	var lRecord [][]string
	lFile, _, _, lErr := GetFileDetails(r, pFile)
	if lErr != nil {
		return lRecord, fmt.Errorf("ReadCSV:001" + lErr.Error())
	} else {
		// Step 2: Create a CSV reader for the input string reader
		lRows := csv.NewReader(lFile)

		// Step 3: Read the CSV file row by row
		for {
			// Step 4: Read a row from the CSV
			lRecordRow, lErr := lRows.Read()

			// Step 4: Check for the end of the file
			if lErr == io.EOF {
				break // Exit the loop when we reach the end of the file
			} else {
				// Step 5: Append the read row to the 2D slice
				lRecord = append(lRecord, lRecordRow)
			}
		}
	}
	// Step 6: Return the 2D slice containing the CSV data and no error
	return lRecord, nil
}

// ReadText reads the contents of a text file from a strings.Reader and returns the data as a 2D slice of strings.

// Step 1: Initialize a 2D slice to store the text data
// Step 2: Create a CSV reader for the input string reader (assuming it's CSV-formatted text)
// Step 3: Read the text file row by row
// Step 4: Check for the end of the file
// Step 5: Append each row to the 2D slice
// Step 6: Return the 2D slice containing the text data and no error
func ReadText(r *http.Request, pFile string) ([][]string, error) {
	var lRecord [][]string
	lFile, _, _, lErr := GetFileDetails(r, pFile)
	if lErr != nil {
		return lRecord, fmt.Errorf("ReadText:001" + lErr.Error())
	} else {
		// Step 2: Create a CSV reader for the input string reader (assuming it's CSV-formatted text)
		lRows := csv.NewReader(lFile)
		lRows.Comma = '|'
		// Step 3: Read the text file row by row
		for {
			// Step 4: Read a row from the text
			lRecordRow, lErr := lRows.Read()

			// Step 4: Check for the end of the file
			if lErr == io.EOF {
				break // Exit the loop when we reach the end of the file
			} else {
				// Step 5: Append the read row to the 2D slice
				lRecord = append(lRecord, lRecordRow)
			}
		}
	}
	// Step 6: Return the 2D slice containing the text data and no error
	return lRecord, nil
}

// ReadXlsxFile reads an uploaded XLSX file from an HTTP request, extracts its contents,
// and performs specific operations on the data.
// It takes an HTTP request (r), the name of the uploaded file (pFile) as inputs.
// The function first retrieves the uploaded file and saves it to the server. It then opens
// the XLSX file, reads the data from a specified tab, and stores it in a 2D array (record).
// After processing the data, the function removes the temporary file created during the operation.
// Any encountered errors are returned as error messages.

// File Creation and Reading:(Step-by-Step Process)
// 1. The uploaded file is retrieved from the HTTP request.
// 2. It is saved to a server location with a unique file name based on the provided Header.Filename.
// 3. The content of the uploaded file is read, and a temporary file is created on the server.
// 4. The XLSX file is opened using excelize, and data is extracted from a specified tab (e.g., "TabName").
// 5. The data is stored in a 2D array (record) for further processing.

func ReadXlsxFile(r *http.Request, pFile string) error {
	log.Println("ReadXlsFile +")
	var record [][]string
	lFile, _, Header, lErr := GetFileDetails(r, pFile)
	if lErr != nil {
		return fmt.Errorf("ReadXlsFile:001" + lErr.Error())
	} else {
		path := "./"
		fileName := path + Header.Filename
		//Creating file's in specific server path
		out, lErr := os.Create(fileName)
		if lErr != nil {
			return fmt.Errorf("ReadXlsFile : 002" + lErr.Error())
		} else {
			datas, _ := ioutil.ReadAll(lFile)
			fileStr := string(datas)

			file := strings.NewReader(fileStr)

			_, lErr = io.Copy(out, file) // Copy the file's content to the server file
			if lErr != nil {
				return fmt.Errorf("ReadXlsFile : 003" + lErr.Error())
			} else {
				lNewFile, lErr := excelize.OpenFile(fileName)
				if lErr != nil {
					return fmt.Errorf("ReadXlsFile : 004" + lErr.Error())
				} else {
					rows, lErr := lNewFile.GetRows("TabName") // Specify the tab name in the XLSX file
					if lErr != nil {
						return fmt.Errorf("ReadXlsFile : 005" + lErr.Error())
					} else {
						for _, row := range rows {
							record = append(record, row)
						}
					}

					// Your condition to filter records can be applied here.

					// This method is used to remove the temporary file created during the operation.
					lErr = os.Remove(fileName)
					if lErr != nil {
						return fmt.Errorf("ReadXlsFile : 005-" + lErr.Error())
					}
				}
			}
		}
	}
	log.Println("ReadXlsFile-")
	return nil
}

// Join2DArray joins two 2D arrays by appending the rows of the second array to the first array.
// It takes two 2D arrays (array1 and array2) as input and returns a new 2D array.
// The function iterates through the rows of array2 and appends each row to the end of array1.
// The resulting combined array is returned as the output.

// Step-by-Step Process:
// 1. Iterate through the rows of array2.
// 2. For each row in array2, create a copy (lRow) by appending it to array1.
// 3. Append the copied row (lRow) to array1, effectively combining the two arrays.
// 4. Repeat the process for all rows in array2.
// 5. The final combined array is returned as the output.

func Join2DArray(pArray1 [][]string, pArray2 [][]string) [][]string {
	// Iterate through the rows of array2
	for ArrayIndex := 0; ArrayIndex < len(pArray2); ArrayIndex++ {
		// Append the current lRow of array2 to array1
		lRow := append(pArray2[ArrayIndex])
		pArray1 = append(pArray1, lRow)
	}
	return pArray1
}

// Filter2DArray1 filters a 2D array by searching for a specified start value in the first column.
// It takes a start value (pStartvalue) and a 2D array (pRows) as input.
// The function iterates through the rows of the input array to find the first occurrence
// of the start value in the first column. It then creates a new array that includes all
// rows following the matching row until it encounters a row with an empty or null first column.
// The resulting filtered array is returned as the output.

// Step-by-Step Process:
// 1. Initialize variables: lIndex (to store the index of the matching row) and lNewArr (the filtered array).
// 2. Iterate through the rows of pRows to find the first occurrence of pStartvalue in the first column.
// 3. When a match is found, store the index in lIndex.
// 4. Iterate through the rows of pRows starting from the matching row (lIndex).
// 5. For each row, check if it has data in the first column and if it's not empty or null.
// 6. If the conditions are met, append the row to lNewArr.
// 7. Continue this process until an empty or null value is encountered in the first column.
// 8. Return the filtered array (lNewArr) as the output.

func Filter2DArray1(pStartvalue string, pRows [][]string) [][]string {
	log.Println("Filter2DArray +")
	var lIndex int
	var lNewArr [][]string

	// Iterate through the rows of pRows to find the first occurrence of pStartvalue
	for i := 0; i < len(pRows); i++ {
		if len(pRows[i]) > 0 {
			if strings.Contains(strings.ToLower(pRows[i][0]), strings.ToLower(pStartvalue)) {
				lIndex = i
			}
		}
	}

	// Iterate through the rows of pRows to store records following the matching row
	// until an empty or null value is encountered in the first column.
	for j := lIndex; j < len(pRows); j++ {
		if len(pRows[j]) != 0 {
			if pRows[j][0] != "" {
				lNewArr = append(lNewArr, pRows[j])
			} else {
				break
			}
		}
	}

	log.Println("Filter2DArray - ")
	return lNewArr
}

// Filter2DArray2 filters a 2D array by searching for a specified start value in the first column.
// It takes a start value (pStartvalue) and a 2D array (pRows) as input.
// The function iterates through the rows of the input array to find all occurrences
// of the start value in the first column. It then creates a new array that includes all
// rows following each matching row until it encounters a row with an empty or null first column.
// The resulting filtered array is returned as the output.

// Step-by-Step Process:
// 1. Initialize variables: lIndex (to store indices of matching rows) and lNewArr (the filtered array).
// 2. Iterate through the rows of pRows to find all occurrences of pStartvalue in the first column.
// 3. When a match is found, store the index in lIndex.
// 4. Iterate through the indices in lIndex.
// 5. For each matching row, iterate through the rows of pRows to store records following that row.
// 6. Continue storing rows until an empty or null value is encountered in the first column.
// 7. Repeat this process for all matching rows.
// 8. Return the filtered array (lNewArr) as the output.

func Filter2DArray2(pStartvalue string, pRows [][]string) [][]string {
	log.Println("Filter2DArray +")
	var lIndex []int
	var lNewArr [][]string

	// Iterate through the rows of pRows to find all occurrences of pStartvalue
	for i := 0; i < len(pRows); i++ {
		if len(pRows[i]) > 0 {
			if strings.Contains(strings.ToLower(pRows[i][0]), strings.ToLower(pStartvalue)) {
				lIndex = append(lIndex, i)
			}
		}
	}

	// Iterate through the indices in lIndex
	for srIndex := 0; srIndex < len(lIndex); srIndex++ {
		// Iterate through the rows of pRows to store records following each matching row
		// until an empty or null value is encountered in the first column.
		for j := lIndex[srIndex] + 1; j < len(pRows); j++ {
			if len(pRows[j]) != 0 {
				if pRows[j][0] != "" {
					lNewArr = append(lNewArr, pRows[j])
				} else {
					break
				}
			}
		}
	}

	log.Println("Filter2DArray - ")
	return lNewArr
}
