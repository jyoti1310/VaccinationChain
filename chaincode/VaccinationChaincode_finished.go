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
*/

package main

import (
	"errors"
	"fmt"
	"strconv"
	"encoding/json"
	
	"github.com/hyperledger/fabric/core/chaincode/shim"
)

// Chaincode to implement Vaccination schedule of a child
type SimpleChaincode struct {
}

type ChildDetails struct{
		BirthCertID string `json:"BirthCertID"`
        ChildName string `json:"ChildName"`
        FatherName string `json:"FatherName"`
        FatherDOB string `json:"FatherDOB"`
		FatherID string `json:"FatherID"`
		FatherMobile int `json:"FatherMobile"`
		BirthDate string `json:"BirthDate"`
		BirthPlace string `json:"BirthPlace"`
		BirthTime int `json:"BirthTime"`
		Address string `json:"Address"`
}

type VaccinationInfo struct{
		VaccineID string `json:"VaccineID"`
        VaccineName string `json:"VaccineName"`
        QuantiyPrescribed int `json:"QuantiyPrescribed"` //in ml eg 10 ml 
        DosesPrescribed int `json:"DosesPrescribed"`
        DoseInterval int `json:"DoseInterval"`  //in days
        Active int `json:"DoseInterval"`
}

type VaccinationPlan struct{
		VaccineID string `json:"VaccineID"`
        VaccineName string `json:"VaccineName"`
        ChildName string `json:"ChildName"`
        ChildID string `json:"ChildID"`
        ChildDOB string `json:"ChildDOB"`
		PrescribedLocation string `json:"PrescribedLocation"`
		ActualLocation string `json:"ActualLocation"`
		LastDoseDate string `json:"LastDoseDate"`
		NextDoseDate string `json:"NextDoseDate"`
		TotalDosePrescribed int `json:"TotalDosePrescribed"`
		PendingDose int `json:"PendingDose"`
		LastReminderDate string `json:"LastReminderDate"`
		NextReminderDate string `json:"NextReminderDate"`
		
}

//var employeeLogBog map[string]SKATEmployee
type ChildDetailsRepository struct {
	Child_Details []ChildDetails `json:"Child_Details"`
}
type VaccinationInfoRepository struct {
	Vaccination_Info []VaccinationInfo `json:"Vaccination_Info"`
}

type VaccinationPlanRepository struct {
	Vaccination_Plan []VaccinationPlan `json:"Vaccination_Plan"`
}

func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}

// Init resets all the things
func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1")
	}

	err := stub.PutState("hello_Block", []byte(args[0]))
	if err != nil {
		return nil, err
	}

	return nil, nil
}

// Invoke isur entry point to invoke a chaincode function
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("invoke is running " + function)

	// Handle different functions
	if function == "init" {
		return t.Init(stub, "init", args)
	} else if function == "write" {
		return t.write(stub, args)
	} else if  function == "addChildInfo" {
		return t.addNewChildInfo(stub, args)
	} else if  function == "updateVaccinationInfo" {
		return t.addUpdateVaccination(stub, args)
	}
	fmt.Println("invoke did not find func: " + function)

	return nil, errors.New("Received unknown function invocation: " + function)
}

// Query is our entry point for queries
func (t *SimpleChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("query is running " + function)

	// Handle different functions
	if function == "read" { //read a variable
		return t.read(stub, args)
	} 
	fmt.Println("query did not find func: " + function)

	return nil, errors.New("Received unknown function query: " + function)
}

// write - invoke function to write key/value pair
func (t *SimpleChaincode) write(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var key, value string
	var err error
	fmt.Println("running write()")

	if len(args) != 2 {
		return nil, errors.New("Incorrect number of arguments. Expecting 2. name of the key and value to set")
	}

	key = args[0] //rename for funsies
	value = args[1]
	err = stub.PutState(key, []byte(value)) //write the variable into the chaincode state
	if err != nil {
		return nil, err
	}
	return nil, nil
}

// read - query function to read key/value pair
func (t *SimpleChaincode) read(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
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

// ============================================================================================================================
// Init Employee - create a new Employee, store into chaincode state
// ============================================================================================================================
func (t *SimpleChaincode) addNewChildInfo(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var err error
	var key string
	//,jsonResp string
	
	//   0       1       2     3
	// "asdf", "blue", "35", "bob"
	if len(args) != 10 {
		return nil, errors.New("Incorrect number of arguments. Expecting 10 parameters")
	}

	//input sanitation
	fmt.Println("- start adding new Child Info-new")
	
	NewChildInfo := ChildDetails{}
	NewChildInfo.BirthCertID = args[0]
	NewChildInfo.ChildName = args[1]
	NewChildInfo.FatherName = args[2]
    NewChildInfo.FatherDOB = args[3]
	NewChildInfo.FatherID = args[4]
	NewChildInfo.FatherMobile, err = strconv.Atoi(args[5])
	/*if err != nil {
		return nil, errors.New("mobile number must be a numeric string at argument 6")
	}*/
	NewChildInfo.BirthDate = args[6]
	NewChildInfo.BirthPlace = args[7]
	NewChildInfo.BirthTime,err = strconv.Atoi(args[8])
	NewChildInfo.Address = args[9]

	fmt.Println("adding Child Info @ " + NewChildInfo.BirthCertID + ", " + NewChildInfo.FatherID);
	fmt.Println("- end add Child Info 1")
	jsonAsBytes, _ := json.Marshal(NewChildInfo)

	if err != nil {
		return jsonAsBytes, err
	}
	
	key = NewChildInfo.BirthCertID + "_2"
	t.appendtoChildDetailsRepository(stub,key,NewChildInfo)
	fmt.Println("- end add Child Info 2")
	return jsonAsBytes, nil
}

//==================================================================================================================================
//Append to EmployeeRepository
//===================================================================================================================================

func (t *SimpleChaincode) appendtoChildDetailsRepository(stub shim.ChaincodeStubInterface,  key string,newChildInfo ChildDetails) (bool, error){

var jsonResp string 
repositoryJsonAsBytes, err := stub.GetState(key)
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + "Child Info Repository" + "\"}"
		return false, errors.New(jsonResp)
	}
	var childinfoRepository ChildDetailsRepository
	json.Unmarshal(repositoryJsonAsBytes, &childinfoRepository)	

	childinfoRepository.Child_Details = append(childinfoRepository.Child_Details,newChildInfo)
	//update Child Repository
	updatedRepositoryJsonAsBytes, _  := json.Marshal(childinfoRepository)
	err = stub.PutState(key, updatedRepositoryJsonAsBytes)	//store child with id as key
	if err != nil {
		return false, err
	}		
	return true, nil
}

// ============================================================================================================================
// Init vaccination info - add/update vaccination information reference data
// ============================================================================================================================
func (t *SimpleChaincode) addUpdateVaccination(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var err error
	var key,jsonResp string
	
	//   0       1       2     3
	// "asdf", "blue", "35", "bob"
	if len(args) != 6 {
		return nil, errors.New("Incorrect number of arguments. Expecting 6")
	}
	fmt.Println("- adding new vaccination")
	fmt.Println("VaccineID-"+args[0])
	
	NewVaccination := VaccinationInfo{}
	NewVaccination.VaccineID = args[0]
	NewVaccination.VaccineName=args[1]
	NewVaccination.QuantiyPrescribed, err =strconv.Atoi(args[2])
	
	if err != nil {
		return nil, errors.New("QuantiyPrescribed on argument 3 must be a numeric string")
	}
	NewVaccination.DosesPrescribed, err =strconv.Atoi(args[3])
	if err != nil {
		return nil, errors.New("DosesPrescribed on argument 4 must be a numeric string")
	}
	
	NewVaccination.DoseInterval, err=strconv.Atoi(args[4])
	
	if err != nil {
		return nil, errors.New("DoseInterval on argument 5 must be a numeric string")
	}
	
	NewVaccination.Active, err=strconv.Atoi(args[5])
	
	if err != nil {
		return nil, errors.New("Active on argument 6 must be a numeric string")
	}
	
	//NewVaccination,err = t.setApprovalStatus(stub,NewVaccination);
	/*if err != nil {
		jsonResp = "{\"Error\":\"Failed to set Approval " + "\"}"
		return nil, errors.New(jsonResp)
	}*/
	jsonAsBytes, _ := json.Marshal(NewVaccination)
	//Adding new Vaccine to Repository
	key = NewVaccination.VaccineID + "_2"
	fmt.Println("Calling appendtoVaccinationInfo");
	_, err = t.appendtoVaccinationInfo(stub,key,NewVaccination)
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to append to  vaccinationInfoRepository" + "\"}"
		return nil, errors.New(jsonResp)
	}
	fmt.Println("- end add Vaccination 2")
	return jsonAsBytes, nil
}
//==================================================================================================================================
//Append to Vaccination Info repository
//===================================================================================================================================
func (t *SimpleChaincode) appendtoVaccinationInfo(stub shim.ChaincodeStubInterface,  key string,newVaccine VaccinationInfo) (bool, error){
var jsonResp string 
repositoryJsonAsBytes, err := stub.GetState(key)
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + "VaccinationInfoRepository" + "\"}"
		return false, errors.New(jsonResp)
	}
	var vaccinationInfoRepository VaccinationInfoRepository
	json.Unmarshal(repositoryJsonAsBytes, &vaccinationInfoRepository)	
	vaccinationInfoRepository.Vaccination_Info = append(vaccinationInfoRepository.Vaccination_Info,newVaccine)
	//update Employee Repository
	updatedRepositoryJsonAsBytes, _  := json.Marshal(vaccinationInfoRepository)
	err = stub.PutState(key, updatedRepositoryJsonAsBytes)	//store vaccine with id as key
	
	if err != nil {	
		jsonResp = "{\"Error\":\"Failed to init vaccination info " + "\"}"
		return false, errors.New(jsonResp)
	}		
	return true, nil
}

	
