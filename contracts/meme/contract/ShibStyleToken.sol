// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "./Ownable.sol";
import "./ReentrancyGuard.sol";
import "./IUniswapV2Factory.sol";
import "./IUniswapV2Router02.sol";
//import "./IUniswapV2Pair.sol";

contract ShibStyleToken is Ownable, ReentrancyGuard {
    // ---------------- ERC20 基本 ----------------
    string  public name = "ShibStyle Token";
    string  public symbol = "SHIBX";
    uint8   public constant decimals = 18;
    uint256 public immutable totalSupply;

    mapping(address => uint256) public balanceOf;
    mapping(address => mapping(address => uint256)) public allowance;

    event Transfer(address indexed from, address indexed to, uint256 value);
    event Approval(address indexed owner_, address indexed spender, uint256 value);

    // ---------------- DEX/LP 集成 ----------------
    IUniswapV2Router02 public router;
    address public pair;         // token-WETH pair
    address public immutable WETH;

    // ---------------- 税收/分配 ----------------
    struct TaxBps { uint16 buy; uint16 sell; uint16 transfer_; } // 单位：bp（1% = 100）
    TaxBps public taxBps = TaxBps({buy: 300, sell: 300, transfer_: 100}); // 默认：买3%，卖3%，转1%

    struct SplitBps { uint16 marketing; uint16 liquidity; uint16 burn; }   // 三者之和 <= 10000
    SplitBps public splitBps = SplitBps({marketing: 4000, liquidity: 4000, burn: 2000}); // 默认：40%营销，40%加池，20%销毁

    uint16 public constant MAX_TAX_PER_SIDE = 1000; // 安全上限：每类税率<=10%
    address public marketingWallet;                 // 接收营销资金（ETH）
    address public constant DEAD = address(0xdead);  // 锁仓地址

    // 税收累积与自动加池
    uint256 public swapThreshold;    // 触发阈值（以 token 计）
    bool    public swapEnabled = true;
    bool    private _inSwap;

    // ---------------- 交易限制 ----------------
    bool    public tradingEnabled = false;
    uint256 public tradingStartBlock;

    uint256 public maxTxAmount;      // 单笔最大交易额
    uint16  public dailyTradeLimit;  // 每地址每日最多交易次数（买/卖/转账均计数）
    mapping(address => bool) public isExemptFromLimits; // 白名单（税/限额均豁免，可单独扩展）
    mapping(address => bool) public isExemptFromTax;

    struct DailyCounter { uint64 day; uint32 count; }
    mapping(address => DailyCounter) private _daily;

    // ---------------- 事件 ----------------
    event TaxesUpdated(uint16 buy, uint16 sell, uint16 transfer_);
    event SplitsUpdated(uint16 marketing, uint16 liquidity, uint16 burn);
    event SwapSettingsUpdated(bool enabled, uint256 threshold);
    event TradingEnabled(uint256 startBlock);
    event LimitsUpdated(uint256 maxTxAmount, uint16 dailyTrades);
    event RouterPairUpdated(address router, address pair);
    event MarketingWalletUpdated(address wallet);

    // ---------------- 构造函数 ----------------
    constructor(
        uint256 _supply,
        address _router,
        address _marketingWallet
    ) {
        require(_router != address(0) && _marketingWallet != address(0), "zero addr");

        totalSupply = _supply;
        balanceOf[msg.sender] = _supply;
        emit Transfer(address(0), msg.sender, _supply);

        router = IUniswapV2Router02(_router);
        WETH = router.WETH();
        marketingWallet = _marketingWallet;

        address factory = router.factory();
        pair = IUniswapV2Factory(factory).createPair(address(this), WETH);

        // 初始参数
        swapThreshold = _supply / 100000; // 0.001% 触发
        maxTxAmount  = _supply / 100;     // 单笔 1%
        dailyTradeLimit = 20;             // 每日 20 笔

        // 项目方与合约豁免
        isExemptFromLimits[msg.sender] = true;
        isExemptFromLimits[address(this)] = true;
        isExemptFromLimits[marketingWallet] = true;
        isExemptFromTax[msg.sender] = true;
        isExemptFromTax[address(this)] = true;
        isExemptFromTax[marketingWallet] = true;
    }

    // ---------------- ERC20 标准 ----------------
    function approve(address spender, uint256 amount) external returns (bool){
        // 允许 spender 账户从 msg.sender 账户中转移 amount 个代币
        allowance[msg.sender][spender] = amount;
        emit Approval(msg.sender, spender, amount);
        return true;
    }
    
    // 减去已经授权的额度
    function _spendAllowance(address from, address spender, uint256 amount) internal {
        uint256 cur = allowance[from][spender];
        if (cur != type(uint256).max){
            require(cur >= amount, "insufficient allowance");
            // 减少授权额度
            unchecked { allowance[from][spender] = cur - amount; }
        }
    }

    // 转账
    function transfer(address to, uint256 amount) external returns (bool){
        _transfer(msg.sender, to, amount);
        return true;
    }

    // 转账
    function transferFrom(address from, address to, uint256 amount) external returns (bool){
        // 检查授权额度 扣除授权额度
        _spendAllowance(from, msg.sender, amount);
        // 执行转账逻辑
        _transfer(from, to, amount);
        return true;
    }

    // ---------------- 内部转账逻辑（含税/限额/加池） ----------------
    function _transfer(address from, address to, uint256 amount) internal {
        require(from != address(0) && to != address(0), "zero addr");
        require(amount > 0, "zero amount");

        // 交易是否开放（白名单/流动性添加等可豁免）
        if (!tradingEnabled) {
            require(isExemptFromLimits[from] || isExemptFromLimits[to], "trading disabled");
        }

        // 基本余额检查
        uint256 fromBal = balanceOf[from];
        require(fromBal >= amount, "insufficient balance");

        // 单笔限额
        if (!isExemptFromLimits[from] && !isExemptFromLimits[to]) {
            require(amount <= maxTxAmount, "over maxTx");
        }

        // 每日交易次数限制
        if (!isExemptFromLimits[from] && !isExemptFromLimits[to] && dailyTradeLimit > 0) {
            uint64 day = uint64(block.timestamp / 1 days);
            DailyCounter memory dc = _daily[msg.sender];
            if (dc.day != day) { dc.day = day; dc.count = 0; }
            require(dc.count < dailyTradeLimit, "daily trades limit");
            unchecked { dc.count += 1; }
            _daily[msg.sender] = dc;
        }

        // 自动加池（仅在卖出前触发，减少干扰）
        if (
            swapEnabled &&
            !_inSwap &&
            to == pair &&                 // 卖出方向
            !isExemptFromLimits[from] && // 避免在初期/白名单时触发
            balanceOf[address(this)] >= swapThreshold
        ) {
            _inSwap = true;
            _processFeesAndLiquidity(balanceOf[address(this)]);
            _inSwap = false;
        }

        // 计算税费
        uint256 fee;
        if (!isExemptFromTax[from] && !isExemptFromTax[to]) {
            bool isBuy  = from == pair;
            bool isSell = to == pair;
            uint16 rate = isBuy ? taxBps.buy : (isSell ? taxBps.sell : taxBps.transfer_);
            if (rate > 0) {
                fee = (amount * rate) / 10000;
            }
        }

        // 转账 & 进税池
        unchecked {
            balanceOf[from] = fromBal - amount;
            uint256 receiveAmt = amount - fee;
            balanceOf[to] += receiveAmt;
            emit Transfer(from, to, receiveAmt);
        }
        if (fee > 0) {
            balanceOf[address(this)] += fee;
            emit Transfer(from, address(this), fee);
        }
    }

    // 处理累积税费：按 splitBps 分为燃烧/营销/加池
    function _processFeesAndLiquidity(uint256 amount) internal nonReentrant {
        if (amount == 0) return;
        uint16 tot = splitBps.marketing + splitBps.liquidity + splitBps.burn;
        if (tot == 0) return;

        uint256 toBurn = (amount * splitBps.burn) / tot;
        if (toBurn > 0) {
            unchecked { balanceOf[address(this)] -= toBurn; }
            balanceOf[DEAD] += toBurn;
            emit Transfer(address(this), DEAD, toBurn);
        }

        uint256 remaining = balanceOf[address(this)];
        if (remaining == 0) return;

        // 按比例拆分：营销（换成 ETH 转给金库） & 流动性（50%换ETH + 50%留在合约配对）
        uint256 forMarketing = (remaining * splitBps.marketing) / (tot - splitBps.burn);
        uint256 forLiquidity = remaining - forMarketing;

        // 营销：换 ETH -> 转金库
        if (forMarketing > 0) {
            _swapTokensForETH(forMarketing, address(this));
            uint256 ethBal = address(this).balance;
            if (ethBal > 0) {
                (bool ok,) = payable(marketingWallet).call{value: ethBal}("");
                require(ok, "marketing transfer fail");
            }
        }

        // 加池：一半换 ETH，一半留 token
        if (forLiquidity > 1) {
            uint256 half = forLiquidity / 2;
            uint256 otherHalf = forLiquidity - half;

            _swapTokensForETH(half, address(this));
            uint256 ethAmt = address(this).balance;
            if (ethAmt > 0 && otherHalf > 0) {
                _addLiquidity(otherHalf, ethAmt, owner); // LP 发给 owner（可改为 DEAD 进行永久锁池）
            }
        }
    }

    // 交换 Token -> ETH
    function _swapTokensForETH(uint256 tokenAmount, address to) internal {
        if (tokenAmount == 0) return;
        address;
        path[0] = address(this);
        path[1] = WETH;

        _approve(address(this), address(router), tokenAmount);
        router.swapExactTokensForETHSupportingFeeOnTransferTokens(
            tokenAmount, 0, path, to, block.timestamp
        );
    }

    // 加池（Token + ETH）
    // to 就是添加流动性的接收者
    function _addLiquidity(uint256 tokenAmount, uint256 ethAmount, address to) internal {
        _approve(address(this), address(router), tokenAmount);
        router.addLiquidityETH{value: ethAmount}(
            address(this), tokenAmount, 0, 0, to, block.timestamp
        );
    }

    // ---------------- 用户可调用：添加/移除流动性（封装） ----------------
    // 注意：用户也可直接使用 Router。本函数仅提供便捷封装。
    function addLiquidityETHForUser(uint256 tokenAmount, uint256 minToken, uint256 minETH, uint256 deadline)
        external payable nonReentrant returns (uint amountToken, uint amountETH, uint liquidity)
    {
        require(tradingEnabled, "trading disabled");
        _transfer(msg.sender, address(this), tokenAmount);
        _approve(address(this), address(router), tokenAmount);
        (amountToken, amountETH, liquidity) = router.addLiquidityETH{value: msg.value}(
            address(this), tokenAmount, minToken, minETH, msg.sender, deadline
        );
        // 退回未用尽的 token（若有）
        uint256 remain = allowance[address(this)][address(router)] > 0 ? 0 : 0; // 占位：Router已消费完授权，无需退回
    }

    function removeLiquidityETHForUser(uint256 liquidity, uint256 minToken, uint256 minETH, uint256 deadline)
        external nonReentrant returns (uint amountETH)
    {
        // 用户需提前把 LP 代币授权给 Router（合约不保管用户 LP）
        amountETH = router.removeLiquidityETHSupportingFeeOnTransferTokens(
            address(this), liquidity, minToken, minETH, msg.sender, deadline
        );
    }

    // ---------------- 管理/治理（建议交给多签或 Timelock） ----------------
    function enableTrading() external onlyOwner {
        require(!tradingEnabled, "already enabled");
        tradingEnabled = true;
        tradingStartBlock = block.number;
        emit TradingEnabled(tradingStartBlock);
    }

    function setRouter(address newRouter) external onlyOwner {
        require(newRouter != address(0), "zero");
        router = IUniswapV2Router02(newRouter);
        emit RouterPairUpdated(newRouter, pair);
    }

    function setPair(address newPair) external onlyOwner {
        require(newPair != address(0), "zero");
        pair = newPair;
        emit RouterPairUpdated(address(router), newPair);
    }

    function setMarketingWallet(address w) external onlyOwner {
        require(w != address(0), "zero");
        marketingWallet = w;
        emit MarketingWalletUpdated(w);
    }

    function setTaxes(uint16 buy, uint16 sell, uint16 transfer_) external onlyOwner {
        require(buy <= MAX_TAX_PER_SIDE && sell <= MAX_TAX_PER_SIDE && transfer_ <= MAX_TAX_PER_SIDE, "tax too high");
        taxBps = TaxBps({buy: buy, sell: sell, transfer_: transfer_});
        emit TaxesUpdated(buy, sell, transfer_);
    }

    function setSplits(uint16 marketing, uint16 liquidity, uint16 burn) external onlyOwner {
        require(marketing + liquidity + burn <= 10000, "sum>100%");
        splitBps = SplitBps({marketing: marketing, liquidity: liquidity, burn: burn});
        emit SplitsUpdated(marketing, liquidity, burn);
    }

    function setSwapSettings(bool enabled, uint256 threshold) external onlyOwner {
        swapEnabled = enabled;
        swapThreshold = threshold;
        emit SwapSettingsUpdated(enabled, threshold);
    }

    function setLimits(uint256 _maxTxAmount, uint16 _dailyTrades) external onlyOwner {
        require(_maxTxAmount > 0, "zero");
        maxTxAmount = _maxTxAmount;
        dailyTradeLimit = _dailyTrades;
        emit LimitsUpdated(_maxTxAmount, _dailyTrades);
    }

    function setExempt(address account, bool limitExempt, bool taxExempt) external onlyOwner {
        isExemptFromLimits[account] = limitExempt;
        isExemptFromTax[account] = taxExempt;
    }

    // 手动触发处理税费（用于手动清算/回收）
    function manualProcess(uint256 tokenAmount) external onlyOwner {
        require(!_inSwap, "busy");
        _inSwap = true;
        _processFeesAndLiquidity(tokenAmount == 0 ? balanceOf[address(this)] : tokenAmount);
        _inSwap = false;
    }

    // 紧急取回卡在合约的 ETH（例如因四舍五入）
    function rescueETH(address to, uint256 amount) external onlyOwner {
        (bool ok,) = payable(to).call{value: amount}("");
        require(ok, "rescue fail");
    }

    // ---------------- 内部工具 ----------------
    function _approve(address owner_, address spender, uint256 amount) internal {
        allowance[owner_][spender] = amount;
        emit Approval(owner_, spender, amount);
    }

    // 接收 Router swap 的 ETH
    receive() external payable {}
}