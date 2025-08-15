// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract voting {
    // 用一个mapping来存储候选人得票信息
    mapping(string => uint256) private votes;

    // 数组来确定所有存储的候选人
    string[] private condidates;

    // 防止重复添加候选人
    mapping(string => bool) private is_candidate;

    // 一个vote函数，允许用户投票给某个候选人
    function vote(string calldata name) public {
        if (!is_candidate[name]) {
            condidates.push(name);
            is_candidate[name] = true;
        }
        votes[name]++;
    }

    //一个getVotes函数，返回某个候选人的得票数
    function getVotes(string calldata name) public view returns(uint256){
        return votes[name];
    }

    // 一个resetVotes函数，重置所有候选人的得票数
    function resteVotes() public {
        for(uint i = 0; i < condidates.length;i++){
            votes[condidates[i]] = 0;
        }
    }
}