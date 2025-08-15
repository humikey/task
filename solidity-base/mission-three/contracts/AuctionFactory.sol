// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "./Auction.sol";

contract AuctionFactory {
    // Array to store all auction addresses
    address[] public allAuctions;

    // Event emitted when a new auction is created
    event AuctionCreated(address indexed auctionAddress, address indexed seller, address nftContract, uint256 tokenId);

    /**
     * @dev Create a new auction instance.
     * @param nftContract Address of the NFT contract.
     * @param tokenId ID of the NFT to be auctioned.
     * @param paymentToken Address of the ERC20 token for payment, or address(0) for ETH.
     * @param startingPrice Minimum bid price.
     */
    function createAuction(
        address nftContract,
        uint256 tokenId,
        address paymentToken,
        uint256 startingPrice,
        address ethUsdPriceFeed,
        address erc20UsdPriceFeed
    ) external {
        // Deploy a new Auction contract
        Auction newAuction = new Auction(msg.sender, nftContract, tokenId, paymentToken, startingPrice,
            ethUsdPriceFeed, erc20UsdPriceFeed);

        // Store the auction address
        allAuctions.push(address(newAuction));

        // Emit event
        emit AuctionCreated(address(newAuction), msg.sender, nftContract, tokenId);
    }

    /**
     * @dev Get the total number of auctions created.
     * @return The total number of auctions.
     */
    function getAllAuctions() external view returns (address[] memory) {
        return allAuctions;
    }
}