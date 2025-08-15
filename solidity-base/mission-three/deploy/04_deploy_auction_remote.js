const { getNamedAccounts, network } = require("hardhat")
const { developmentChains, networkConfig } = require("../helper-hardhat-config")

module.exports = async ({ getNamedAccounts, deployments }) => {
    const { firstAccount } = await getNamedAccounts()
    const { deploy, log } = deployments

    // get parameters for constructor
    let destChainRouter
    // let auctionAddress
    // let bidder
    // let nftAddr
    // get router and linktoken based on network
    destChainRouter = networkConfig[network.config.chainId].router
    linkToken = networkConfig[network.config.chainId].linkToken
    log(`non local environment: sourcechain router: ${destChainRouter}, link token: ${linkToken}`)

    log("deploying the AuctionRemote")
    await deploy("AuctionRemote", {
        contract: "AuctionRemote",
        from: firstAccount,
        log: true,
        args: [destChainRouter]
    })
    log("AuctionRemote deployed")
}

module.exports.tags = ["all", "destchain"]
