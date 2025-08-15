const { getNamedAccounts, network } = require("hardhat")
const { developmentChains, networkConfig } = require("../helper-hardhat-config")

module.exports = async ({ getNamedAccounts, deployments }) => {
    const { firstAccount } = await getNamedAccounts()
    const { deploy, log } = deployments

    // get parameters for constructor
    let sourceChainRouter
    // let auctionAddress
    // let bidder
    // let nftAddr
    // get router and linktoken based on network
    sourceChainRouter = networkConfig[network.config.chainId].router
    linkToken = networkConfig[network.config.chainId].linkToken
    log(`non local environment: sourcechain router: ${sourceChainRouter}, link token: ${linkToken}`)

    log("deploying the AuctionRoot")
    await deploy("AuctionRoot", {
        contract: "AuctionRoot",
        from: firstAccount,
        log: true,
        args: [sourceChainRouter]
    })
    log("AuctionRoot deployed")
}

module.exports.tags = ["all", "sourcechain"]

// (uint256 auctionId, address bidder, uint256 bidAmount)