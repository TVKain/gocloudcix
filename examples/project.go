package main

import (
	"context"
	"fmt"

	"github.com/TVKain/go-cloudcix/cloudcix"
	"github.com/TVKain/go-cloudcix/cloudcix/compute/projects"
)

func main() {
	opts, err := cloudcix.AuthOptionsFromEnv()

	if err != nil {
		fmt.Printf("Error reading auth options: %v\n", err)
		return
	}

	cloudcixClient, err := cloudcix.AuthenticatedClient(context.TODO(), opts)

	if err != nil {
		fmt.Printf("Error authenticating: %v\n", err)
		return
	}
	computeClient, err := cloudcix.NewCompute(context.TODO(), cloudcixClient)
	if err != nil {
		fmt.Printf("Error creating compute client: %v\n", err)
		return
	}

	project, err := projects.Get(context.TODO(), computeClient, 200)
	if err != nil {
		fmt.Printf("Error getting project: %v\n", err)
	}

	if project != nil {
		fmt.Printf("Project ID: %d, Name: %s\n", project.Content.Id, project.Content.Name)
	} else {
		fmt.Println("No project found or an error occurred.")
	}

}
