// SPDX-License-Identifier: MIT
pragma solidity ^0.8.28;

contract AssetEventRegistry {

    enum EventType { CREATED, UPDATED, DELETED }

    struct AssetEvent {
        EventType eventType;
        string assetId;
        string assetCode;
        string assetName;
        string category;
        uint256 value;
        uint256 timestamp;
        address performedBy;
    }

    AssetEvent[] public events;

    event AssetEventRecorded(
        EventType eventType,
        string assetId,
        string assetCode,
        uint256 timestamp,
        address performedBy
    );

    function recordEvent(
        EventType eventType,
        string memory assetId,
        string memory assetCode,
        string memory assetName,
        string memory category,
        uint256 value
    ) public {
        events.push(
            AssetEvent({
                eventType: eventType,
                assetId: assetId,
                assetCode: assetCode,
                assetName: assetName,
                category: category,
                value: value,
                timestamp: block.timestamp,
                performedBy: msg.sender
            })
        );

        emit AssetEventRecorded(
            eventType,
            assetId,
            assetCode,
            block.timestamp,
            msg.sender
        );
    }

    function getEventsCount() public view returns (uint256) {
        return events.length;
    }
}
