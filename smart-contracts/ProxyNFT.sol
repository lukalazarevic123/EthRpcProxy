// SPDX-License-Identifier: MIT
pragma solidity 0.8.19;

import "@openzeppelin/contracts@4.9.3/token/ERC721/ERC721.sol";
import "@openzeppelin/contracts/utils/Counters.sol";

contract ProxyNFT is ERC721 {
    using Counters for Counters.Counter;
    Counters.Counter private _tokenIds;

    constructor() ERC721("ProxyNFT", "PXY") {}

    function mint(address _holder)
        public
        returns (uint256)
    {
        _tokenIds.increment();

        uint256 newTokenId= _tokenIds.current();
        _mint(_holder, newTokenId);

        return newTokenId;
    }
}