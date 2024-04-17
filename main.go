package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/Jeffail/gabs/v2"
	"github.com/fatih/color"
	"github.com/rodaine/table"
)

var (
	// # Define the flags
	port      = flag.String("port", "31000", "Port to listen on")
	output    = flag.String("output", "output.jsonl", "Output file")
	verbose   = flag.Bool("verbose", false, "Print the data to the console")
	pretty    = flag.Bool("pretty", false, "Pretty print the JSON data as a table")
	doctor    = flag.Bool("doctor", false, "Run the doctor to check the health of the system")
	logEvents = []gabs.Container{}
	gabsData  *gabs.Container

	successMessage = color.New(color.FgGreen).SprintFunc()
	errorMessage   = color.New(color.FgRed).SprintFunc()
	warningMessage = color.New(color.FgYellow).SprintFunc()
	titleMessage   = color.New(color.FgHiWhite).SprintFunc()
)

func init() {
	flag.Parse()
	log.Println("Starting Cloudflow UDP server on port", *port)
	fmt.Println(titleMessage(motd_banner))

	if *doctor {
		// check for AWS CSM Mode variables
		csmEnabledEnv := os.Getenv("AWS_CSM_ENABLED")
		csmHost := os.Getenv("AWS_CSM_HOST")
		csmPort := os.Getenv("AWS_CSM_PORT")

		if csmEnabledEnv == "" {
			// color error messages
			log.Fatal(errorMessage("Fatal Error: AWS_CSM_ENABLED is not set"))
			if csmHost == "" || csmPort == "" {
				log.Println(warningMessage("Warning: AWS_CSM_HOST or AWS_CSM_PORT is not set"))
			}
		} else {
			log.Println(successMessage("AWS_CSM_ENABLED is set to ", csmEnabledEnv))
			log.Println(successMessage("AWS_CSM_HOST is set to ", csmHost))
			log.Println(successMessage("AWS_CSM_PORT is set to ", csmPort))
		}

		os.Exit(0)
	}

	log.Println("Writing output to", *output)
	log.Println("Waiting for AWS API events...")
}

// manage shutdown signal to gracefully close the server
func shutdownSignal() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
	log.Println("Shutting down...")
	// check file for events
	log.Printf("Check %s for log events\n", *output)
	os.Exit(0)
}

func printEventTable() {
	headerFmt := color.New(color.FgGreen, color.Underline).SprintfFunc()
	columnFmt := color.New(color.FgYellow).SprintfFunc()

	tbl := table.New("Timestamp", "AccessKey", "Service", "Api", "Region", "UserAgent")
	tbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)
	if len(logEvents) == 0 {
		return
	}
	if logEvents[0].Path("AccessKey").Data() == nil {
		return
	}
	for _, event := range logEvents {
		timestamp := int(event.Path("Timestamp").Data().(float64))
		accessKey := event.Path("AccessKey").Data()
		if accessKey == nil {
			continue
		}
		service := event.Path("Service").Data().(string)
		api := event.Path("Api").Data().(string)
		region := event.Path("Region").Data().(string)
		userAgent := event.Path("UserAgent").Data().(string)
		tbl.AddRow(timestamp, accessKey, service, api, region, userAgent)
	}
	// clear the screen
	fmt.Print("\033[H\033[2J")
	tbl.Print()

}

func main() {
	// # Handle shutdown signals
	go shutdownSignal()
	// # Open the output file
	file, err := os.OpenFile(*output, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}

	// # Close the file when the program exits
	defer file.Close()

	// # Create a new UDP listener
	addr, err := net.ResolveUDPAddr("udp", ":"+*port)
	if err != nil {
		log.Fatal(err)
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		log.Fatal(err)
	}

	// # Read from the UDP connection
	for {
		buffer := make([]byte, 2048)
		n,
			_, err := conn.ReadFromUDP(buffer)
		if err != nil {
			log.Fatal(err)
		}
		// # Print the data to the console
		if *verbose {
			log.Println(string(buffer[:n]))
		} else {

			if !*pretty {
				fmt.Println(string(buffer[:n]))
			}
		}
		// # Write the data to the output file
		data := strings.TrimSpace(string(buffer[:n]))
		if data != "" {
			gabsData, err = gabs.ParseJSON([]byte(data))
			if err != nil {
				log.Fatal(err)
			}
			// # Remove the session token from the data
			gabsData.Delete("SessionToken")
			_, err = file.WriteString(gabsData.String() + "\n")

			if err != nil {
				log.Fatal(err)
			}
		}
		// append the data to the logEvents slice for later processing
		if *pretty {
			printEventTable()
		}
		logEvents = append(logEvents, *gabsData)

	}
}
