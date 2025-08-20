// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

contract Math {
    uint256 public lastResult;

    function add(uint256 a, uint256 b) public returns (uint256) {
        uint256 result = a + b;
        lastResult = result;
        return result;
    }

    function sub(uint256 a, uint256 b) public returns (uint256) {
        require(a >= b, "Underflow");
        uint256 result = a - b;
        lastResult = result;
        return result;
    }
}
