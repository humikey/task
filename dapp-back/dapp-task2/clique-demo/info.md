mkdir D:\chains\clique-demo\node1
mkdir D:\chains\clique-demo\node2

# 节点1账户
geth account new --datadir .\node1
输入密码：905000080hukui
0x64b4c81299891Ef660dF12d057834D4a6e7B9495

# 节点2账户
geth account new --datadir .\node2
输入密码：905000080hukui
0x45d84c099d9Fe6F58Ac4763635fca80dDE290D8D

# 生成 extraData（Clique 必需字段）
0x000000000000000000000000000000000000000000000000000000000000000064b4c81299891Ef660dF12d057834D4a6e7B949545d84c099d9Fe6F58Ac4763635fca80dDE290D8D0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000


# 编写genesis.json

# 初始化两个节点
geth --datadir .\node1 init .\genesis.json
geth --datadir .\node2 init .\genesis.json

# 启动第两个节点
geth --datadir .\node1 `
  --networkid 1337 `
--port 30311 `
  --http --http.addr 127.0.0.1 --http.port 8545 --http.api "admin,eth,net,web3,personal,txpool" `
--unlock 0x64b4c81299891Ef660dF12d057834D4a6e7B9495 `
  --password .\node1\passwd.txt `
--mine `
--allow-insecure-unlock

geth --datadir .\node1 `
  --networkid 1337 `
--port 30312 `
  --http --http.addr 127.0.0.1 --http.port 9545 --http.api "admin,eth,net,web3,personal,txpool" `
--unlock 0x45d84c099d9Fe6F58Ac4763635fca80dDE290D8D `
  --password .\node2\passwd.txt `
--mine `
--allow-insecure-unlock

# 节点控制台

admin.peers
eth.blockNumber