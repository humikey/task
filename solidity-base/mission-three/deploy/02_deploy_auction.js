// scripts/deployAuction.js

const { getNamedAccounts } = require("hardhat");

module.exports = async ({ getNamedAccounts, deployments }) => {
    const { firstAccount } = await getNamedAccounts()
    const { deploy, log } = deployments

    const auctionNFT = await deployments.get("AuctionNFT")

    // Required parameters - these would need to be configured based on your network
    const seller = firstAccount; // The seller account
    const nftContract = auctionNFT.address; // Address of the NFT contract (e.g. CryptoKitties, etc.)
    const tokenId = 1; // Token ID of the NFT to be auctioned
    const paymentToken = "0x0000000000000000000000000000000000000000"; // ETH (address(0)) or ERC20 token address
    const startingPrice = ethers.utils.parseEther("100"); // Starting price in USD with 18 decimals
    const ethUsdPriceFeed = "0x694AA1769357215DE4FAC081bf1f309aDC325306"; // Mainnet ETH/USD price feed
    const erc20UsdPriceFeed = "0x694AA1769357215DE4FAC081bf1f309aDC325306"; // Price feed for the ERC20 token if used (or can be any address if using ETH)

    log("Deploying the Auction contract")
    await deploy("Auction", {
        contract: "Auction",
        from: firstAccount,
        log: true,
        args: [seller, nftContract, tokenId, paymentToken, startingPrice, ethUsdPriceFeed, erc20UsdPriceFeed]
    })

    console.log("Auction deployed to:", auction.address);
    console.log("Seller:", seller);
    console.log("NFT Contract:", nftContract);
    console.log("Token ID:", tokenId.toString());
    console.log("Payment Token:", paymentToken);
    console.log("Starting Price (USD):", ethers.utils.formatEther(startingPrice));
    console.log("ETH/USD Price Feed:", ethUsdPriceFeed);
    console.log("ERC20/USD Price Feed:", erc20UsdPriceFeed);

    log("Auction is deployed!")
}

module.exports.tags = ["all", "sourcechain"]