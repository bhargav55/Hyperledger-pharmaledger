# PharmaLedger

The application tracks the medicine supply chain right from origin at the manufacturer till the medicine gets sold to any
customer by retailer. Each flow in the process is updated in the ledger which is immutable and is recorded as a transaction.
This way any medicine can be traced along with the provenance. This smart contract layer of the  blockchain automates the flow.
This solution also reduces the paper work involved in the process

The organizations in the network are :

1. Manufacture: 

    a. Add new medicine to the ledger
    b. Creates shipment based on the product order/ purchase order from distributor

2.  Distributor

    a. Raise purchase order for a medicine to the manufacturer
    b. Creates shipment based on the product order/ purchase order from retailer
    
3. Transporter

    a. Update shipment 
    
4. Retailer

    a. Raise purchase order for a medicine to the distributor
    b. Sell medicine to the customer with customer details
    
Any Org member can retrieve status of the medicine and get the history records for the medicine.

