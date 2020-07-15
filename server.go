package main

import (
	"flag"
	"log"
	"net"
	
///////////////////
	//"bufio"
	"encoding/json"
	"encoding/csv"
	"fmt"
	"io"
	//"log"
	"os"
	"strconv"
	//"reflect"

)
/////////////////////////ReadMe/////////////////////////////////////////////////////
////////(functions.go) file contains functions we used in server.go for testing separately 
//////////////(TimeSeries_KeyIndicators.csv) contains data set for query
//////////////working of each function is explained in their space 
//////////idea for writing tcp server is taken from provided yotube-tcp-video 
////////////i was unable to perform proper testing of server despite using four different installatons of Netcat on windows platform
/////////////so i tested my functions in a separate file (functions.go) they are working fine and 
////////////////are providing desired output as two screenshots of the output are attached in repository
//////////////////////////////////////////////////////////////////////////////////////////
func main() {
	var addr string
	var network string
	flag.StringVar(&addr, "e", ":4040", "service endpoint [ip addr or socket path]")
	flag.StringVar(&network, "n", "tcp", "network protocol [tcp,unix]")
	flag.Parse()

	// validate supported network protocols
	switch network {
	case "tcp", "tcp4", "tcp6", "unix":
	default:
		log.Fatalln("unsupported network protocol:", network)
	}

	// create a listener for provided network and host address
	ln, err := net.Listen(network, addr)
	if err != nil {
		log.Fatal("failed to create listener:", err)
	}
	defer ln.Close()
	log.Println("**** COVID-19 ***")
	log.Printf("Service started: (%s) %s\n", network, addr)

	// connection-loop - handle incoming requests
	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println(err)
			if err := conn.Close(); err != nil {
				log.Println("failed to close listener:", err)
			}
			continue
		}
		log.Println("Connected to", conn.RemoteAddr())

		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer func() {
		if err := conn.Close(); err != nil {
			log.Println("error closing connection:", err)
		}
	}()

	if _, err := conn.Write([]byte("Connected...\n")); err != nil {
		log.Println("error writing:", err)
		return
	}

	// loop to stay connected with client until client breaks connection
	for {
		///////////////////////////////buffer start////////////////
		// buffer for client command
		cmdLine := make([]byte, 0, 4096) // big buffer to store all input data
   		 tmp := make([]byte, 256)     // using small tmo buffer for reading chunks of input bytes
		for {
			n, err := conn.Read(tmp)
			if err != nil {
				if err != io.EOF {
					fmt.Println("read error:", err)
					return
				}
				break
			}
			//fmt.Println("got", n, "bytes.")
			cmdLine = append(cmdLine, tmp[:n]...)
	
		}
		//////////////////////////////////buffer end///////////////
		var key=findKey(string(cmdLine))
		var result=LoadAndSearch("TimeSeries_KeyIndicators.csv",key)
		var count=len(result)
			if count == 0 {
				if _, err := conn.Write([]byte("Nothing found\n")); err != nil {
					log.Println("failed to write:", err)
				}
				continue
			}
			if count > 0 {
				
			  	for i := 0; i < count; i++{
					  ///convert my each result data to byte to transmit over connection
					fmt.Println(result[i])
					  data := []byte(result[i])		
			 	if _, err := conn.Write(data); 
			  	err != nil {
			  		log.Println("failed to write response:", err)
			  	return
			 	}
			   }
			}
		
	}
}

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
				 
				 return key
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
		var f=string(JsonData)
		//add Json data to map
		table = append(table, f)
		
					}
					
				}
								
		}
	
		}

		//return final calculated result - array of json-encoded-arrays
		return table
		//fmt.Println(table)		
		
	}