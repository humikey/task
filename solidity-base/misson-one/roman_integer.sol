// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract roman_integer{

    function roman_to_integer(string memory roman) public pure returns (uint256){
        bytes memory b = bytes(roman);
        uint total = 0;
        uint prev = 0;

        for (uint i = b.length; i > 0; i--) {
            uint val = romanCharToUint(b[i - 1]);

            // 如果当前字符的值小于之前字符，减去（如 IV = 5-1）
            if (val < prev) {
                total -= val;
            } else {
                total += val;
            }

            prev = val;
        }

        return total;

    }

    function romanCharToUint(bytes1 c) internal pure returns (uint) {
        if (c == "I") return 1;
        if (c == "V") return 5;
        if (c == "X") return 10;
        if (c == "L") return 50;
        if (c == "C") return 100;
        if (c == "D") return 500;
        if (c == "M") return 1000;
        revert("Invalid Roman character");
    }

    function intToRoman(uint256 num) public pure returns (string memory) {
        require(num >= 1 && num <= 3999, "Out of range");

        // 数值和对应罗马字符
        uint16[13] memory values = [1000, 900, 500, 400, 100, 90, 50, 40, 10, 9, 5, 4, 1];
        string[13] memory symbols = ["M", "CM","D","CD","C","XC","L","XL","X","IX","V","IV","I"];

        bytes memory result;
        for (uint i = 0; i < values.length; i++) {
            while (num >= values[i]) {
                num -= values[i];
                result = bytes.concat(result, bytes(symbols[i]));
            }
        }

        return string(result);
    }

}