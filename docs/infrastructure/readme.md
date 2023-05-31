## Overview

Some external dependencies or subsystems are needed during SP running, and these 
dependencies and subsystems are abstracted as interfaces for SP providers to 
customize their own implementation, for example: SP providers can use different 
storage vendors as the underlying storage of PieceStore, such as AWS S3, Ali OSS
etc. More resources will be consumed during uploading and downloading object payloads. 
In order to ensure reliable service quality, resource management is a necessary 
component. SP providers can customize their own resource managers. 

Of course, any components will have a default implementation if the SP provider is 
not customized. The default implementation of `PieceStore` supports local disk file,
AWS S3 and Min.IO three types. The default implementation of `SPDB` and `BSDB` is MySQL.
The default implementation of `Consensus` is greenfield full node, and offer an 
`ResourceManager` implementation that is Pre-reserved based on resource limits configuration.


### Main infrastructure includes:
* [PieceStore](./01-piece_store.md): PieceStore is the component to access piece store that 
  store the object payload data, any storage vendors can access SP as a piece store, as long 
  as the interface of PieceStore is implemented, the storage vendors agnostic to the user.
* [SPDB](./02-sp_db.md): SPDB is the component to records the SP metadata.
* [BSDB](./03-bs_db.md): BSDB is the component to records the greenfield metadata.
* [Consensus](./04-consensus.md): is the component to query greenfield consensus
  data, the consensus data can come from validator, full node, or other off-chain data service.
* [ResourceManager](../../core/rcmgr/README.md): ResourceManager is the interface to the resource
  management subsystem. The ResourceManager tracks and accounts for resource usage in the stack,
  from the internals to the application, and provides a mechanism to limit resource usage 
  according to a user configurable policy.
