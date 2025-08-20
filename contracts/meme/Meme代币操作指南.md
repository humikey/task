# SHIB 风格 Meme 代币合约——部署与运维操作指南

> 适用对象：负责代币发行、部署与日常运维的工程师与运营人员。
>
> 合约特性（前置）：交易税（买/卖/转账）、自动加池（auto-liquidity）、营销钱包、销毁、单笔额度与按日交易次数限制、白名单豁免、手动处理税池等。

------

## 目录

1. 背景与目标
2. 环境准备（Hardhat/Foundry 二选一）
3. 合约与仓库结构
4. 链与 Router 选择（务必自校验）
5. 部署流程（一步步）
6. 初始流动性与开盘流程
7. 运营期参数调优与日常操作
8. 交易限制与白名单管理
9. 自动加池与税池清算机制
10. 常见操作脚本（Hardhat）
11. Etherscan/区块浏览器交互（UI 步骤）
12. 监控与告警（事件/指标）
13. 风险与合规注意事项
14. 常见问题（FAQ）

------

## 1. 背景与目标

本指南围绕一份 SHIB 风格 ERC20 代币合约展开，覆盖从**部署、初始化、提供流动性、开盘、税率与限额运营、监控与风控**到**升级迁移**的完整周期操作，强调「安全优先」「可观测」「可回滚」。

------

## 2. 环境准备（Hardhat）

### 2.1 基础依赖

- Node.js ≥ 18.x，pnpm 或 npm
- 一个或多个链的 RPC（Alchemy/Infura/自建节点）
- 私钥（建议使用**测试钱包**演练；生产使用多签/硬件钱包）

### 2.2 选项 A：Hardhat

```bash
# 初始化
mkdir shibx && cd shibx
npm init -y
npm i -D hardhat typescript ts-node @types/node @nomicfoundation/hardhat-toolbox @openzeppelin/contracts
npx hardhat # 选择 TypeScript 项目
```

## 3. 合约与仓库结构

```
shibx/
├─ contracts/
│  └─ ShibStyleToken.sol    # 代币合约（前文示例）
├─ scripts/
│  ├─ deploy.ts             # 部署脚本
│  ├─ init.ts               # 初始化参数脚本（税率/份额/限额）
│  ├─ ops.ts                # 常用运营脚本（调整税、白名单、手动清算）
│  └─ lp.ts                 # 用户 LP 封装脚本（加/退流动性）
├─ test/                    # 建议添加 e2e 测试
├─ hardhat.config.ts
├─ .env                     # RPC/私钥/路由器等
└─ package.json
```

------

## 4. 链与 Router 选择（务必自校验）

合约示例基于 **UniswapV2-兼容 Router** 接口（`IUniswapV2Router02`）。不同公链/DEX 的 Router 地址可能不同且会更新：

- **务必以目标 DEX 官方文档/推特/Discord/区块浏览器验证最新 Router 地址**。
- 主网/测试网地址常见于：以太坊（Uniswap/Sushiswap）、BSC（Pancake）、Polygon、Arbitrum、Base、Avalanche、Optimism 等。
- 你也可以选择部署在测试网（Sepolia/BSC Testnet 等）完成演练。

> 在 `.env` 中以 `ROUTER=` 占位，并在部署时注入，杜绝硬编码。

------

## 5. 部署流程（一步步）

### 5.1 准备 `.env`

```
# .env （示例）
PRIVATE_KEY=0xabc...            # 部署私钥（测试钱包）
RPC_URL=https://...             # 目标链 RPC
ROUTER=0x...                    # 目标链 DEX Router（务必自校验）
MARKETING=0x...                 # 营销钱包（多签地址）
SUPPLY=1000000000000000000000000000  # 1e27 = 1B * 1e18
ETHERSCAN_KEY=...               # 可选：验证合约
```

### 5.2 Hardhat 配置 `hardhat.config.ts`

```ts
import "@nomicfoundation/hardhat-toolbox";
import * as dotenv from "dotenv"; dotenv.config();

const config = {
  solidity: {
    version: "0.8.24",
    settings: { optimizer: { enabled: true, runs: 200 } }
  },
  networks: {
    target: {
      url: process.env.RPC_URL || "",
      accounts: process.env.PRIVATE_KEY ? [process.env.PRIVATE_KEY] : [],
    }
  },
  etherscan: { apiKey: process.env.ETHERSCAN_KEY || "" }
};
export default config;
```

### 5.3 部署脚本 `scripts/deploy.ts`

```ts
import { ethers } from "hardhat";

async function main() {
  const supply = BigInt(process.env.SUPPLY!);
  const router = process.env.ROUTER!;
  const marketing = process.env.MARKETING!;

  const F = await ethers.getContractFactory("ShibStyleToken");
  const c = await F.deploy(supply, router, marketing);
  await c.waitForDeployment();

  const addr = await c.getAddress();
  console.log("Token deployed:", addr);
  console.log("Pair:", await c.pair());
  console.log("Router:", router);
}

main().catch((e)=>{ console.error(e); process.exit(1); });
```

### 5.4 执行部署

```bash
npx hardhat run scripts/deploy.ts --network target
```

输出：代币地址、交易对地址、路由器地址。

### 5.5 （可选）区块浏览器验证

```bash
npx hardhat verify --network target <TokenAddress> <SUPPLY> <ROUTER> <MARKETING>
```

> 若合约构造参数与源码匹配、编译器版本一致，即可通过验证。

------

## 6. 初始流动性与开盘流程

> 建议在**小范围私测**确认：加池、买卖、税与自动加池、限额与白名单逻辑正常后，再进行公开开盘。

### 6.1直接用 Router 加池

1. 在区块浏览器打开 Router 合约（`addLiquidityETH` 写入界面）。
2. `token = <你的代币地址>`；
3. 先在代币合约 `approve(router, tokenAmount)`；
4. 在 Router 写入：`amountTokenDesired`、`amountTokenMin`、`amountETHMin`、`to=你的地址`、`deadline`，并**附带 ETH**；
5. 成功后你将获得 LP 代币（发送到 `to`）。

### 6.2 开盘顺序

1. **参数预设**（见第 7 节）：较高限额限制 + 低中等税率；
2. **加池**（小额试运行 → 正式加池）；
3. **白名单设置**：运营钱包/做市脚本/营销钱包豁免税与限额；
4. **开启交易**：`enableTrading()`；
5. **逐步放宽限制**：根据成交与社区反馈动态调整 `maxTxAmount` 与 `dailyTradeLimit`；
6. **上线公告**：透明披露税率/分配/是否锁池/营销钱包。

------

## 7. 运营期参数调优与日常操作

### 7.1 税率与分配（营销/加池/销毁）

- `setTaxes(buy, sell, transfer)`（bp，1% = 100）：单侧上限 10%（合约内置）。
- `setSplits(marketing, liquidity, burn)`（总和 ≤ 10000）：控制税池清算后的分配比例。
- 建议：
  - 上线初期：买卖 2%~4%，转账 0%~1%；
  - 分配：营销 30%~50%，加池 30%~50%，销毁 0%~20%；
  - 随成交量与社区反馈动态微调。

### 7.2 自动加池与阈值

- `setSwapSettings(enabled, threshold)`：开启与阈值（以 token 计）。
- 建议阈值：总量的 0.001%~0.01%，避免频繁触发或长时间不触发。
- 手动清算：`manualProcess(amount)`，amount=0 表示清算全部税池。

### 7.3 营销钱包与资金流

- `setMarketingWallet(addr)`：建议设为**多签**地址。
- 自动清算后，营销份额会**换成 ETH** 转至该地址；请建立出入账台账与披露机制。

### 7.4 日志与事件

- 关键事件：`Transfer`、`TaxesUpdated`、`SplitsUpdated`、`SwapSettingsUpdated`、`TradingEnabled`、`LimitsUpdated`、`RouterPairUpdated`、`MarketingWalletUpdated`。
- 运营期需订阅并记录（见第 12 节）。

------

## 8. 交易限制与白名单管理

### 8.1 交易限制

- 单笔最大额：`setLimits(maxTxAmount, dailyTrades)`。
- 每日交易次数：`dailyTrades`（每地址，按 `block.timestamp/1days` 计日）。
- 关闭/放宽：动态下调限制或将关键地址加入白名单。

### 8.2 白名单

- `setExempt(account, limitExempt, taxExempt)`：
  - `limitExempt`：豁免交易额度/频率限制；
  - `taxExempt`：豁免税收。
- 典型白名单对象：做市机器人、桥接/金库地址、运营/空投分发地址。

### 8.3 交易开关

- `enableTrading()`：只可开启一次，开启前仅白名单可动。

------

## 9. 自动加池与税池清算机制（内部流程）

1. 正常交易收税：按买/卖/转账税率计入合约余额；
2. 触发条件：在**卖出路径**、`balanceOf(this) ≥ swapThreshold` 且 `swapEnabled` 时；
3. 清算流程：
   - 按 `splitBps` 拆分：燃烧 → 营销（换成 ETH）→ 加池（50% 换 ETH + 50% token 配对）；
   - 加池 LP 接收方：当前为 `owner`（可改为 `DEAD` 达到锁池效果，或发送至锁仓合约）。

> 设计权衡：自动加池有助于积累深度，但也会在触发块引入额外交易与可被感知的 MEV 风险。阈值与触发频率需权衡。

------

## 10. 常见操作脚本（Hardhat）

> 以下脚本以 `npx hardhat run scripts/xxx.ts --network target` 执行。

### 10.1 初始化参数 `scripts/init.ts`

```ts
import { ethers } from "hardhat";

async function main(){
  const token = await ethers.getContractAt("ShibStyleToken", process.env.TOKEN!);
  await (await token.setTaxes(300, 300, 100)).wait();   // 买3% 卖3% 转1%
  await (await token.setSplits(4000, 4000, 2000)).wait(); // 营销40/加池40/销毁20
  await (await token.setSwapSettings(true, BigInt(process.env.THRESHOLD!))).wait();
  await (await token.setLimits(BigInt(process.env.MAX_TX!), 20)).wait();
  console.log("init done");
}
main().catch(console.error);
```

### 10.2 白名单与开盘 `scripts/ops.ts`

```ts
import { ethers } from "hardhat";

async function main(){
  const token = await ethers.getContractAt("ShibStyleToken", process.env.TOKEN!);
  // 豁免示例：做市与营销
  await (await token.setExempt(process.env.MARKETING!, true, true)).wait();
  // 开盘
  await (await token.enableTrading()).wait();
  console.log("trading enabled");
}
main().catch(console.error);
```

### 10.3 手动清算与应急 `scripts/ops.ts`（扩展）

```ts
// 手动清算全部税池
await (await token.manualProcess(0)).wait();
// 应急转移合约内 ETH（例如四舍五入残余）
await (await token.rescueETH("0xYourMultisig", ethers.parseEther("0.5"))).wait();
```

------

## 11. Etherscan/区块浏览器交互（UI 步骤）

1. 打开代币合约页面 → **Contract → Write**：连接钱包；
2. 参数修改：调用 `setTaxes`/`setSplits`/`setSwapSettings`/`setLimits`/`setExempt`；
3. 开盘：`enableTrading`；
4. 手动清算：`manualProcess`；
5. 增/退 LP：优先在 Router 合约页面调用（`addLiquidityETH` / `removeLiquidityETHSupportingFeeOnTransferTokens`）。

------

## 12. 监控与告警（事件/指标）

建议自建告警（Telegram/Discord/Bark）：

- **事件**：`TaxesUpdated`、`SplitsUpdated`、`TradingEnabled`、`LimitsUpdated`、`RouterPairUpdated`；
- **阈值**：合约余额（税池）、营销钱包入账、`pair` 储备变化、LP 代币持有人分布；
- **价格/深度**：基于 DEX 子图或自抓 `getReserves`；
- **安全**：owner 变更、营销钱包变更、Router/Pair 变更即时告警。

------

## 13. 风险与合规注意事项

- **合约安全**：务必审计；上线前在测试网完整走通所有路径（买/卖/税/加池/限额/白名单）。
- **所有权治理**：生产环境使用**多签**/**时间锁**；敏感操作设定延时窗口与公开披露。
- **税与限制披露**：公开说明税率、阈值、用途、LP 是否锁定、营销钱包归属与支出。
- **市场风险**：自动加池与销毁并不能保证价格，仅改善深度与筹码结构；谨防过高税率导致流动性枯竭。
- **MEV/三明治**：自动清算会在卖出时段触发，注意阈值、分散触发与外部私有交易通道（如自带保护的路由）。

------

## 14. 常见问题（FAQ）

**Q1：交易提示 `trading disabled`？** 需先 `enableTrading()`，或将地址加入 `isExemptFromLimits`。

**Q2：`over maxTx`？** 调整 `maxTxAmount` 或给予 `limitExempt`。

**Q3：自动加池不触发？** 检查 `swapEnabled`、`swapThreshold`、合约余额≥阈值、且处于**卖出**路径。

**Q4：营销资金未到账？** 查看合约 `balance` 是否产生 ETH、营销钱包是否正确、`manualProcess(0)` 尝试手动清算。

**Q5：加池失败/滑点报错？** 提高 `amountTokenMin/amountETHMin` 容忍度或调整交易滑点；确认先 `approve`。

**Q6：每日限次如何计算？** 以 `block.timestamp/1 days` 为日界，UTC 计日，跨时区注意差异。

**Q7：如何锁池？** 将 LP 代币发送至 `DEAD` 或专用锁仓合约（第三方 Locker），并公开锁定期与交易哈希。

