const {
    FileSystemWallet,
    Gateway,
    X509WalletMixin
  } = require("fabric-network");
  const path = require("path");
  const fs = require('fs');

  

  exports.addToWallet = async (req,res) => {
    let certificatePath = req.body.certificatePath;
    let privateKeyPath = req.body.privateKeyPath;
    let companyRole = req.body.companyRole;
    

    try {

      // Fetch the credentials from our previously generated Crypto Materials required to create identity
      const certificate = fs.readFileSync(certificatePath).toString();
      
      // IMPORTANT: Change the private key name to the key generated on your computer
      const privatekey = fs.readFileSync(privateKeyPath).toString();
  
      // Load credentials into wallet
      const identityLabel = companyRole.toUpperCase() + '_ADMIN';
      const companyMsp = companyRole.charAt(0).toUpperCase() + companyRole.slice(1)
      console.log("companyymsp : ", companyMsp)
         const identity = X509WalletMixin.createIdentity(companyMsp + 'MSP', certificate, privatekey);
  
      const wallet = new FileSystemWallet('./wallet/'+companyRole);
      await wallet.import(identityLabel, identity);

      res.json({
        status: "success",
        message: `Added identity to wallet`
      });
  
    } catch (error) {
      console.log(`Error adding to wallet. ${error}`);
      res.json({
        status: "failed",
        message: `Error adding to wallet`
      });
      
    }

  }

  exports.registerCompany = async (req, res) => {
    let companyName = req.body.companyName;
    let location = req.body.location;
    let companyRole = req.body.companyRole;
    

    try{
        const wallet = new FileSystemWallet("./wallet/" + companyRole);
        const connectionProfilePath = path.resolve(__dirname, "..", "..","..", "connection-"+companyRole+".json");
        // Check to see if we've already enrolled the user.

        
        const fabricIdentity= companyRole.toUpperCase()+"_ADMIN"
        const userExists = await wallet.exists(fabricIdentity);
        if (!userExists) {
            console.log(`An identity for the company ${companyRole} does not exist in the wallet`);
            console.log('Add the identity to wallet');
            return;
        }

        console.log("connecting to getway")
        // Create a new gateway for connecting to our peer node.
        const gateway = new Gateway();
        await gateway.connect(connectionProfilePath, { wallet, identity: fabricIdentity, discovery: { enabled: true, asLocalhost: true } });

        // Get the network (channel) our contract is deployed to.
        const network = await gateway.getNetwork('pharmachannel');
        console.log("connected to getway")
        // Get the contract from the network.
        const contract = network.getContract('medicinecontract');

        console.log("received contract")
        const result= await contract.submitTransaction('registerCompany',companyName, location, companyRole);
        console.log('Transaction has been submitted');

        // Disconnect from the gateway.
        await gateway.disconnect();

        res.json({
            status: "success",
            message: `Registered company ${companyName} successfully`,
            data:result.toString()
          });

    }
    catch (error) {
        res.json({
          status: "failed",
          message: `Failed to register company ${companyName}: ${error}`
        });
    }

    
}

exports.addMedicine = async (req, res) => {
  let drugName = req.body.drugName;
  let serialNo = req.body.serialNo;
  let mfgDate = req.body.mfgDate;
  let expDate = req.body.expDate;
  let owner = req.body.owner;
  let manufactureId = req.body.manufactureId;
  let manufactureName = req.body.manufactureName;
  let companyRole = req.body.companyRole;

  try{
      const wallet = new FileSystemWallet("./wallet/" + companyRole);
      const connectionProfilePath = path.resolve(__dirname, "..", "..","..", "connection-"+companyRole+".json");
      // Check to see if we've already enrolled the user.
      const fabricIdentity= companyRole.toUpperCase()+"_ADMIN"
      const userExists = await wallet.exists(fabricIdentity);
      if (!userExists) {
          console.log(`An identity for the company ${companyRole} does not exist in the wallet`);
          console.log('Add the identity to wallet');
          return;
      }

      // Create a new gateway for connecting to our peer node.
      const gateway = new Gateway();
      await gateway.connect(connectionProfilePath, { wallet, identity: fabricIdentity, discovery: { enabled: true, asLocalhost: true } });

      // Get the network (channel) our contract is deployed to.
      const network = await gateway.getNetwork('pharmachannel');

      // Get the contract from the network.
      const contract = network.getContract('medicinecontract');

     
      const result=await contract.submitTransaction('addMedicine',drugName, serialNo, mfgDate, expDate, owner, manufactureId, manufactureName);
      console.log('Transaction has been submitted');

      // Disconnect from the gateway.
      await gateway.disconnect();

      res.json({
          status: "success",
          message: `Added medicine successfully`,
          data:result.toString()
        });

  }
  catch (error) {
      res.json({
        status: "failed",
        message: `Failed to add medicine: ${error}`
      });
  }
}

exports.createPurchaseOrder = async (req, res) => {
  let buyerId = req.body.buyerId;
  let sellerId = req.body.sellerId;
  let drugName = req.body.drugName;
  let quantity = req.body.quantity;
  let companyRole = req.body.companyRole;
 

  try{
      const wallet = new FileSystemWallet("./wallet/" + companyRole);
      const connectionProfilePath = path.resolve(__dirname, "..", "..","..", "connection-"+companyRole+".json");
      // Check to see if we've already enrolled the user.
      const fabricIdentity= companyRole.toUpperCase()+"_ADMIN"
      const userExists = await wallet.exists(fabricIdentity);
      if (!userExists) {
          console.log(`An identity for the company ${companyRole} does not exist in the wallet`);
          console.log('Add the identity to wallet');
          return;
      }

      // Create a new gateway for connecting to our peer node.
      const gateway = new Gateway();
      await gateway.connect(connectionProfilePath, { wallet, identity: fabricIdentity, discovery: { enabled: true, asLocalhost: true } });

      // Get the network (channel) our contract is deployed to.
      const network = await gateway.getNetwork('pharmachannel');

      // Get the contract from the network.
      const contract = network.getContract('medicinecontract');

      
      const result=await contract.submitTransaction('createPurchaseOrder',buyerId, sellerId, drugName, quantity);
      console.log('Transaction has been submitted');

      // Disconnect from the gateway.
      await gateway.disconnect();

      res.json({
          status: "success",
          message: `Created Purchase order successfully`,
          data:result.toString()
        });

  }
  catch (error) {
      res.json({
        status: "failed",
        message: `Failed to create purchase order: ${error}`
      });
  }
}

exports.createShipment = async (req, res) => {
  let buyerId = req.body.buyerId;
  let drugName = req.body.drugName;
  let listOfAssets = req.body.listOfAssets;
  let transporterId = req.body.transporterId;
  let companyRole = req.body.companyRole;
 

  try{
      const wallet = new FileSystemWallet("./wallet/" + companyRole);
      const connectionProfilePath = path.resolve(__dirname, "..", "..","..", "connection-"+companyRole+".json");
      
      // Check to see if we've already enrolled the user.
      const fabricIdentity= companyRole.toUpperCase()+"_ADMIN"
      const userExists = await wallet.exists(fabricIdentity);
      if (!userExists) {
          console.log(`An identity for the company ${companyRole} does not exist in the wallet`);
          console.log('Add the identity to wallet');
          return;
      }

      // Create a new gateway for connecting to our peer node.
      const gateway = new Gateway();
      await gateway.connect(connectionProfilePath, { wallet, identity: fabricIdentity, discovery: { enabled: true, asLocalhost: true } });

      // Get the network (channel) our contract is deployed to.
      const network = await gateway.getNetwork('pharmachannel');

      // Get the contract from the network.
      const contract = network.getContract('medicinecontract');

      
      const result=await contract.submitTransaction('createShipment',buyerId, drugName, listOfAssets, transporterId);
      console.log('Transaction has been submitted');

      // Disconnect from the gateway.
      await gateway.disconnect();

      res.json({
          status: "success",
          message: `Created shipment successfully`,
          data:result.toString()
        });

  }
  catch (error) {
      res.json({
        status: "failed",
        message: `Failed to create shipment: ${error}`
      });
  }
}


exports.updateShipment = async (req, res) => {
  let buyerId = req.body.buyerId;
  let drugName = req.body.drugName;
  let transporterId = req.body.transporterId;
  let companyRole = req.body.companyRole;
 

  try{
      const wallet = new FileSystemWallet("./wallet/" + companyRole);
      const connectionProfilePath = path.resolve(__dirname, "..", "..","..", "connection-"+companyRole+".json");
      // Check to see if we've already enrolled the user.
      const fabricIdentity= companyRole.toUpperCase()+"_ADMIN"
      const userExists = await wallet.exists(fabricIdentity);
      if (!userExists) {
          console.log(`An identity for the company ${companyRole} does not exist in the wallet`);
          console.log('Add the identity to wallet');
          return;
      }

      // Create a new gateway for connecting to our peer node.
      const gateway = new Gateway();
      await gateway.connect(connectionProfilePath, { wallet, identity: fabricIdentity, discovery: { enabled: true, asLocalhost: true } });

      // Get the network (channel) our contract is deployed to.
      const network = await gateway.getNetwork('pharmachannel');

      // Get the contract from the network.
      const contract = network.getContract('medicinecontract');

      
      const result=await contract.submitTransaction('updateShipment',buyerId, drugName,  transporterId);
      console.log('Transaction has been submitted');

      // Disconnect from the gateway.
      await gateway.disconnect();

      res.json({
          status: "success",
          message: `Updated shipment successfully`,
          data: result.toString()
        });

  }
  catch (error) {
      res.json({
        status: "failed",
        message: `Failed to update shipment: ${error}`
      });
  }
}

exports.retailMedicine = async (req, res) => {

  let drugName = req.body.drugName;
  let serialNo = req.body.serialNo;
  let retailerId = req.body.retailerId;
  let customerAadhar = req.body.customerAadhar;
  let companyRole = req.body.companyRole;
 

  try{
      const wallet = new FileSystemWallet("./wallet/" + companyRole);
      const connectionProfilePath = path.resolve(__dirname, "..", "..","..", "connection-"+companyRole+".json");
      // Check to see if we've already enrolled the user.
      const fabricIdentity= companyRole.toUpperCase()+"_ADMIN"
      const userExists = await wallet.exists(fabricIdentity);
      if (!userExists) {
          console.log(`An identity for the company ${companyRole} does not exist in the wallet`);
          console.log('Add the identity to wallet');
          return;
      }

      // Create a new gateway for connecting to our peer node.
      const gateway = new Gateway();
      await gateway.connect(connectionProfilePath, { wallet, identity: fabricIdentity, discovery: { enabled: true, asLocalhost: true } });

      // Get the network (channel) our contract is deployed to.
      const network = await gateway.getNetwork('pharmachannel');

      // Get the contract from the network.
      const contract = network.getContract('medicinecontract');

      
      const result=await contract.submitTransaction('retailMedicine',drugName, serialNo,retailerId, customerAadhar);
      console.log('Transaction has been submitted');

      // Disconnect from the gateway.
      await gateway.disconnect();

      res.json({
          status: "success",
          message: `Medicine sold entry by retailer created successfully`,
          data: result.toString()
        });

  }
  catch (error) {
      res.json({
        status: "failed",
        message: `Failed to create entry for medicine sold by retailer: ${error}`
      });
  }
}


exports.viewHistory = async (req, res) => {

  if (!req.params) {
		res.status(404).json({
			message: "Parameters are not supplied"
        });
    }

    else{

  let serialNo = req.params.serialNo;
  let drugName = req.params.drugName;
  let companyRole = req.params.companyRole;

 

  try{
      const wallet = new FileSystemWallet("./wallet/" + companyRole);
      const connectionProfilePath = path.resolve(__dirname, "..", "..","..", "connection-"+companyRole+".json");
      // Check to see if we've already enrolled the user.
      const fabricIdentity= companyRole.toUpperCase()+"_ADMIN"
      const userExists = await wallet.exists(fabricIdentity);
      if (!userExists) {
          console.log(`An identity for the company ${companyRole} does not exist in the wallet`);
          console.log('Add the identity to wallet');
          return;
      }

      // Create a new gateway for connecting to our peer node.
      const gateway = new Gateway();
      await gateway.connect(connectionProfilePath, { wallet, identity: fabricIdentity, discovery: { enabled: true, asLocalhost: true } });

      // Get the network (channel) our contract is deployed to.
      const network = await gateway.getNetwork('pharmachannel');

      // Get the contract from the network.
      const contract = network.getContract('medicinecontract');

      
      const result= await contract.evaluateTransaction('viewHistory',serialNo,drugName );
      
        
        console.log(`Transaction has been evaluated, result is: ${result.toString()}`);

      

      // Disconnect from the gateway.
      await gateway.disconnect();

      res.json({
          status: "success",
          data: result.toString()
        });

  }
  catch (error) {
      res.json({
        status: "failed",
        message: `Failed to retrieve history for medicine: ${error}`
      });
  }
}
}


exports.viewMedicineCurrentState = async (req, res) => {

  if (!req.params) {
		res.status(404).json({
			message: "Parameters are not supplied"
        });
    }

    else{

  let serialNo = req.params.serialNo;
  let drugName = req.params.drugName;
  let companyRole = req.params.companyRole;

 

  try{
      const wallet = new FileSystemWallet("./wallet/" + companyRole);
      const connectionProfilePath = path.resolve(__dirname, "..", "..","..", "connection-"+companyRole+".json");
      // Check to see if we've already enrolled the user.
      const fabricIdentity= companyRole.toUpperCase()+"_ADMIN"
      const userExists = await wallet.exists(fabricIdentity);
      if (!userExists) {
          console.log(`An identity for the company ${companyRole} does not exist in the wallet`);
          console.log('Add the identity to wallet');
          return;
      }

      // Create a new gateway for connecting to our peer node.
      const gateway = new Gateway();
      await gateway.connect(connectionProfilePath, { wallet, identity: fabricIdentity, discovery: { enabled: true, asLocalhost: true } });

      // Get the network (channel) our contract is deployed to.
      const network = await gateway.getNetwork('pharmachannel');

      // Get the contract from the network.
      const contract = network.getContract('medicinecontract');

      
      const result= await contract.evaluateTransaction('viewMedicineCurrentState',serialNo,drugName );
      
        
        console.log(`Transaction has been evaluated, result is: ${result.toString()}`);

      

      // Disconnect from the gateway.
      await gateway.disconnect();

      res.json({
          status: "success",
          data: result.toString()
        });

  }
  catch (error) {
      res.json({
        status: "failed",
        message: `Failed to retrieve medicine current state: ${error}`
      });
  }
}
}



exports.getMedicineByManufacturer = async (req, res) => {
  if (!req.params) {
		res.status(404).json({
			message: "Parameters are not supplied"
        });
    }

    else{
 
  let companyName = req.params.companyName;
  let companyRole = req.params.companyRole;

 

  try{
      const wallet = new FileSystemWallet("./wallet/" + companyRole);
      const connectionProfilePath = path.resolve(__dirname, "..", "..","..", "connection-"+companyRole+".json");
      // Check to see if we've already enrolled the user.
      const fabricIdentity= companyRole.toUpperCase()+"_ADMIN"
      const userExists = await wallet.exists(fabricIdentity);
      if (!userExists) {
          console.log(`An identity for the company ${companyRole} does not exist in the wallet`);
          console.log('Add the identity to wallet');
          return;
      }

      // Create a new gateway for connecting to our peer node.
      const gateway = new Gateway();
      await gateway.connect(connectionProfilePath, { wallet, identity: fabricIdentity, discovery: { enabled: true, asLocalhost: true } });

      // Get the network (channel) our contract is deployed to.
      const network = await gateway.getNetwork('pharmachannel');

      // Get the contract from the network.
      const contract = network.getContract('medicinecontract');

     
      const result= await contract.evaluateTransaction('getMedicineByManufacturer',companyName );
     
        
        console.log(`Transaction has been evaluated, result is: ${result.toString()}`);

      

      // Disconnect from the gateway.
      await gateway.disconnect();

      res.json({
          status: "success",
          data: result.toString()
        });

  }
  catch (error) {
      res.json({
        status: "failed",
        message: `Failed to retrieve medicines from manufacturer: ${error}`
      });
  }
}
}

exports.getCompany = async (req, res) => {

  if (!req.params) {
		res.status(404).json({
			message: "Parameters are not supplied"
        });
    }

    else{
  let companyName = req.params.companyName;
  let companyRole = req.params.companyRole;

 

  try{
      const wallet = new FileSystemWallet("./wallet/" + companyRole);
      const connectionProfilePath = path.resolve(__dirname, "..", "..","..", "connection-"+companyRole+".json");
      // Check to see if we've already enrolled the user.
      const fabricIdentity= companyRole.toUpperCase()+"_ADMIN"
      const userExists = await wallet.exists(fabricIdentity);
      if (!userExists) {
          console.log(`An identity for the company ${companyRole} does not exist in the wallet`);
          console.log('Add the identity to wallet');
          return;
      }

      // Create a new gateway for connecting to our peer node.
      const gateway = new Gateway();
      await gateway.connect(connectionProfilePath, { wallet, identity: fabricIdentity, discovery: { enabled: true, asLocalhost: true } });

      // Get the network (channel) our contract is deployed to.
      const network = await gateway.getNetwork('pharmachannel');

      // Get the contract from the network.
      const contract = network.getContract('medicinecontract');

      
      const result= await contract.evaluateTransaction('getCompany',companyName );
     
        
        console.log(`Transaction has been evaluated, result is: ${result.toString()}`);

      

      // Disconnect from the gateway.
      await gateway.disconnect();

      res.json({
          status: "success",
          data: result.toString()
        });

  }
  catch (error) {
      res.json({
        status: "failed",
        message: `Failed to retrieve company details: ${error}`
      });
  }
}
}