// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract binary_search{


    function binarySearch(uint[] memory arr, uint target) public pure returns (int result){
        int left = 0;
        int right = int(arr.length) -1;
        while(left <= right){
            int mid = left + (right -left)/2;
            uint midVal = arr[uint(mid)];
            if(midVal == target){
                return mid;
            } else if (midVal < target){
                left = mid + 1;
            } else {
                right = mid -1;
            }
        }

        return -1;
    }
}