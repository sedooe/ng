# ng

**ng** is a kubectl plugin that introduces **namespace group** concept for managing access controls of different teams/working groups within a single Kubernetes cluster.

It marks specified namespace with a **namespace group** such as dev, frontend, backend, production and provides easy to use commands for managing RBAC.

Every namespace can belong to at most 1 namespace group.

### Usage

```
USAGE:
  ng add <NAME>                            : add namespace to ng <NAME>
  ng remove                                : remove namespace from ng
  ng grant <USER> -g <NG> -r <ROLE>        : grant <USER> role <ROLE> for ng <NG>
  ng revoke <USER> -g <NG> -r <ROLE>       : revoke <USER>'s role <ROLE> for ng <NG>
  ng show [<NAME>]                         : show ng with its associated namespaces
```
