---
marp: true
style: @import url('../tailwind.utilities.css');
size: 16:9
---
<!-- _class: lead -->
# Using HashiCorp Raft to distribute data between nodes when a database server is not an option

---

# Motivation
- We are a Zero Trust Network Access (ZTNA) company
- A ZTNA architecture is comprised of a control-plane and a data-plane
- What makes us unique is the fact that our control-plane is in the customer's on-prem and not in the cloud

---
# Problems we face
- Customer may have multiple networks (which means we need multiple connectors)
- Customer may want HA (which means we need multiple connectors per network)
- Customer networks are not always inter-connected (which means we cannot just use a shared database)

---

# Design Patterns
- stateless tokens 
  - requires a shared-secret between the nodes (or even PKI)
  - requires a client-side store (such as cookies or javascript)
  - each node can read/write state securely and easily
  
  
- "sticky" transactions
  - requires a method of "calling" a node directly
  - best for managing a resource that is only available on one node
