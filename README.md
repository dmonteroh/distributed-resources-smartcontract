# Distributed Edge Resrouces - SmartContracts
Three distinct GOLANG Fabric Smart Contracts for Inventory Management, Edge Server Resource Collection and Latency Collection.

### Inventory Management
Keep different assets in the blockchain with their properties, e.g. Edge Servers & Robots, and functions associated with listing the different kinds of assets.

### Edge Server Resource Collection
Stores the data created by the [Distributed Resource Collector & Heartbeat](https://github.com/dmonteroh/distributed-resource-collector). Currently 2/3 configurations have been finished: Unique resources, Updatable resources. Resource Offloading is still a work in progress.

### Latency Collection
The latency collector Smart Contract stores the results of the Latency Measurement included in the [Distributed Resource Collector & Heartbeat](https://github.com/dmonteroh/distributed-resource-collector). It is also responsible for directly interacting with the **Inventory Management** Smart Contracts to get the necessary details and properties of the inventory assets.

# v0.1
Inventory Management, Edge Server Resource Collection and Latency Collection. Offloading data from the blockchain, data verirification functions and result pagination are still a Work In Progress.