// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract BeggingContract {
    address public owner;

    // 记录每个捐赠者的总捐赠金额
    mapping(address => uint256) public donations;

    // 记录所有捐赠者地址
    address[] public donors;

    // 事件记录
    event DonationReceived(address indexed donor, uint256 amount);
    event Withdrawal(address indexed owner, uint256 amount);

    // 构造函数，部署者即为合约所有者
    constructor() {
        owner = msg.sender;
    }

    // 任何人都可以调用此函数向合约捐赠 ETH
    receive() external payable {
        require(msg.value > 0, "Donation must be more than 0");

        // 如果是第一次捐赠，则记录地址
        if (donations[msg.sender] == 0) {
            donors.push(msg.sender);
        }

        donations[msg.sender] += msg.value;

        emit DonationReceived(msg.sender, msg.value);
    }

    // 仅合约所有者可以调用，用于提取所有捐赠资金
    function withdraw() external {
        require(msg.sender == owner, "Only owner can withdraw");
        uint256 balance = address(this).balance;
        require(balance > 0, "No funds to withdraw");

        payable(owner).transfer(balance);

        emit Withdrawal(owner, balance);
    }

    // 获取所有捐赠者地址列表
    function getAllDonors() external view returns (address[] memory) {
        return donors;
    }

    // 合约当前余额
    function getBalance() external view returns (uint256) {
        return address(this).balance;
    }
}
