# beacon
The beacon serves to guide components to their partners.

## Base Functionality
The beacon interacts with the storage system to maintain what components are
registered to the system. When a request comes for a given component type,
the beacon sends the connection information back.

## Current
The storage system will initially be a docker volume connected to the beacon
container on the same machine. There will be a single beacon and it is assumed
that the churn rate is 0.

The next evolution is the external storage system to allow beacon distribution
and handling nebula component churn. 
