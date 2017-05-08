package main


/* Easy to use recursive directory filesystem changes notification tool that can send output to 
   remote or local syslog servers, based on github.com/rjeczalik/notify
   Floris Meester floris@grid6.io This version is specific for Docker containers
*/


import (
    "log"
    "encoding/json"
    "github.com/rjeczalik/notify"
    "flag"
    "os"
    "fmt"
    "strings"
    "log/syslog"
    "time"
)

type Configuration struct {
    Sysloghost string
    Syslogproto string
    Syslogport string
    Stdout bool
    Localonly bool
    Paths []string
}

// Ugly but we create a global variable holding the Docker container mounts

var mpath string
var c chan notify.EventInfo

func main(){

        // Open configuration file
        conf := flag.String("config","/etc/gape-docker.conf", "Path to configuration file")
        flag.Parse()
        file, err := os.Open(*conf)
        if err != nil {
                log.Fatal("Can't find configuration file, try 'gape-docker -config <path> ", *conf)
        }
        decoder := json.NewDecoder(file)
        configuration := Configuration{}
        err = decoder.Decode(&configuration)
        if err != nil {
                fmt.Println("error opening configuration:", err)
        }

	// Create the notification channel
	c = make(chan notify.EventInfo, 1)
	cc := make(chan notify.EventInfo, 1)

        // Create a syslog writer for logging local or remote
	if configuration.Localonly{
	        logger, err := syslog.New(syslog.LOG_NOTICE, "Gape-docker")
                if err == nil {
                	log.SetOutput(logger)
		} else {
			log.Fatal(err)
		}
        }else {
        	logger, err := syslog.Dial(configuration.Syslogproto, configuration.Sysloghost +
			":" + configuration.Syslogport, syslog.LOG_NOTICE, "Gape")
                	if err == nil {
                	log.SetOutput(logger)
        	}else{
			log.Fatal(err)
		}	
	}
	
	// Try to get the container mountpoints
	mpoints := getmountpoints()
	
	// Create the notification watches
	for _, mp := range mpoints {
		createwatches(mp, configuration)
	}

	// We create a special non-recursive watch for the mount point path so we can detect the creation of new containers
	if err := notify.Watch(mpath, cc, notify.Create); err != nil {
    		log.Fatal(err)
	}
	defer notify.Stop(c)

	// Loop forever and receive  events from the channels.
	var ei notify.EventInfo

	for {
		select {
			case ei = <-c:
				log.Print(ei)
			case nc := <- cc:
				if !strings.HasSuffix(nc.Path(), "-init") && !strings.HasSuffix(nc.Path(), "-removing") {
				fmt.Println("New container detected", nc.Path())

				// give the container some time
				time.Sleep(3 * time.Second)
				createwatches(nc.Path(), configuration)
			}
		}
                if configuration.Stdout {
                        fmt.Println(ei)
                }
	}
}

// Setup watches 
func createwatches(mp string, configuration Configuration){
	for _,item := range configuration.Paths {

                // Check if item exists and is a directory otherwise bailout
                fd, err := os.Stat(item)
                if err != nil {
                        log.Fatal(err)
                }
                if !fd.IsDir(){
                        log.Fatal("Not a directory: ", item)
                }
                // I might move the notification options to the config file
                // The main difference with standard gape is that we build the paths within
                // the container mountpoints
                path := mp + item
                fmt.Println(path)
                if err := notify.Watch(path + "...", c, notify.Remove, notify.Create, notify.Write, notify.Rename ); err != nil {
                        log.Fatal(err)
                }
        }
}
