const express = require('express');

const router = express.Router();

//const mongoose = require('mongoose');

//const Contract = require('../model/contracts');
const chaincodeMethods = require('../controllers/chaincodemethods');

router.get('/', (req,res,next) => {
    res.json({
        result: "ok",
        message: "initialized"
    })
});

module.exports = router;


router.post('/registerCompany', chaincodeMethods.registerCompany);
router.post('/addToWallet', chaincodeMethods.addToWallet);
router.post('/addMedicine', chaincodeMethods.addMedicine);
router.post('/createPurchaseOrder', chaincodeMethods.createPurchaseOrder);
router.post('/createShipment', chaincodeMethods.createShipment);
router.post('/updateShipment', chaincodeMethods.updateShipment);
router.post('/retailMedicine', chaincodeMethods.retailMedicine);
router.get('/viewHistory/:serialNo/:drugName/:companyRole', chaincodeMethods.viewHistory);
router.get('/viewMedicineCurrentState/:serialNo/:drugName/:companyRole', chaincodeMethods.viewMedicineCurrentState);
router.get('/getMedicineByManufacturer/:companyName/:companyRole', chaincodeMethods.getMedicineByManufacturer);
router.get('/getCompany/:companyName/:companyRole', chaincodeMethods.getCompany);
