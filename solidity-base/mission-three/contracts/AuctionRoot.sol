// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import {CCIPReceiver} from "@chainlink/contracts-ccip/src/v0.8/ccip/applications/CCIPReceiver.sol";
import {Client} from "@chainlink/contracts-ccip/src/v0.8/ccip/libraries/Client.sol";
import {Auction} from "./Auction.sol"; // Explicitly import Auction to avoid conflicts

contract AuctionRoot is CCIPReceiver {

    event CrossChainBidReceived(
        address indexed auctionAddress,
        address indexed bidder,
        uint256 bidAmount
    );

    constructor(address _ccipRouter) CCIPReceiver(_ccipRouter) {}

    /**
     * @dev Handle incoming cross-chain bids.
     */
    function _ccipReceive(
        Client.Any2EVMMessage memory message
    ) internal override {
        // Decode the payload
        (address auctionAddress, address bidder, uint256 bidAmount) = abi.decode(
            message.data,
            (address, address, uint256)
        );

        // Process the bid on the remote chain
        //Auction auction = address(auctionAddress);
        //require(address(auction) != address(0), "Auction does not exist");

        // Place the bid on the auction
        Auction(auctionAddress).placeBid{value: bidAmount}(bidAmount, bidder);

        emit CrossChainBidReceived(auctionAddress, bidder, bidAmount);
    }

}
