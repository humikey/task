// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract reverse_string{

    function reverse(string memory str) public pure returns (string memory) {
        bytes memory strb = bytes(str);
        bytes memory output = new bytes(strb.length);

        for (uint i = 0; i < strb.length;i++){
            output[i] = strb[strb.length -i -1];
        }

        return string(output);
    }

    // 固定字节数组转动态数组
    function Todynamic(bytes6 memory name) view public returns(bytes){
        //return bytes(name);
        bytes memory newName = new bytes(name.length);

        //for循环挨个赋值
        for(uint i = 0;i<name.length;i++){
           newName[i] =  name[i];
        }
        return newName;
    }
}