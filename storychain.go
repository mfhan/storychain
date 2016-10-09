/*
Copyright IBM Corp 2016 All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

		 http://www.apache.org/licenses/LICENSE-2.0

// MAIN METHODS:
DEPLOY
QUERY
INVOKE
*/

package main

import (
	"errors"
	"fmt"
    "encoding/json"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

// SimpleChaincode example simple Chaincode implementatijon


// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}


type story struct {
    name string
    state string
}



var storiesString = "stories"
// You need to create a short main function that will execute when each peer deploys their instance of the chaincode. It just starts the chaincode and registers it with the peer. You don’t need to add any code for this function. Both chaincode_start.go and chaincode_finished.go have a main function that lives at the top of the file. The function looks like this:

func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}

// Init resets all the things

// func (t *SimpleChaincode) Init(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
// 	if len(args) != 1 {
// 		return nil, errors.New("Incorrect number of arguments. Expecting 1")
// 	}
// 	return nil, nil
// }
//change the Init function so that it stores the first element in the args argument to the key "hello_world".
// ('t' is the OBJECT that implements SimpleChaincode )

func (t *SimpleChaincode) Init(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
    if len(args) != 1 {
        return nil, errors.New("Incorrect number of arguments. Expecting 1")
    }

    var stories []string
    storiesJson, _ := json.Marshal(stories)

    // Write a blank list to ledger
    err := stub.PutState(storiesString, storiesJson)
    if err != nil {
        return nil, err
    }

    return nil, nil
}


// Invoke is our entry point to invoke a chaincode function

//Initial version:
// func (t *SimpleChaincode) Invoke(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
// 	fmt.Println("invoke is running " + function)

// 	// Handle different functions
// 	if function == "init" {													//initialize the chaincode state, used as reset
// 		return t.Init(stub, "init", args)
// 	}
// 	fmt.Println("invoke did not find func: " + function)					//error

// 	return nil, errors.New("Received unknown function invocation")
// }

//In your chaincode_start.go file, change the Invoke function so that it calls a generic write function.

func (t *SimpleChaincode) Invoke(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
    fmt.Println("invoke is running " + function)

    // Handle different functions
    if function == "init" {
        return t.Init(stub, "init", args)
    } else if function == "addstory" {
        return t.addstory(stub, args)
    } else if function == "firstedit" {
        return t.firstedit(stub, args)
    } else if function == "approve" {
        return t.approve(stub, args)
    }

    fmt.Println("invoke did not find func: " + function)

    return nil, errors.New("Received unknown function invocation")
}


func (t *SimpleChaincode) addstory(stub *shim.ChaincodeStub, args []string) ([]byte, error) {

    if len(args) != 1 {
        return nil, errors.New("Incorrect number of arguments. Expecting 1")
    }

    storyName := args[0]

    //get all the stories
    storiesAsBytes, err := stub.GetState(storiesString)
    if err != nil {
        return nil, errors.New("Failed to get all stories")
    }
    var stories []string
    json.Unmarshal(storiesAsBytes, &stories)

    for i:=0; i < len(stories); i++ {
        if stories[i] == storyName {
            return nil, errors.New("Cannot add this story, it already exists")
        }
    }

    var newStory story
    newStory.name = storyName
    newStory.state = "written"

    thisStoryJson, _ := json.Marshal(newStory)
    err = stub.PutState(storyName, thisStoryJson)
    if err != nil {
        return nil, err
    }

    stories = append(stories, storyName)
    storiesJson, _ := json.Marshal(stories)

    // Write a blank list to ledger
    err = stub.PutState(storiesString, storiesJson)
    if err != nil {
        return nil, err
    }

    return nil, nil
}

func (t *SimpleChaincode) changeState(stub *shim.ChaincodeStub, args []string) ([]byte, error) {

    storyName := args[0]
    fromState := args[1]
    toState := args[2]

    var theStory story

    storyAsBytes, err := stub.GetState(storyName)
    if err != nil {
        return nil, errors.New("Story not found")
    }
    json.Unmarshal(storyAsBytes, &theStory)

    if theStory.state != fromState {
        return nil, errors.New("Story not in the from State: " + fromState)
    }

    theStory.state = toState
    theStoryJson, _ := json.Marshal(theStory)
    err = stub.PutState(storyName, theStoryJson)
    if err != nil {
        return nil, err
    }

    return nil, nil
}

func (t *SimpleChaincode) firstedit(stub *shim.ChaincodeStub, args []string) ([]byte, error) {

    storyName := args[0]

    var arguments []string

    arguments = append(arguments, storyName)
    arguments = append(arguments, "written")
    arguments = append(arguments, "firstedited")

    t.changeState(stub, arguments)

    return nil, nil
}


func (t *SimpleChaincode) approve(stub *shim.ChaincodeStub, args []string) ([]byte, error) {

    storyName := args[0]
    var arguments []string

    arguments = append(arguments, storyName)
    arguments = append(arguments, "firstedited")
    arguments = append(arguments, "approved")

    t.changeState(stub, arguments)

    return nil, nil
}


// We're changing the Query function so that it calls a generic read function.

func (t *SimpleChaincode) Query(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
    fmt.Println("query is running " + function)

    // Handle different functions
    if function == "read" {                            //read a variable
        return t.read(stub, args)
    }
    fmt.Println("query did not find func: " + function)

    return nil, errors.New("Received unknown function query")
}

// Now that it’s looking for read, make that function somewhere in your chaincode_start.go file.

func (t *SimpleChaincode) read(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
    var storyName, jsonResp string
    var err error

    if len(args) != 1 {
        return nil, errors.New("Incorrect number of arguments. Expecting name of the key to query")
    }

    storyName = args[0]
    valAsbytes, err := stub.GetState(storyName)
    if err != nil {
        jsonResp = "{\"Error\":\"Failed to get state for " + storyName + "\"}"
        return nil, errors.New(jsonResp)
    }

    return valAsbytes, nil
}

// This read function is using the complement to PutState called GetState. This shim function just takes one string argument. The argument is the name of the key to retrieve. Next, this function returns the value as an array of bytes back to Query, who in turn sends it back to the REST handler.

