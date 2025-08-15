// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "@openzeppelin/contracts/token/ERC721/IERC721.sol";
import "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import "@chainlink/contracts/src/v0.8/shared/interfaces/AggregatorV3Interface.sol";

contract Auction {
    address public seller;
    address public nftContract;
    uint256 public tokenId;
    address public paymentToken; // Address of ERC20 token, or address(0) for ETH
    uint256 public startingPrice; // In USD (scaled by 1e18)
    uint256 public highestBid; // In USD (scaled by 1e18)
    address public highestBidder;
    bool public active;

    AggregatorV3Interface public ethUsdPriceFeed;
    AggregatorV3Interface public erc20UsdPriceFeed;

    event BidPlaced(address indexed bidder, uint256 amountInUsd);
    event AuctionEnded(address indexed winner, uint256 amountInUsd);

    constructor(
        address _seller,
        address _nftContract,
        uint256 _tokenId,
        address _paymentToken,
        uint256 _startingPrice,
        address _ethUsdPriceFeed,
        address _erc20UsdPriceFeed
    ) {
        seller = _seller;
        nftContract = _nftContract;
        tokenId = _tokenId;
        paymentToken = _paymentToken;
        startingPrice = _startingPrice;
        active = true;

        ethUsdPriceFeed = AggregatorV3Interface(_ethUsdPriceFeed);
        erc20UsdPriceFeed = AggregatorV3Interface(_erc20UsdPriceFeed);

        // Transfer NFT to this contract
        IERC721(nftContract).transferFrom(seller, address(this), tokenId);
    }

    function placeBid(uint256 bidAmount, address remoteAddress) external payable {
        require(active, "Auction is not active");
        require(msg.sender != seller, "Seller cannot bid");

        address bidder = msg.sender;
        if(paymentToken != address(0)){
            bidder = remoteAddress; // Use remote address for cross-chain bids
        }

        uint256 bidAmountInUsd = paymentToken == address(0)
            ? _convertEthToUsd(msg.value)
            : _convertErc20ToUsd(bidAmount);

        require(bidAmountInUsd > highestBid, "Bid must be higher than the current highest bid");

        // Refund previous highest bidder
        if (highestBidder != address(0)) {
            if (paymentToken == address(0)) {
                payable(highestBidder).transfer(_convertUsdToEth(highestBid));
            } else {
                IERC20(paymentToken).transfer(highestBidder, _convertUsdToErc20(highestBid));
            }
        }

        // Update highest bid
        if (paymentToken == address(0)) {
            highestBid = bidAmountInUsd;
        } else {
            IERC20(paymentToken).transferFrom(bidder, address(this), bidAmount);
            highestBid = bidAmountInUsd;
        }
        highestBidder = bidder;

        emit BidPlaced(bidder, highestBid);
    }

    function endAuction() external {
        require(active, "Auction is not active");
        require(msg.sender == seller, "Only the seller can end the auction");

        active = false;

        if (highestBidder != address(0)) {
            // Transfer NFT to the highest bidder
            IERC721(nftContract).transferFrom(address(this), highestBidder, tokenId);

            // Transfer funds to the seller
            if (paymentToken == address(0)) {
                payable(seller).transfer(_convertUsdToEth(highestBid));
            } else {
                IERC20(paymentToken).transfer(seller, _convertUsdToErc20(highestBid));
            }

            emit AuctionEnded(highestBidder, highestBid);
        } else {
            // No bids, return NFT to the seller
            IERC721(nftContract).transferFrom(address(this), seller, tokenId);
        }
    }

    function _convertEthToUsd(uint256 ethAmount) internal view returns (uint256) {
        (, int256 price, , , ) = ethUsdPriceFeed.latestRoundData();
        require(price > 0, "Invalid ETH/USD price");
        return (ethAmount * uint256(price)) / 1e18;
    }

    function _convertErc20ToUsd(uint256 erc20Amount) internal view returns (uint256) {
        (, int256 price, , , ) = erc20UsdPriceFeed.latestRoundData();
        require(price > 0, "Invalid ERC20/USD price");
        return (erc20Amount * uint256(price)) / 1e18;
    }

    function _convertUsdToEth(uint256 usdAmount) internal view returns (uint256) {
        (, int256 price, , , ) = ethUsdPriceFeed.latestRoundData();
        require(price > 0, "Invalid ETH/USD price");
        return (usdAmount * 1e18) / uint256(price);
    }

    function _convertUsdToErc20(uint256 usdAmount) internal view returns (uint256) {
        (, int256 price, , , ) = erc20UsdPriceFeed.latestRoundData();
        require(price > 0, "Invalid ERC20/USD price");
        return (usdAmount * 1e18) / uint256(price);
    }
}