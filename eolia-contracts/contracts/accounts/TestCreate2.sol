// SPDX-License-Identifier: MIT
pragma solidity ^0.8.28;

import "@openzeppelin/contracts/utils/Create2.sol";
import "@openzeppelin/contracts/proxy/ERC1967/ERC1967Proxy.sol";

/**
 * Test contract to verify Create2 computation
 * 测试合约来验证Create2计算
 */
contract TestCreate2 {
    
    /**
     * Test the exact same logic as SimpleAccountFactory
     * 测试与SimpleAccountFactory完全相同的逻辑
     */
    function testGetAddress(address owner, uint256 salt) public view returns (address) {
        return Create2.computeAddress(bytes32(salt), keccak256(abi.encodePacked(
                type(ERC1967Proxy).creationCode,
                abi.encode(
                    address(0x1234567890123456789012345678901234567890), // dummy address
                    abi.encodeCall(this.dummyFunction, (owner))
                )
            )));
    }
    
    /**
     * Alternative implementation using manual concatenation
     * 使用手动连接的替代实现
     */
    function testGetAddressManual(address owner, uint256 salt) public view returns (address) {
        bytes memory creationCode = type(ERC1967Proxy).creationCode;
        bytes memory initData = abi.encode(
            address(0x1234567890123456789012345678901234567890), // dummy address
            abi.encodeCall(this.dummyFunction, (owner))
        );
        
        // Manual concatenation
        bytes memory combined = new bytes(creationCode.length + initData.length);
        for (uint i = 0; i < creationCode.length; i++) {
            combined[i] = creationCode[i];
        }
        for (uint i = 0; i < initData.length; i++) {
            combined[creationCode.length + i] = initData[i];
        }
        
        return Create2.computeAddress(bytes32(salt), keccak256(combined));
    }
    
    /**
     * Dummy function for testing
     * 用于测试的虚拟函数
     */
    function dummyFunction(address owner) external pure returns (bool) {
        return owner != address(0);
    }
    
    /**
     * Test with simple data
     * 使用简单数据测试
     */
    function testSimpleCreate2(uint256 salt) public view returns (address) {
        bytes memory data = abi.encodePacked("Hello", "World");
        return Create2.computeAddress(bytes32(salt), keccak256(data));
    }
    
    /**
     * Compute the same CREATE2 address as the factory using a real implementation and initializer
     * 使用真实的实现合约地址与初始化数据，计算与工厂相同的 CREATE2 地址
     *
     * Why this function?
     * - The earlier test function used a dummy implementation and a dummy initializer, so the hash differed
     * - This function mirrors the factory logic exactly: ERC1967Proxy + (implementation, abi.encodeWithSignature("initialize(address)", owner))
     *
     * 为什么需要这个函数？
     * - 之前的测试函数使用了假的实现地址和假的初始化函数，导致哈希不同
     * - 该函数严格镜像工厂逻辑：ERC1967Proxy + (实现地址, abi.encodeWithSignature("initialize(address)", owner))
     */
    function testGetAddressWithImpl(address implementation, address owner, uint256 salt) public view returns (address) {
        // 1) Build the constructor calldata for ERC1967Proxy
        // 1) 构造 ERC1967Proxy 的构造参数数据：实现地址 + 初始化数据
        bytes memory initData = abi.encode(
            implementation,
            // Use the same initializer signature as SimpleAccountFactory
            // 与 SimpleAccountFactory 中相同的初始化函数签名保持一致
            abi.encodeWithSignature("initialize(address)", owner)
        );

        // 2) Compute the CREATE2 address using the same bytecode hash as the factory
        // 2) 使用与工厂一致的字节码哈希计算 CREATE2 地址
        return Create2.computeAddress(
            bytes32(salt),
            keccak256(abi.encodePacked(
                type(ERC1967Proxy).creationCode,
                initData
            ))
        );
    }
}
