// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

// This contract is an optimized version of the Math contract
// It uses pure functions to avoid state changes and gas costs associated with storage
contract MathOptimized2 {
    function add(uint256 a, uint256 b) public pure returns (uint256) {
        unchecked {
            return a + b;
        }
    }

    function sub(uint256 a, uint256 b) public pure returns (uint256) {
        unchecked {
            return a - b;
        }
    }
}
