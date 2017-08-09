package main

import (
	"fmt"
	"net/http"

	//"io/ioutil"
	//"reflect"
	"io/ioutil"
	//"reflect"
	"log"
	"strings"
	//"os"
	//"io"
	"io"
	"os"
)

type findObject struct {
	initialCapture string
	prefix         string
}

var referenceObjects []findObject

func extractAndSaveObject(downloadFolder string, objectLocation string) {

	// don't worry about errors
	response, e := http.Get(objectLocation)
	if e != nil {
		log.Fatal(e)
	}

	fmt.Println("In extractAndSaveObject - " + objectLocation)
	defer response.Body.Close()

	filename := objectLocation[strings.LastIndex(objectLocation, "/")+1:]
	fmt.Println("Last location of / is - ", strings.LastIndex(objectLocation, "/")+1)
	fmt.Println("filename is - ", filename)

	//open a file for writing
	file, err := os.Create(downloadFolder + filename)
	if err != nil {
		log.Fatal(err)
	}
	// Use io.Copy to just dump the response body to the file. This supports huge files
	_, err = io.Copy(file, response.Body)
	if err != nil {
		log.Fatal(err)
	}
	file.Close()

	fmt.Println("Success!")
}

//
func lookForNextObjectTag(searchableData string) (foundObject int) {

	//Initialise to -1.
	//If find an object then store the object index.
	var startOfObjectLocation int
	//var passBackObject findObject
	minimumStartOfObjectLocation := -1

	//fmt.Println(searchableData)
	for _, referenceObject := range referenceObjects{
		startOfObjectLocation = strings.Index(searchableData, referenceObject.initialCapture)
		if startOfObjectLocation != -1 {
			if startOfObjectLocation < minimumStartOfObjectLocation || minimumStartOfObjectLocation == -1 {
				//passBackObject = referenceObject
				minimumStartOfObjectLocation = startOfObjectLocation
			}
		}
	}
	return minimumStartOfObjectLocation
}

func extractImages(body string, downloadFolder string) {
	var startOfObjectLocation int
	var endOfObjectLocation int
	var endOfObjectLocationSingleQuote int
	var endOfObjectLocationDoubleQuote int

	startOfObjectLocation = lookForNextObjectTag(body)

	for startOfObjectLocation != -1 {
		body = body[startOfObjectLocation:]
		fmt.Println("yoh! - ", startOfObjectLocation)
		endOfObjectLocationSingleQuote = strings.Index(body, "'")
		endOfObjectLocationDoubleQuote = strings.Index(body, "\"")

		//Single quotes and double quotes can be used in html.
		//Work out which comes soonest, and this will be the same as the
		//one at the start of the string.
		if endOfObjectLocationSingleQuote < endOfObjectLocationDoubleQuote {
			endOfObjectLocation = endOfObjectLocationSingleQuote
		} else{
			endOfObjectLocation = endOfObjectLocationDoubleQuote
		}

		fmt.Println("endOfObjectLocation = ", endOfObjectLocation)
		fmt.Println("Object location is " + "https://" + body[0:endOfObjectLocation])
		body = body[endOfObjectLocation:]
/*		fmt.Println("Left over body = ", manipulatedBody)*/

		fmt.Println("startOfObjectLocation = ", startOfObjectLocation)
		startOfObjectLocation = lookForNextObjectTag(body)

	}

}

func initialiseObjectReferences() {

	//Populate the slice.
	referenceObjects = append(referenceObjects, findObject{initialCapture: "/img.php?loc=loc", prefix:"https://"})
	referenceObjects = append(referenceObjects, findObject{initialCapture: "imagevenue.com/img.php?", prefix:"https://"})
	referenceObjects = append(referenceObjects, findObject{initialCapture: "hifiwigwam.com/"})
	referenceObjects = append(referenceObjects, findObject{initialCapture: "picbux.com/image.php?id=", prefix:"https://"})
	referenceObjects = append(referenceObjects, findObject{initialCapture: "picturesupload.com/show.php/", prefix:"https://"})
	referenceObjects = append(referenceObjects, findObject{initialCapture: "imagehigh.com/", prefix:"https://"})
	referenceObjects = append(referenceObjects, findObject{initialCapture: "image2share.com/", prefix:"https://"})
	referenceObjects = append(referenceObjects, findObject{initialCapture: "paintedover.com/uploads/show.php", prefix:"https://"})
	referenceObjects = append(referenceObjects, findObject{initialCapture: "10pix.com/", prefix:"https://"})
}

func main() {

	downloadFolder := "./"

	initialiseObjectReferences()

	//resp, err := http.Get("http://www.bbc.co.uk/")
	resp, err := http.Get("https://hifiwigwam.com/forum/topic/125807-what-does-%C2%A33900-buy-you/")

	if err != nil {
		// handle error
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if err == nil {
		searchableBody := string(body)

		extractImages(searchableBody, downloadFolder)
	}
	/*
		ioutil.WriteFile("dump", body, 0600)

		for i:= 0; i < len(body); i++ {
			fmt.Println( body[i] ) // This logs uint8 and prints numbers
		}

		fmt.Println( reflect.TypeOf(body) )
		fmt.Println("done")
	*/
}
