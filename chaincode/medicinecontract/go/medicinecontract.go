
package main

import (
	"encoding/json"
    "fmt"
    "time"
    "strconv"
	"strings"
	"bytes"
    "github.com/hyperledger/fabric/core/chaincode/shim/ext/cid"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
	"crypto/x509"
)

type MedicineContract struct {

}

//DocType used to distinguish records in state database
type Company struct {
	CompanyId              string           `json:"companyid"`
	DocType                string           `json:"doctype"`  
	Name                   string           `json:"name"`
	Location               string           `json:"location"`
	OrganizationRole       string           `json:"organizationrole"`
	HierarchyKey           int              `json:"hierarchykey"`
	CreatedAt              string           `json:"createdat"`

}



type Medicine struct {
	ProductId              string           `json:"productid"`
	DocType                string           `json:"doctype"`
	Name                   string           `json:"name"`
	SerialNo               string           `json:"serialno"`
	ManufacturingDate      string           `json:"manufacturingdate"`
	ExpiryDate             string           `json:"expirydate"`
	Owner                  string           `json:"owner"`
	ManufacturerId         string           `json:"manufacturerid"`
	ManufacturerName       string           `json:"manufacturername"`
	ShipmentList           []string         `json:"shipmentlist"`
	CreatedAt              string           `json:"createdat"`
}

type ProductOrder struct {
	ProductOrderId         string           `json:"productorderid"`
	DocType                string           `json:"doctype"`
	MedicineName           string           `json:"medicinename"`
	Quantity               string           `json:"quantity"`
	BuyerKey               string           `json:"buyerkey"`
	SellerKey              string           `json:"sellerkey"`
	CreatedAt              string           `json:"createdat"`
}

type Shipment struct {
	BuyerId                string           `json:"buyerid"`
	DocType                string           `json:"doctype"`
	MedicineName           string           `json:"medicinename"`
	TransporterId          string           `json:"TransporterId"`
	Creator                string           `json:"creator"`
	AssetList              []string         `json:"assetlist"`
	CreatedAt              string           `json:"createdat"`
	UpdatedAt              string           `json:"updatedat"`
	Status                 string           `json:"status"`
}

type CounterNO struct {
	Counter                int              `json:"counter"`
}

// enum types

type DocumentType int

const (
	CompanyDocType DocumentType = 1 + iota
	MedicineDocType 
	ProductOrderDocType 
	ShipmentDocType 
)

type ShipmentStatus int

const (
	InTransitStatus ShipmentStatus = 1 + iota
	DeliveredStatus 
)

var statuses= [...]string {
    "intransit",
    "delivered",
}

var doctypes= [...]string {
	"company",
	"medicine",
	"productorder",
    "shipment",
}

func (s ShipmentStatus) String() string {return statuses[s-1]}

func (d DocumentType) String() string { return doctypes[d-1]}

var logger = shim.NewLogger("medicinecontract")

func (s *MedicineContract) Init(APIstub shim.ChaincodeStubInterface) sc.Response {
	logger.Info("Inside Init")
	// Initializing Company Counter
	CompanyCounterBytes, _ := APIstub.GetState("CompanyCounterNO")
	if CompanyCounterBytes == nil {
		var CompanyCounter = CounterNO{Counter: 0}
		CompanyCounterBytes, _ := json.Marshal(CompanyCounter)
		err := APIstub.PutState("CompanyCounterNO", CompanyCounterBytes)
		if err != nil {
			return shim.Error(fmt.Sprintf("Failed to Initiate Company Counter"))
		}
	} 
	logger.Info("Init executed")
	
    return shim.Success(nil)
}

func (s *MedicineContract) Invoke(APIstub shim.ChaincodeStubInterface) sc.Response {

	// Retrieve the requested Smart Contract function and arguments
	function, args := APIstub.GetFunctionAndParameters()
	fmt.Println("function name, ", function)
	// Route to the appropriate handler function to interact with the ledger appropriately
	if function == "registerCompany" {
		return s.registerCompany(APIstub, args)
	} else if function == "addMedicine" {
		return s.addMedicine(APIstub, args)
	} else if function == "getCompany" {
		return s.getCompany(APIstub, args)
	}else if function == "createPurchaseOrder" {
		return s.createPurchaseOrder(APIstub, args)
	} else if function == "createShipment" {
		return s.createShipment(APIstub, args)
	} else if function == "updateShipment" {
		return s.updateShipment(APIstub, args)
	}else if function == "retailMedicine" {
		return s.retailMedicine(APIstub, args)
	} else if function == "viewHistory" {
		return s.viewHistory(APIstub, args)
	} else if function == "viewMedicineCurrentState" {
		return s.viewMedicineCurrentState(APIstub, args)
	}else if function == "getMedicineByManufacturer" {
		return s.getMedicineByManufacturer(APIstub, args)
	}

	return shim.Error("Invalid Smart Contract function name.")
}


//args
//args[0]=companyName
//args[1]=location
//args[2]=organizationrole


func (s *MedicineContract) registerCompany(APIstub shim.ChaincodeStubInterface, args []string) sc.Response{
	logger.Info("Init registerCompany")

	if len(args) != 3 {
		logger.Error("invalid number of arguments. Expecting 3")
		return shim.Error("invalid number of arguments. Expecting 3")
	}
	
	var hierarchyKey int

	if strings.ToLower(args[2]) == "manufacturer" {
		hierarchyKey = 1
	} else if strings.ToLower(args[2]) == "distributor" {
		hierarchyKey = 2
	} else if strings.ToLower(args[2]) == "retailer" {
		hierarchyKey = 3
	}

	
	companyCounter := getCounter(APIstub,"CompanyCounterNO")
	logger.Debug("Retrieved company counter")
	companyCounter++
	

	
	txTimeAsPtr, errTx := s.GetTxTimestampChannel(APIstub)
	if errTx != nil {
		logger.Error("Error in gettng timestamp fo channel: "+errTx.Error())
		return shim.Error("Error in gettng timestamp fo channel: "+errTx.Error())
	}
	var company = Company{CompanyId:strconv.Itoa(companyCounter),DocType:CompanyDocType.String(), Name:args[0], Location:args[1], OrganizationRole:args[2], HierarchyKey:hierarchyKey, CreatedAt:txTimeAsPtr}
   
	companyAsBytes, errMarshal := json.Marshal(company)
	
	if errMarshal != nil {
		logger.Error("Marshal Error for Company")
		return shim.Error(fmt.Sprintf("Marshal Error for Company: %s", errMarshal))
	}

	errPut := APIstub.PutState(company.CompanyId, companyAsBytes)
	if errPut != nil {
		logger.Error("Failed to add company to ledger")
		return shim.Error(fmt.Sprintf("Failed to add Company: %s", company.CompanyId))
	}
	
    logger.Info("Registered Company")
	//TO Increment the Company Counter
	incrementCounter(APIstub,"CompanyCounterNO")
	
    logger.Debug("Incremented company counter")
	return shim.Success(companyAsBytes)


}

//args
//args[0]=name
//args[1]=serialno
//args[2]=manufacturedate
//args[3]=expirydate
//args[4]=owner
//args[5]=manufacturerid
//args[6]=manufacturername


func (s *MedicineContract) addMedicine(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	logger.Info("Inside add Medicine")

	if len(args) != 7 {
		logger.Error("Invalid number of arguments. Expecting 7")
		return shim.Error("Invalid number of arguments. Expecting 7")
	}
	isValid:=isInitiatorValid(APIstub,[]string{"manufacturer.pharma-network.com"})
	if !isValid {
		logger.Error("initiator is not valid. Should be manufacturer")
		return shim.Error("Initiator is not valid. Should be manufacturer")
	}

	shipmentList := []string{}

	//date := time.Now().Format("02-01-2006 15:04:05")
	// To get the transaction Timestamp from the channel header
	txTimeAsPtr, errTx := s.GetTxTimestampChannel(APIstub)
	if errTx != nil {
		logger.Error("Error in gettng timestamp fo channel: "+errTx.Error())
		return shim.Error("Error in gettng timestamp fo channel: "+errTx.Error())
	}
	objectType := "serialno~name"
	medicineKey,_ := APIstub.CreateCompositeKey(objectType, []string{args[1],args[0]})
	
	logger.Debug("Medicine composite key created")
	
	var medicine = Medicine{ProductId:medicineKey, DocType:MedicineDocType.String(), Name:args[0], SerialNo:args[1], ManufacturingDate:args[2], ExpiryDate:args[3], Owner:args[4], ManufacturerId:args[5], ManufacturerName:args[6], ShipmentList:shipmentList, CreatedAt:txTimeAsPtr}

	medicineAsBytes, errMarshal := json.Marshal(medicine)
	
	
	if errMarshal != nil {
		logger.Error("Marshal error for Medicine")
		return shim.Error(fmt.Sprintf("Marshal Error for Medicine: %s", errMarshal))
	}

	errPut := APIstub.PutState(medicineKey, medicineAsBytes)
	if errPut != nil {
        logger.Error("Error adding medicine to ledger")
		return shim.Error(fmt.Sprintf("Error adding Medicine to ledger: %s", errMarshal))
	}
	logger.Info("added medicine to ledger")
	

	
   return shim.Success(nil)

}

func (s *MedicineContract) getCompany(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	logger.Info("Inside getCompany")
	if len(args) != 1 {
		logger.Error("Invalid number of arguments. Expecting 1")
		return shim.Error("Invalid number of arguments. Expecting 1")
	}
	companyAsBytes, err := APIstub.GetState(args[0])
	if err!=nil {
		logger.Error("error in getting company")
		return shim.Error("error in getting company")
	}
	logger.Info("retrieved company")

	return shim.Success(companyAsBytes)
}

//args
//args[0]=buyerKey
//args[1]=sellerkey
//args[2]=medicinename
//args[3]=quantity
func (s *MedicineContract) createPurchaseOrder(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	logger.Info("Inside createPurchaseOrder")  
	if len(args) != 4{
		logger.Error("Invalid number of arguments. Expecting 4")
		shim.Error("Invalid number of arguments. Expecting 4")
	}
	isValid:= isInitiatorValid(APIstub, []string{"retailer.pharma-network.com", "distributor.pharma-network.com"})
	   
	  if !isValid {
		  logger.Error("Initiator is not valid. Expecting retailer or distributor")
		return shim.Error("Initiator is not valid. Expecting retailer or distributor")
	}
	  
	   objectType:="buyerkey~medicinename"
	   purchaseOrderKey,_ := APIstub.CreateCompositeKey(objectType, []string{args[0],args[2]})
       logger.Debug("Purchase order key created")
	   
	   buyerAsBytes,_ := APIstub.GetState(args[0])
	   buyerCompany := Company{}
	   json.Unmarshal(buyerAsBytes, &buyerCompany)

	   sellerAsBytes,_ := APIstub.GetState(args[1])
	   sellerCompany := Company{}
	   json.Unmarshal(sellerAsBytes, &sellerCompany)

	   if buyerCompany.HierarchyKey-sellerCompany.HierarchyKey!=1{
		   logger.Error("Cannot create Purchase order as hierarchy of transfer is not correct")
		   return shim.Error("Cannot create Purchase order as hierarchy of transfer is not correct")
	   }

	   purchaseOrder:= ProductOrder{ProductOrderId:purchaseOrderKey, DocType:ProductOrderDocType.String(), MedicineName:args[2], Quantity: args[3], BuyerKey:args[0], SellerKey: args[1]}

	   purchaseOrderAsBytes, _ := json.Marshal(purchaseOrder)
	   errPut := APIstub.PutState(purchaseOrderKey, purchaseOrderAsBytes)
	   
	   if errPut != nil {
		   logger.Error("error in adding purchasing order to ledger")
		   return shim.Error("error in creating purchase order")
	   }
	   logger.Info("Purchase order created")


	   return shim.Success(purchaseOrderAsBytes)


}

//args
//args[0]=buyerkey
//args[1]=medicinename
//args[2]=listofassets
//args[3]=transporterkey

func (s *MedicineContract) createShipment(APIstub shim.ChaincodeStubInterface, args []string) sc.Response{
	logger.Info("inside createShipment")

	if len(args) != 4{
		logger.Error("Invalid number of arguements. Expecting 4")
		return shim.Error("Invalid number of arguments. Expecting 4")
	}
	isValid:=isInitiatorValid(APIstub, []string{"manufacturer.pharma-network.com", "distributor.pharma-network.com"})

	if !isValid {
		fmt.Println("Initiator is not valid. Expecting manufacturer or distributor")
		return shim.Error("Initiator is not valid. Expecting manufacturer or distributor")
	}
	objectType:= "buyerkey~medicinename"
	purchaseOrderKey,_ := APIstub.CreateCompositeKey(objectType,[]string{args[0],args[1]})
	logger.Debug("Created purchase composite key")
	purchaseOrderAsBytes,_:=APIstub.GetState(purchaseOrderKey)
	purchaseOrder:=ProductOrder{}
	json.Unmarshal(purchaseOrderAsBytes, &purchaseOrder)

	listofassets:=strings.Split(args[2],",")

	if strconv.Itoa(len(listofassets)) != purchaseOrder.Quantity {
		logger.Error("Quantity of medicine mentioned in purchase order does not match the number of assets")
		return shim.Error("Quantity of medicine mentioned in purchase order does not match the number of assets")

	}

	for i:=0 ;i<len(listofassets); i++ {
        medicineObjectType:= "serialno~name"
		medicineCompositeKey,_ := APIstub.CreateCompositeKey(medicineObjectType, []string{listofassets[i],args[1]})
		logger.Debug("Created medicine composite key")
		medicineAsBytes, _ := APIstub.GetState(medicineCompositeKey)
		if medicineAsBytes == nil {
			return shim.Error("error in getting asset details from ledger")
		}

	}

	sellerAsBytes,_:=APIstub.GetState(purchaseOrder.SellerKey)
	sellerCompany:=Company{}
	json.Unmarshal(sellerAsBytes, &sellerCompany)

	txTimeAsPtr, errTx := s.GetTxTimestampChannel(APIstub)

	if errTx != nil {
		logger.Error("Error in getting timestamp fro channel")
		return shim.Error("Error in gettng timestamp fo channel: "+errTx.Error())
	}
	shipment := Shipment {
		BuyerId:args[0],
		DocType:ShipmentDocType.String(),
		MedicineName:args[1],
		TransporterId:args[3],
		AssetList:listofassets,
		Status:InTransitStatus.String(),
		Creator:purchaseOrder.SellerKey,
		CreatedAt:txTimeAsPtr,
		UpdatedAt:txTimeAsPtr }

	shipmentObjectType:="buyerkey~medicinename~doctype"
	shipmentCompositeKey,_:=APIstub.CreateCompositeKey(shipmentObjectType,[]string{args[0],args[1],ShipmentDocType.String()})
    logger.Debug("Created shipment composite key")
	shipmentAsBytes,_:=json.Marshal(shipment)
	errPut:=APIstub.PutState(shipmentCompositeKey, shipmentAsBytes)
	if errPut!=nil{
        logger.Error("Error in adding shipment to ledger")
		return shim.Error("Error in shipment put state.Error: "+errPut.Error())
	}
	logger.Info("Created shipment ")
	return shim.Success(shipmentAsBytes)
}

//args
//args[0]=buyerkey
//args[1]=medicinename
//args[2]=transporterkey

func (s *MedicineContract) updateShipment(APIstub shim.ChaincodeStubInterface, args []string) sc.Response{
	logger.Info("Inside updateShipment")

	if len(args) != 3{
		logger.Error("invalid number of arguments. Expecting 3")
		return shim.Error("invalid number of arguments. Expecting 3")
	}
	isValid:=isInitiatorValid(APIstub,[]string{"transporter.pharma-network.com"})

	if !isValid {
		logger.Error("Initiator is not valid. Expecting transporter")
		return shim.Error("Initiator is not valid. Expecting transporter")
	}
	objectType:="buyerkey~medicinename~doctype"
	shipmentKey,_:=APIstub.CreateCompositeKey(objectType,[]string{args[0],args[1],ShipmentDocType.String()})
    logger.Debug("Create shipment composite key")
	shipmentAsBytes,_:= APIstub.GetState(shipmentKey)
	shipment:=Shipment{}
	json.Unmarshal(shipmentAsBytes,&shipment)

	if args[2]!=shipment.TransporterId{
		logger.Error("transporter key passed doesnt match the transportId of shipment object")
		return shim.Error("transporter key passed doesnt match the transportId of shipment object")

	}
	txTimeAsPtr, errTx := s.GetTxTimestampChannel(APIstub)
	if errTx != nil {
		logger.Error("Error in getting timestamp of channel")
		return shim.Error("Error in getting timestamp of channel: "+errTx.Error())
	}
	shipment.Status=DeliveredStatus.String()
	shipment.UpdatedAt=txTimeAsPtr
	newShipmentAsBytes,_:=json.Marshal(shipment)
	errPut:= APIstub.PutState(shipmentKey,newShipmentAsBytes)
	if errPut!=nil{
		logger.Error("error in updating the shipment in ledger")
		return shim.Error("error in put state of update shipment")
	}

	for i:=0;i<len(shipment.AssetList);i++{
		medicineObjectType:="serialno~name"
		medicineKey,_:=APIstub.CreateCompositeKey(medicineObjectType,[]string{shipment.AssetList[i],args[1]})
        logger.Info("create medicine composite key")
		medicineAsBytes,_:=APIstub.GetState(medicineKey)
		if medicineAsBytes==nil{
			logger.Error("Error in gettting asset details from ledger")
			return shim.Error("Error in getting asset details from ledger")
		}
		medicine:=Medicine{}
		json.Unmarshal(medicineAsBytes, &medicine)
		medicine.Owner=args[0]
		medicine.ShipmentList=append(medicine.ShipmentList,shipmentKey)  
		newMedicineAsBytes,_:=json.Marshal(medicine)
		errMedicinePut := APIstub.PutState(medicineKey,newMedicineAsBytes)
		if errMedicinePut != nil {
			logger.Error("Error in updating medicine in ledger")
			return shim.Error("Error in put state of medicine")
		}
        
		
	}
    logger.Info("updated shipment")
	return shim.Success(newShipmentAsBytes)

}

//args
//args[0]=medicinename
//args[1]=serialno
//args[2]=retailerKey
//args[3]=customerIdentityNumber
func (s *MedicineContract) retailMedicine(APIstub shim.ChaincodeStubInterface, args []string) sc.Response{
	logger.Info("inside retailMedicine")
	if len(args) != 4 {
		logger.Error("invalid number of arguments. Expecting 4")
		return shim.Error("invalid number of arguments. Expecting 4")
	}
	isValid:=isInitiatorValid(APIstub, []string{"retailer.pharma-network.com"})
	if !isValid {
		logger.Error("Initiator is not valid. Expecting retailer")
		return shim.Error("Initiator is not valid. Expecting retailer")
	}
	objectType:="serialno~name"
	compositeKey,_ := APIstub.CreateCompositeKey(objectType,[]string{args[1],args[0]})
	logger.Debug("create medicine composite key")
	medicineAsBytes,_ := APIstub.GetState(compositeKey)
	if medicineAsBytes ==nil{
		logger.Error("Error in retrieving medicine using key")
		return shim.Error("Error in retrieivng the medicine using key")
	}
	medicine:=Medicine{}
	json.Unmarshal(medicineAsBytes,&medicine)

	if medicine.Owner != args[2]{
		logger.Error("retialer key does not match with owner from medicine object")
		return shim.Error("Retailer key does not match with owner from medicine object")
	}
	medicine.Owner=args[3]
	newMedicineAsBytes,_:= json.Marshal(medicine)
	errPut := APIstub.PutState(compositeKey,newMedicineAsBytes)
	if errPut!=nil {
		logger.Error("error in updating the medicine to ledger")
		return shim.Error("Error in Put state of updated medicine")
	}
    logger.Info("Medicine updated in retailMedicine")
	return shim.Success(newMedicineAsBytes)

}

//args
//args[0]=serialno
//args[1]=medicinename

func (s *MedicineContract) viewHistory(APIstub shim.ChaincodeStubInterface, args []string) sc.Response{
	logger.Info("inside viewHistory")
	if len(args) != 2 {
		logger.Error("Invalid number of arguments. Expecting 2")
		return shim.Error("Invalid number of arguments. Expecting 2")
	}
	isValid:=isInitiatorValid(APIstub,[]string{"retailer.pharma-network.com","transporter.pharma-network.com","distributor.pharma-network.com","manufacturer.pharma-network.com"})
	if !isValid {
		logger.Error("initiator is not valid. Expecting distributor, trnasporter or retailer")
		return shim.Error("Initiator is not valid. Expecting distributor, trnasporter or retailer")
	}
	objectType:="serialno~name"
	medicineCompositeKey,_:=APIstub.CreateCompositeKey(objectType, []string{args[0],args[1]})
    logger.Debug("medicine compisite key created")
	historyQueryIterator,err := APIstub.GetHistoryForKey(medicineCompositeKey)
	if err!=nil{
		logger.Error("Error in fetching history")
		return shim.Error(" Error in fetcing history" + err.Error())
	}

	

	defer historyQueryIterator.Close()

	var buffer bytes.Buffer
	buffer.WriteString("[")
	bArrayMemberAlreadyWritten := false

	for historyQueryIterator.HasNext() {
		response, err := historyQueryIterator.Next()
		if err != nil {
			logger.Error("Error in reading next history record")
			return shim.Error("Error in reading next history record: "+err.Error())
		}
		if bArrayMemberAlreadyWritten ==true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"txn\":")
		buffer.WriteString("\"")
		buffer.WriteString(response.TxId)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Value\":")

		//if it is deleted then set to null
		if response.IsDelete {
			buffer.WriteString("null")
		} else {
			buffer.WriteString(string(response.Value))
		}

		buffer.WriteString(", \"Timestamp\":")
		buffer.WriteString("\"")
		buffer.WriteString(time.Unix(response.Timestamp.Seconds, int64(response.Timestamp.Nanos)).String())
		buffer.WriteString("\"")

		buffer.WriteString(",\"IsDelete\":")
		buffer.WriteString("\"")
		buffer.WriteString(strconv.FormatBool(response.IsDelete))
		buffer.WriteString("\"")

		buffer.WriteString("}")
		bArrayMemberAlreadyWritten=true
	}
	buffer.WriteString("]")
	logger.Info("viewHistory: " + buffer.String())
	return shim.Success(buffer.Bytes())

}

//args
//args[0]=serialno
//args[1]=medicinename

func (s *MedicineContract) viewMedicineCurrentState(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	logger.Info("inside viewMedicineCurrentState")
	if len(args) != 2 {
		logger.Error("Invalid number of arguments. Expecting 2")
		return shim.Error("Invalid number of arguments. Expecting 2")
	}
	
	objectType:="serialno~name"
	medicineCompositeKey,_:=APIstub.CreateCompositeKey(objectType,[]string{args[0],args[1]})
	logger.Debug("Created medicine composite key")
	medicineAsBytes,_:=APIstub.GetState(medicineCompositeKey)
	if medicineAsBytes==nil{
		logger.Error("Invalid medicine retrieved")
		return shim.Error("Invalid medicine retrieved ")
	}
	logger.Info("Retrieved medicine current state")
	return shim.Success(medicineAsBytes)
}

func (s *MedicineContract) getMedicineByManufacturer(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	logger.Info("inside getMedicineByManufacturer")
	if len(args) != 1 {
		logger.Error("invalid number of arguments. Expecting 1")
		return shim.Error("invalid number of arguments. Expecting 1")
	}
	isValid:=isInitiatorValid(APIstub, []string{"retailer.pharma-network.com","transporter.pharma-network.com","distributor.pharma-network.com"})
	if !isValid {
        logger.Error("Initiator is not valid. Expecting retailer, transporter or distributor")
		return shim.Error("Initiator is not valid. Expecting retailer, trnasporter or distributor")
	}
	// {
	// 	"selector": {
	// 		"doctype": "medicine"
    // 	   "manufacturerName":args[0]
	//    }
	// 	}


	queryString:= fmt.Sprintf("{\"selector\":{\"doctype\":\"medicine\",\"manufacturername\":\""+args[0]+"\"}}")
	var pagesize int32 =20
	bookmark:=""

	

	resultsIterator, responseMetadata, err := APIstub.GetQueryResultWithPagination(queryString, pagesize, bookmark)
	if err != nil {
		logger.Error("Error in pagination")
		return shim.Error("Error in pagination :"+ err.Error())
	}

	defer resultsIterator.Close()

	var buffer bytes.Buffer

	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			logger.Error("Error in reading next record")
			return shim.Error("Error in reading next record: "+err.Error())
		}

		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"Key\":")
		buffer.WriteString("\"")
		buffer.WriteString(queryResponse.Key)
		buffer.WriteString("\"")

		buffer.WriteString(",\"Record\":")
		buffer.WriteString(string(queryResponse.Value))
		buffer.WriteString("}")
		bArrayMemberAlreadyWritten=true
	}
	buffer.WriteString("]")

	var bufferWithPagination bytes.Buffer 

	bufferWithPagination.WriteString("[{\"ResponseMetadata\":{\"RecordsCount\":")
	bufferWithPagination.WriteString("\"")
	bufferWithPagination.WriteString(fmt.Sprintf("%v", responseMetadata.FetchedRecordsCount))
	bufferWithPagination.WriteString("\"")
	bufferWithPagination.WriteString(", \"Bookmark\":")
	bufferWithPagination.WriteString("\"")
	bufferWithPagination.WriteString(responseMetadata.Bookmark)
	bufferWithPagination.WriteString("\"}}]")


	

	logger.Info("get queryresult with pagination result: " + bufferWithPagination.String())

	return shim.Success(buffer.Bytes())

}


func getCounter(APIstub shim.ChaincodeStubInterface, AssetType string) int {
	counterAsBytes, _ := APIstub.GetState(AssetType)
	counterAsset := CounterNO{}

	json.Unmarshal(counterAsBytes, &counterAsset)
	fmt.Sprintf("Counter Current Value %d of Asset Type %s",counterAsset.Counter,AssetType)

	return counterAsset.Counter
}

func incrementCounter(APIstub shim.ChaincodeStubInterface,  AssetType string) int {
	counterAsBytes, _ := APIstub.GetState(AssetType)
	counterAsset := CounterNO{}

	json.Unmarshal(counterAsBytes, &counterAsset)
	counterAsset.Counter++
	counterAsBytes, _ = json.Marshal(counterAsset)

	err := APIstub.PutState(AssetType, counterAsBytes)
	if err != nil {

		fmt.Sprintf("Failed to Increment Counter")

	}
	return counterAsset.Counter
}

func isInitiatorValid(APIstub shim.ChaincodeStubInterface, initiators []string) bool{
	var cert *x509.Certificate
	
	cert,_ = cid.GetX509Certificate(APIstub)


	for i:=0 ; i<len(initiators) ;i++ {

		for _, a := range cert.Issuer.Organization{
			if a == initiators[i] {
				return true
			}
		}
		
		
	}
	return false
	
}

func (s *MedicineContract) GetTxTimestampChannel(APIstub shim.ChaincodeStubInterface) (string,error){
	txTimeAsPtr, err := APIstub.GetTxTimestamp()
	if err!=nil {
		fmt.Println("Returning error in TimeStamp")
		return "Error", err
	}
	fmt.Printf("\t returned value from APIstub: %v\n",txTimeAsPtr)
	timeStr := time.Unix(txTimeAsPtr.Seconds, int64(txTimeAsPtr.Nanos)).String()
	return timeStr, nil
}

func main() {

    // Create a new Smart Contract
    err := shim.Start(new(MedicineContract))
    if err != nil {
        logger.Error("Error creating new Medicine Contracts:",err)
        
    }
}



