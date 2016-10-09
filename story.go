/*
Copyright IBM Corp 2016 All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

		 http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
// MAIN METHODS:
DEPLOY
QUERY
INVOKE
*/

package main

import (
	"errors"
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}


// You need to create a short main function that will execute when each peer deploys their instance of the chaincode. It just starts the chaincode and registers it with the peer. You don’t need to add any code for this function. Both chaincode_start.go and chaincode_finished.go have a main function that lives at the top of the file. The function looks like this:

func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}



func (t *SimpleChaincode) Init(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
    var Aval int
    var err error

    if len(args) != 1 {
        return nil, errors.New("Incorrect number of arguments. Expecting 1")
    }

    // Initialize the chaincode
    Aval, err = strconv.Atoi(args[0])
    if err != nil {
        return nil, errors.New("Expecting integer value for asset holding")
    }

    // Write the state to the ledger
    err = stub.PutState("scoopIndexStr", []byte(strconv.Itoa(Aval)))              //making a test var "abc", I find it handy to read/write to it right away to test the network
    if err != nil {
        return nil, err
    }

    var empty []string
    jsonAsBytes, _ := json.Marshal(empty)                               //marshal an emtpy array of strings to clear the index
    err = stub.PutState(scoopIndexStr, jsonAsBytes)
    if err != nil {
        return nil, err
    }

    // var trades AllTrades
    // jsonAsBytes, _ = json.Marshal(trades)                               //clear the open trade struct
    // err = stub.PutState(openTradesStr, jsonAsBytes)
    // if err != nil {
    //     return nil, err
    // }
    // return nil, nil
}


// Invoke is our entry point to invoke a chaincode function



func (t *SimpleChaincode) Invoke(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
    fmt.Println("invoke is running " + function)

    // Handle different functions
    if function == "init" {
        return t.Init(stub, "init", args)
    } else if function == "write" {
        return t.write(stub, args)
    }
    fmt.Println("invoke did not find func: " + function)

    return nil, errors.New("Received unknown function invocation")
}

//Now that it’s looking for write let’s make that function somewhere in your chaincode_start.go file.

func (t *SimpleChaincode) write(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
    var key, value string
    var err error
    fmt.Println("running write()")

    if len(args) != 2 {
        return nil, errors.New("Incorrect number of arguments. Expecting 2. name of the key and value to set")
    }

    key = args[0]                            //rename for fun
    value = args[1]
    err = stub.PutState(key, []byte(value))  //write the variable into the chaincode state
    if err != nil {
        return nil, err
    }
    return nil, nil
}

// This write function should look similar to the Init change you just did. One major difference is that you can now set the key and value for PutState. This function allows you to store any key/value pair you want into the blockchain ledger.

// Query()

// As the name implies, Query is called whenever you query your chaincode state. Queries do not result in blocks being added to the chain. You will use Query to read the value of your chaincode state's key/value pairs.

// Original query:
// func (t *SimpleChaincode) Query(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
// 	fmt.Println("query is running " + function)

// 	// Handle different functions
// 	if function == "dummy_query" {											//read a variable
// 		fmt.Println("hi there " + function)						//error
// 		return nil, nil;
// 	}
// 	fmt.Println("query did not find func: " + function)						//error

// 	return nil, errors.New("Received unknown function query")
// }

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
    var key, jsonResp string
    var err error

    if len(args) != 1 {
        return nil, errors.New("Incorrect number of arguments. Expecting name of the key to query")
    }

    key = args[0]
    valAsbytes, err := stub.GetState(key)
    if err != nil {
        jsonResp = "{\"Error\":\"Failed to get state for " + key + "\"}"
        return nil, errors.New(jsonResp)
    }

    return valAsbytes, nil
}

// This read function is using the complement to PutState called GetState. This shim function just takes one string argument. The argument is the name of the key to retrieve. Next, this function returns the value as an array of bytes back to Query, who in turn sends it back to the REST handler.

