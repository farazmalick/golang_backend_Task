	package main

	import (
		//"bufio"
		"encoding/json"
		"encoding/csv"
		"fmt"
		"io"
		//"log"
		"os"
		"strconv"
		//"reflect"
		//"bytes"
	//"encoding/gob"
	)

	func main() {
		var a =`{
			"query": {
				"date":"2020-5-25"
			}
		}`
		
		var key=findKey(a)
		var result=LoadAndSearch("TimeSeries_KeyIndicators.csv",key)
		var count=len(result)
		//////////////////
		
		// buf := &bytes.Buffer{}
		// gob.NewEncoder(buf).Encode(result)
		// bs := buf.Bytes()
		// //fmt.Printf("%s", bs)
		// var x=len(bs)
		// fmt.Println(x)
		// for i := 0; i <x; i++{ 
		// 		 fmt.Printf("%s", bs[i])
		// 	}
		///////////////
		if (count<=0){
		fmt.Println(result)}
		if (count>0){
			//fmt.Fprint(w, result)
			// fmt.Sprintf(
		 	// 	"%s",
		 	// 	result,
		    // )
		//fmt.Printf("%v", result)
		 for i := 0; i < count; i++{ 
		// 	//fmt.Println(result[i])
		//data := []byte(result[i])
			 //fmt.Println(data)
			 fmt.Println(result[i])
			 
			}
		}
	}
	/////store:for saving file data
	type Data struct
		{
			Date string `json:"date"`
			Positive int `json:"positive"`
			Tests int `json:"tests"`
			Discharged int `json:"discharged"`
			Expired int `json:"expired"`
			Admitted int `json:"admitted"`
			Region string `json:"region"`
			
		}
		//type Data []DataItem

		//function to find key for searching data from file
	func findKey(value string)interface{}{
			//fmt.Println(reflect.TypeOf(value))
			byt := []byte(value)
			//fmt.Println(reflect.TypeOf(value))
			var dat map[string]interface{}
			if err := json.Unmarshal([]byte(byt), &dat); err != nil {
				panic(err)
			}
			var key=dat["query"].(map[string]interface{})["region"]
			//fmt.Println(key)
			
			 if key == nil{ 
				 key=dat["query"].(map[string]interface{})["date"]	}

				 var X = fmt.Sprintf("%v", key)
				 return X
		}
		
	//Load file and perform search
	func LoadAndSearch(filePath string,key interface{})[]string  {
		
		// Load  file.
		
		f, _ := os.Open(filePath)

		// Create a new reader.
		
		r := csv.NewReader(f)
		//create map using make
		table := make([]string, 0)
		
		for {
			record, err := r.Read()
			// terminate loop at end-of-file
			if err == io.EOF {
				break
			}

			if err != nil {
				panic(err)
			}
			
			//fmt.Println(record)
		// fmt.Println(len(record))
		
		
			for value := range record {
				
				if (record[value]==key) {
					//fmt.Printf("  %v\n", record[value])
					var a =record
					
		//perform conversion string->integer		 
					i1, err := strconv.Atoi(a[0])
					i2, err :=strconv.Atoi(a[1])
					i3, err :=strconv.Atoi(a[3])
					i4, err :=strconv.Atoi(a[4])
					i5, err :=strconv.Atoi(a[5])
					
					if err == nil {
					result := Data{
						Date :a[2],
						Positive :i1,
						Tests:i2,
						Discharged:i3, 
						Expired :i4,
						Admitted :i5,
						Region :a[6],
					}
					
		JsonData, _ := json.MarshalIndent(result, "", "    ")
		//fmt.Println(JsonData)
		var f=string(JsonData)
		//add Json data to map
		table = append(table, f)
		
					}
					
				}
								
		}
	
		}

		//return final calculated result as an array of json-encoded-data
		return table
		//fmt.Println(table)		
		
	}