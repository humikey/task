## 原项目(未优化版本)

```solidity
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
```

## Gas 优化策略

### **优化策略 1：减少存储写入**

- 发现每次运算都写 `lastResult` 到存储（`SSTORE`），这是昂贵操作。
- 改为仅在需要时写入。

```solidity
// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

contract MathOptimized1 {
    function add(uint256 a, uint256 b) public pure returns (uint256) {
        return a + b;
    }

    function sub(uint256 a, uint256 b) public pure returns (uint256) {
        require(a >= b, "Underflow");
        return a - b;
    }
}
```

### **优化策略 2：移除 require，利用 unchecked**

- Solidity 0.8 默认会进行溢出检查，增加 Gas。
- 如果可以确保输入合法，可以使用 `unchecked` 来节省 Gas。

```solidity
// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

contract MathOptimized2 {
    function add(uint256 a, uint256 b) public pure returns (uint256) {
        unchecked {
            return a + b;
        }
    }

    function sub(uint256 a, uint256 b) public pure returns (uint256) {
        unchecked {
            return a - b;
        }
    }
}
```

测试结果

```powershell
[⠊] Compiling...
[⠘] Compiling 4 files with Solc 0.8.30
[⠃] Solc 0.8.30 finished in 670.65ms
Compiler run successful!

Ran 2 tests for test/MathOptimized1.t.sol:MathOptimized1Test
[PASS] testAdd() (gas: 9221)
Logs:
  Gas used add(): 6630

[PASS] testSub() (gas: 9268)
Logs:
  Gas used sub(): 6678

Suite result: ok. 2 passed; 0 failed; 0 skipped; finished in 292.90µs (109.10µs CPU time)

Ran 2 tests for test/Math.t.sol:MathTest
[PASS] testAdd() (gas: 31348)
Logs:
  Gas used add(): 28757

[PASS] testSub() (gas: 31417)
Logs:
  Gas used sub(): 28827

Suite result: ok. 2 passed; 0 failed; 0 skipped; finished in 407.70µs (134.00µs CPU time)

Ran 2 tests for test/MathOptimized2.t.sol:MathOptimized2Test
[PASS] testAdd() (gas: 9042)
Logs:
  Gas used add(): 6451

[PASS] testSub() (gas: 9063)
Logs:
  Gas used sub(): 6473

Suite result: ok. 2 passed; 0 failed; 0 skipped; finished in 344.40µs (208.30µs CPU time)

Ran 2 tests for test/Counter.t.sol:CounterTest
[PASS] testFuzz_SetNumber(uint256) (runs: 256, μ: 28978, ~: 29289)
[PASS] test_Increment() (gas: 28783)
Suite result: ok. 2 passed; 0 failed; 0 skipped; finished in 3.18ms (2.95ms CPU time)

Ran 4 test suites in 10.51ms (4.23ms CPU time): 8 tests passed, 0 failed, 0 skipped (8 total tests)
```



| 版本  | add() Gas | sub() Gas |
| ----- | --------- | --------- |
| 原始  | 28757     | 28827     |
| 优化1 | 6630      | 6678      |
| 优化2 | 6451      | 6473      |

##  分析报告

- **优化前：** 因为每次都写入 `lastResult`，一次运算包含昂贵的 `SSTORE`，导致 gas 接近 29k。
- **优化策略 1（移除存储写入）：** 减少了 `SSTORE`，Gas 降低约 60%（~20k Gas 节省）。
- **优化策略 2（unchecked 移除溢出检查）：** 进一步减少运算时的安全检查，Gas 又下降了。

结论：

1. **存储写入（SSTORE）是 Gas 消耗大头**，能不用就不要。
2. **unchecked 能进一步优化，但要在确保安全的前提下使用**。
3. 组合优化能让合约 Gas 从 **29k → 7k 左右**，节省了 ~60%。