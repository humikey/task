// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "@openzeppelin/contracts/token/ERC721/ERC721.sol";
import "@openzeppelin/contracts/access/Ownable.sol";

contract AuctionNFT is ERC721, Ownable {
    uint256 private _tokenIdCounter;

    constructor() ERC721("MyNFT", "MNFT") Ownable(msg.sender) onlyOwner{
        _tokenIdCounter = 1; // Start token IDs from 1
    }

    /**
     * @dev Mint a new NFT to the specified address.
     * @param to The address to receive the NFT.
     */
    function mint(address to) public onlyOwner {
        uint256 tokenId = _tokenIdCounter;
        _tokenIdCounter++;
        _safeMint(to, tokenId);
    }

    /**
     * @dev Get the current token ID counter.
     * @return The next token ID to be minted.
     */
    function getCurrentTokenId() public view returns (uint256) {
        return _tokenIdCounter;
    }
}