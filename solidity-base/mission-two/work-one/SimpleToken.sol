// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

// ERC20 是以太坊平台上最常见的一种 代币标准（Token Standard），全称是 Ethereum Request for Comments 20。它定义了一组 规范接口（标准函数和事件），用于在智能合约中创建可互换（fungible）的代币。

// 导入代币到钱包
// 打开 MetaMask，进入钱包界面
// 点击 “添加代币”
// 输入合约地址，手动填写符号 SIM 和精度 18
// 点击确认，你就能看到你代币余额了
contract SimpleToken {

    string name = "SimpleToken";
    string symble = "SIM";
    uint8 decimal = 8;

    address _owner;
    uint256 _totalSuply;

    mapping(address => uint256) _balanceOf;
    mapping(address => mapping(address => uint256)) _allowance;

    constructor(){
        _owner = msg.sender;
    }

    event Transfer(address indexed from, address indexed to, uint256 value);

    event Approval(
        address indexed owner,
        address indexed spender,
        uint256 value
    );

    function mint(address to, uint256 value) public {
        require(_owner == msg.sender,"not owner");
        _balanceOf[to] += value;
        _totalSuply += value;
        emit Transfer(address(0), to, value);
    }

    function balanceOf(address account) external view returns (uint256) {
        return _balanceOf[account];
    }

    function transfer(address to, uint256 value) external returns (bool) {
        require(value <= _balanceOf[msg.sender],"not enough balance");
        _balanceOf[msg.sender] -= value;
        _balanceOf[to] += value;
        emit Transfer(msg.sender, to, value);
        return true;
    }

    function allowance(address owner, address spender)
        external
        view
        returns (uint256)
    {
        return _allowance[owner][spender];
    }

    function approve(address spender, uint256 value) external returns (bool) {
        _allowance[msg.sender][spender] = value;
        emit Approval(msg.sender, spender, value);
        return true;
    }

    function transferFrom(
        address from,
        address to,
        uint256 value
    ) external returns (bool) {
        require(value <= _balanceOf[from],"not enough balance");
        require(value <= _allowance[from][msg.sender],"not enough allowance");
        _balanceOf[from] -= value;
        _balanceOf[to] += value;
        emit Transfer(from, to, value);
        return true;
    }
}
