// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import {IRouterClient} from "@chainlink/contracts-ccip/src/v0.8/ccip/interfaces/IRouterClient.sol";
import {OwnerIsCreator} from "@chainlink/contracts-ccip/src/v0.8/shared/access/OwnerIsCreator.sol";
import {Client} from "@chainlink/contracts-ccip/src/v0.8/ccip/libraries/Client.sol";

contract AuctionRemote is OwnerIsCreator {
    IRouterClient public router;

    event MessageSent(bytes32 messageId, uint64 destinationChainSelector, address receiver, bytes data);

    constructor(address _router) {
        router = IRouterClient(_router);
    }


    function sendMessage(
        uint64 destinationChainSelector,
        address receiver,
        uint256 auctionId,
        address bidder,
        uint256 bidAmount
    ) external payable onlyOwner returns (bytes32) {
        bytes memory messageData = abi.encode(auctionId, bidder, bidAmount);
        // 组装跨链消息
        Client.EVM2AnyMessage memory message = Client.EVM2AnyMessage({
            receiver: abi.encode(receiver), // 接收方地址需要 abi.encode
            data: messageData,              // 跨链携带的数据
            tokenAmounts: new Client.EVMTokenAmount[](0) , // 不附带 token
            extraArgs: Client._argsToBytes(
                Client.EVMExtraArgsV1({gasLimit: 200_000})
            ),
            feeToken: address(0) // 使用原生代币支付费用
        });

        // 查询跨链费用
        uint256 fee = router.getFee(destinationChainSelector, message);
        require(msg.value >= fee, "Insufficient fee");

        // 发送消息
        bytes32 messageId = router.ccipSend{value: fee}(destinationChainSelector, message);

        emit MessageSent(messageId, destinationChainSelector, receiver, messageData);
        return messageId;
    }
}
