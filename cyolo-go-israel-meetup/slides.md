---
marp: true
style: @import url('../tailwind.utilities.css');
size: 16:9
---
<style>
img[alt~="center"] {
  display: block;
  margin: 0 auto;
}
</style>
# Using HashiCorp Raft to distribute data between nodes when a database server is not an option
![center](https://github.com/avivklas/talks/assets/6282578/762747e0-46b4-4d38-b006-c857ebfa6db2)

---
# Who are these guys and why can't they use a database server
- We are a Zero Trust Network Access (ZTNA) company
- Our mission is to isolate our customers' private networks and provide network access to it by identifying objects only

---
# Total isolation - the only solution
- Requires no inbound connection to your subnets
- Can operate in air-gapped networks
- Add no latency

  ![center](https://github.com/avivklas/talks/assets/6282578/3e97de60-aceb-43de-bf9e-4abebb1e1d88)

---
# Problems we face
- Customer may have multiple networks (which means we need multiple connectors)
- Customer may want HA (which means we need multiple connectors per network)
- Customer networks are not always inter-connected (which means we cannot just use a shared database)

---
# Session Management
- problem:
  - users need to be authenticated and have their session state needs to be maintained across nodes

- solution:
  - encrypted state tokens allow each node to read and modify the session state securely
  - encryption keys need to be pre-shared among the nodes (we can also leverage PKI)
  - encrypted tokens are stored within the client side (we use cookies) and are tamper-proof

![stateless cookie](cookie.png)

---
# Rendezvous Points
- problem:
  - some transactions require interaction between multiple actors
  - examples: MFA links, workflows and things that start in a browser and end in a different client
- solution:
  - sticky routing
  - host identifiers may be used as keys for sticky routing ("app-1.cyolo.io" or "app-2.cyolo.io" vs "app.cyolo.io")

---
## Raft-based data replication
- Highly consistent, highly available
- Can be used in our network infrastructure
- Allows us to use embedded, disk or in-mem storage

<br />
<br />
<br />
<br />

![center](img_1.png)

---

<!-- _class: lead -->
# Increment the counter

![center](img_4.png)