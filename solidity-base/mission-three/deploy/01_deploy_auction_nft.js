const { getNamedAccounts } = require("hardhat");

module.exports = async({getNamedAccounts, deployments}) => {
    const {firstAccount} = await getNamedAccounts()
    const {deploy, log} = deployments
    
    log("Deploying the nft contract")
    await deploy("AuctionNFT", {
        contract: "AuctionNFT",
        from: firstAccount,
        log: true,
        args: []
    })
    log("AuctionNFT is deployed!")
}

module.exports.tags = ["all", "sourcechain"]