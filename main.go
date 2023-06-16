package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/romanzac/json-mp/mp"
	"github.com/romanzac/json-mp/shape"
	"github.com/spf13/cobra"
	"io"
	"os"
)

var (
	isDecoding            bool
	inputFile, outputFile string

	// JsonMpCmd to starts the application
	JsonMpCmd = &cobra.Command{
		Use:     "json-mp",
		Short:   "JSON <-> MessagePack encoding tool",
		Long:    `Encodes file from JSON to MessagePack format (default mode)`,
		Version: "1.0.0",
		Run:     runJsonMp,
	}
)

func init() {
	JsonMpCmd.PersistentFlags().BoolVarP(&isDecoding, "decode", "d", false, "decodes MessagePack to JSON format")
	JsonMpCmd.PersistentFlags().StringVarP(&inputFile, "input", "i", "", "input file path")
	JsonMpCmd.PersistentFlags().StringVarP(&outputFile, "output", "o", "", "output file path")
	JsonMpCmd.MarkPersistentFlagRequired("input")
	JsonMpCmd.MarkPersistentFlagRequired("output")
	JsonMpCmd.MarkFlagsRequiredTogether("input", "output")

}

func main() {
	err := JsonMpCmd.Execute()
	if err != nil && err.Error() != "" {
		fmt.Println(err)
	}
}

func encodeMessagePack(fileIn *os.File) ([]byte, error) {

	inStat, err := fileIn.Stat()
	if err != nil {
		return nil, err
	}

	// Read the file into a byte slice
	dataIn := make([]byte, inStat.Size())
	_, err = bufio.NewReader(fileIn).Read(dataIn)
	if err != nil && err != io.EOF {
		return nil, err
	}

	// Validate json input
	if !json.Valid(dataIn) {
		return nil, errors.New("invalid JSON input")
	}

	// Assign data shape and deserialize JSON
	var result shape.DataShape
	err = json.Unmarshal(dataIn, &result)
	if err != nil {
		return nil, err
	}

	// Encode to MessagePack
	dataOut, err := mp.Marshal(result)
	if err != nil {
		return nil, err
	}

	return dataOut, nil
}

func decodeMessagePack(fileIn *os.File) ([]byte, error) {

	// Get the file size
	inStat, err := fileIn.Stat()
	if err != nil {
		return nil, err
	}

	// Read the file into a byte slice
	dataIn := make([]byte, inStat.Size())
	_, err = bufio.NewReader(fileIn).Read(dataIn)
	if err != nil && err != io.EOF {
		return nil, err
	}

	// Assign data shape and deserialize MessagePack
	result := shape.DataShape{}
	err = mp.Unmarshal(dataIn, &result)
	if err != nil {
		return nil, err
	}

	// Encode to JSON
	dataOut, err := json.Marshal(result)
	if err != nil {
		return nil, err
	}

	return dataOut, nil
}

func runJsonMp(cmd *cobra.Command, args []string) {

	fileIn, err := os.Open(inputFile)
	if err != nil {
		panic(err)
	}
	defer fileIn.Close()

	if !isDecoding {
		mpData, err := encodeMessagePack(fileIn)
		if err != nil {
			fmt.Printf("Error during encoding to MessagePack: %v", err)
			return
		}
		if err = os.WriteFile(outputFile, mpData, 0666); err != nil {
			fmt.Printf("Error during writing the MessagePack file: %v", err)
			return
		}

	} else {
		jsonData, err := decodeMessagePack(fileIn)
		if err != nil {
			fmt.Printf("Error during decoding to JSON: %v", err)
			return
		}
		if err = os.WriteFile(outputFile, jsonData, 0666); err != nil {
			fmt.Printf("Error during writing the JSON file: %v", err)
			return
		}
	}
}
