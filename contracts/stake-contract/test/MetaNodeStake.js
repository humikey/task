const { expect } = require("chai");
const { ethers, upgrades } = require("hardhat");

describe("MetaNodeStake Contract", function () {
  let MetaNodeStakeFactory, metaNodeStake;
  let MetaNodeTokenFactory, metaNodeToken;
  let owner, admin, user1, user2;

  beforeEach(async function () {
    // 获取测试账户
    [owner, admin, user1, user2] = await ethers.getSigners();

    //     const MetaNodeTokenc = await ethers.getContractFactory('MetaNodeToken')
    // const metaNodeToken = await MetaNodeTokenc.deploy()

    // await metaNodeToken.waitForDeployment();
    // const metaNodeTokenAddress = await metaNodeToken.getAddress();

    // 部署 MetaNodeToken (ERC20)
    MetaNodeTokenFactory = await ethers.getContractFactory("MetaNodeToken");
    metaNodeToken = await MetaNodeTokenFactory.deploy();
    await metaNodeToken.waitForDeployment();
    const metaNodeTokenAddress = await metaNodeToken.getAddress();

    // 部署 MetaNodeStake 合约

    // 1. 获取合约工厂
    MetaNodeStakeFactory = await ethers.getContractFactory("MetaNodeStake");

    // 2. 设置初始化参数（根据你的initialize函数）
    // 例如:
    // IERC20 _MetaNode, uint256 _startBlock, uint256 _endBlock, uint256 _MetaNodePerBlock
    // 你需要替换下面的参数为实际的MetaNode代币地址和区块参数
    // const metaNodeTokenAddress = "0x5FbDB2315678afecb367f032d93F642f64180aa3"; // 替换为实际MetaNode代币地址
    const startBlock = 1; // 替换为实际起始区块
    const endBlock = 999999999999; // 替换为实际结束区块
    const metaNodePerBlock = ethers.parseUnits("1", 18); // 每区块奖励1个MetaNode（18位精度）

    // 3. 部署可升级代理合约
    metaNodeStake = await upgrades.deployProxy(
      MetaNodeStakeFactory,
      [metaNodeTokenAddress, startBlock, endBlock, metaNodePerBlock],
      { initializer: "initialize" }
    );

    //console.log("MetaNodeStake (proxy) deployed to:", stake);

    await metaNodeStake.waitForDeployment();

    // 设置角色
    //await metaNodeStake.grantRole(await metaNodeStake.ADMIN_ROLE(), admin.address);
  });

  it("should initialize the contract correctly", async function () {
    expect(await metaNodeStake.startBlock()).to.equal(1);
    expect(await metaNodeStake.endBlock()).to.equal(999999999999);
    expect(await metaNodeStake.MetaNodePerBlock()).to.equal(ethers.parseUnits("1", 18));
    expect(await metaNodeStake.MetaNode()).to.equal(await metaNodeToken.getAddress());
  });

  it("should allow admin to pause and unpause withdraw", async function () {
    await metaNodeStake.connect(owner).pauseWithdraw();
    expect(await metaNodeStake.withdrawPaused()).to.be.true;

    await metaNodeStake.connect(owner).unpauseWithdraw();
    expect(await metaNodeStake.withdrawPaused()).to.be.false;
  });

  it("should allow admin to add a new pool", async function () {
    await metaNodeStake.connect(owner).addPool(
      "0x0000000000000000000000000000000000000000", // staking token address
      100, // pool weight
      ethers.parseUnits("1", 18), // min deposit amount
      50, // unstake locked blocks
      false // withUpdate
    );

    const pool = await metaNodeStake.pool(0);
    expect(pool.stTokenAddress).to.equal("0x0000000000000000000000000000000000000000");
    expect(pool.poolWeight).to.equal(100);
    expect(pool.minDepositAmount).to.equal(ethers.parseUnits("1", 18));
    expect(pool.unstakeLockedBlocks).to.equal(50);
  });

  it("should allow users to deposit tokens", async function () {
    // 用户授权并存款
    await metaNodeToken.connect(owner).approve(metaNodeStake.address, ethers.parseUnits("10", 18));

    console.log("Pool added successfully", user1);
    // 添加资金池
    await metaNodeStake.connect(user1).addPool(
      user1.address, // staking token address
      100, // pool weight
      ethers.parseUnits("1", 18), // min deposit amount
      50, // unstake locked blocks
      false // withUpdate
    );

    await metaNodeStake.connect(user1).deposit(1, ethers.parseUnits("10", 18));

    const userInfo = await metaNodeStake.user(1, user1.address);
    expect(userInfo.stAmount).to.equal(ethers.parseUnits("10", 18));
  });

  it("should allow users to claim rewards", async function () {
    // 添加资金池
    await metaNodeStake.connect(user1).addPool(
      user1.address, // staking token address
      100, // pool weight
      ethers.parseUnits("1", 18), // min deposit amount
      50, // unstake locked blocks
      false // withUpdate
    );

    // 用户授权并存款
    await metaNodeStake.connect(user1).deposit(1, ethers.parseUnits("1", 18));

    // 快进区块
    await ethers.provider.send("evm_mine", [200]);

    // 用户领取奖励
    await metaNodeStake.connect(user1).claim(1);

    const userInfo = await metaNodeStake.user(1, user1.address);
    expect(userInfo.pendingMetaNode).to.equal(0);
  });

  it("should allow admin to add and update a staking pool", async function () {
    const ADMIN_ROLE = await metaNodeStake.ADMIN_ROLE();

    // 确保只有管理员可以操作
    await metaNodeStake.grantRole(ADMIN_ROLE, admin.address);

    // 添加新的质押池
    await metaNodeStake.connect(admin).addPool(
      await metaNodeToken.getAddress(), // staking token address
      200, // pool weight
      ethers.parseUnits("5", 18), // min deposit amount
      100, // unstake locked blocks
      false // withUpdate
    );

    // 验证质押池是否正确添加
    const pool = await metaNodeStake.pool(0);
    expect(pool.stTokenAddress).to.equal(await metaNodeToken.getAddress());
    expect(pool.poolWeight).to.equal(200);
    expect(pool.minDepositAmount).to.equal(ethers.parseUnits("5", 18));
    expect(pool.unstakeLockedBlocks).to.equal(100);

    // 更新质押池的配置
    await metaNodeStake.connect(admin).updatePool(0, ethers.parseUnits("10", 18), 200);

    // 验证质押池是否正确更新
    const updatedPool = await metaNodeStake.pool(0);
    expect(updatedPool.minDepositAmount).to.equal(ethers.parseUnits("10", 18));
    expect(updatedPool.unstakeLockedBlocks).to.equal(200);

    // 验证非管理员无法操作
    await expect(
      metaNodeStake.connect(user1).addPool(
        await metaNodeToken.getAddress(),
        100,
        ethers.parseUnits("1", 18),
        50,
        false
      )
    ).to.be.revertedWith(`AccessControl: account ${user1.address.toLowerCase()} is missing role ${ADMIN_ROLE}`);
  });
  it("should allow only upgrade role to upgrade the contract", async function () {
    const UPGRADE_ROLE = await metaNodeStake.UPGRADE_ROLE();

    // 确保只有持有升级角色的账户可以升级
    await metaNodeStake.grantRole(UPGRADE_ROLE, admin.address);

    // 模拟升级合约
    const NewMetaNodeStakeFactory = await ethers.getContractFactory("MetaNodeStake");
    const newMetaNodeStake = await upgrades.upgradeProxy(metaNodeStake.address, NewMetaNodeStakeFactory);

    // 验证升级后的合约地址是否一致
    expect(await newMetaNodeStake.address).to.equal(metaNodeStake.address);

    // 验证非升级角色无法升级
    await expect(
      upgrades.upgradeProxy(metaNodeStake.address, NewMetaNodeStakeFactory, { from: user1.address })
    ).to.be.revertedWith(`AccessControl: account ${user1.address.toLowerCase()} is missing role ${UPGRADE_ROLE}`);
  });

  it("should allow admin to pause and unpause staking, unstaking, and claiming", async function () {
    const ADMIN_ROLE = await metaNodeStake.ADMIN_ROLE();

    // 确保只有管理员可以操作
    await metaNodeStake.grantRole(ADMIN_ROLE, admin.address);

    // 暂停质押
    await metaNodeStake.connect(admin).pause();
    await expect(
      metaNodeStake.connect(user1).deposit(0, ethers.parseUnits("1", 18))
    ).to.be.revertedWith("Pausable: paused");

    // 恢复质押
    await metaNodeStake.connect(admin).unpause();
    await metaNodeToken.connect(user1).approve(metaNodeStake.address, ethers.parseUnits("10", 18));
    await metaNodeStake.connect(user1).deposit(0, ethers.parseUnits("1", 18));

    // 暂停解除质押
    await metaNodeStake.connect(admin).pauseWithdraw();
    await expect(
      metaNodeStake.connect(user1).unstake(0, ethers.parseUnits("1", 18))
    ).to.be.revertedWith("withdraw is paused");

    // 恢复解除质押
    await metaNodeStake.connect(admin).unpauseWithdraw();
    await metaNodeStake.connect(user1).unstake(0, ethers.parseUnits("1", 18));

    // 暂停领奖
    await metaNodeStake.connect(admin).pauseClaim();
    await expect(
      metaNodeStake.connect(user1).claim(0)
    ).to.be.revertedWith("claim is paused");

    // 恢复领奖
    await metaNodeStake.connect(admin).unpauseClaim();
    await metaNodeStake.connect(user1).claim(0);
  });
});