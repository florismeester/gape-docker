package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	
)

/* Easy to use recursive directory filesystem changes notification tool that can send output to
   remote or local syslog servers, based on github.com/rjeczalik/notify
   Floris Meester floris@grid6.io This version is specific for Docker containers
*/



// Get information from the Docker engine and build the union mount paths
func getmountpoints() []string {
	
	// Create a client 
	cli, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}

	// List all the running containers
	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil {
		panic(err)
	}

        var info types.Info
	
	// Get an info struct from the engine to fetch the driver type and the Docker root dir
	info, err = cli.Info(context.Background())
	rootdir := info.DockerRootDir
	mpath = rootdir + "/" + info.Driver + "/mnt"
	var containerlist []string		
	
	// Loop through the container list and build the paths
	for _, container := range containers {
		opath := rootdir + "/image/" + info.Driver + "/layerdb/mounts/" +  container.ID + "/mount-id"
		fcontent, err := ioutil.ReadFile(opath)	
		if err != nil {
			fmt.Println(err)
		}
		path := rootdir + "/" + info.Driver + "/mnt/" + string(fcontent)
		containerlist = append(containerlist, path)
		
	}
	return containerlist

}


