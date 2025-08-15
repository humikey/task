const { getNamedAccounts, deployments, ethers } = require("hardhat")
const { expect } = require("chai")

let firstAccount
let nft
let auctionNft
let auctionRoot
let auctionRemote

before(async function () {
    firstAccount = (await getNamedAccounts()).firstAccount
    await deployments.fixture(["all"])
    nft = await ethers.getContract("AuctionNFT", firstAccount)
    auctionNft = await ethers.getContract("AuctionNFT", firstAccount)
    auctionRoot = await ethers.getContract("AuctionRoot", firstAccount)
    auctionRemote = await ethers.getContract("AuctionRemote", firstAccount)
})

describe("test if the nft can be minted successfully",
    async function () {
        it("test if the owner of nft is minter",
            async function () {
                // get nft
                await nft.safeMint(firstAccount)
                // check the owner
                const ownerOfNft = await nft.ownerOf(0)
                expect(ownerOfNft).to.equal(firstAccount)
            })
    })
