// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract merge_array{

function mergeSorted(uint[] memory a, uint[] memory b) public pure returns (uint[] memory) {
        uint lenA = a.length;
        uint lenB = b.length;
        uint total = lenA + lenB;

        uint[] memory result = new uint[](total);

        uint i = 0; // index for a
        uint j = 0; // index for b
        uint k = 0; // index for result

        while (i < lenA && j < lenB) {
            if (a[i] <= b[j]) {
                result[k++] = a[i++];
            } else {
                result[k++] = b[j++];
            }
        }

        // 处理剩余元素
        while (i < lenA) {
            result[k++] = a[i++];
        }

        while (j < lenB) {
            result[k++] = b[j++];
        }

        return result;
    }

}