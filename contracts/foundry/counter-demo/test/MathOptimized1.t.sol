// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "forge-std/Test.sol";
import "../src/MathOptimized1.sol";

contract MathOptimized1Test is Test {
    MathOptimized1 math;

    function setUp() public {
        math = new MathOptimized1();
    }

    function testAdd() public {
        uint256 gasBefore = gasleft();
        uint256 result = math.add(10, 5);
        uint256 gasAfter = gasleft();
        emit log_named_uint("Gas used add()", gasBefore - gasAfter);

        assertEq(result, 15);
    }

    function testSub() public {
        uint256 gasBefore = gasleft();
        uint256 result = math.sub(10, 5);
        uint256 gasAfter = gasleft();
        emit log_named_uint("Gas used sub()", gasBefore - gasAfter);

        assertEq(result, 5);
    }
}
