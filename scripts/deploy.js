const { ethers } = require("hardhat");

async function main() {
  const AssetEventRegistry = await ethers.getContractFactory("AssetEventRegistry");
  const contract = await AssetEventRegistry.deploy();

  await contract.waitForDeployment();

  console.log("AssetEventRegistry deployed to:", await contract.getAddress());
}

main()
  .then(() => process.exit(0))
  .catch((error) => {
    console.error(error);
    process.exit(1);
  });
