# Groups and Partitions

The HMS data model supports arbitrary grouping of nodes into component_groups.  These component groups support unique names and extended descriptions.  In addition, labels and annotations are supported as extended attributes, but not directly queriable.

Groups must be created with one of three types: Shared, Exclusive, and Partition.  That signals to others how they are intended to be used, but nothing is enforced differently at the database level.

## Shared

Shared groups are simply loose collections of nodes associated with a name.  The name is not guaranteed to be unique, but the id is.  A component can be a part of many shared groups and groups may overlap.

## Exclusive

Exclusive groups are intended for non-overlapping groups.  When adding or updating components to a group, the system needs to ensure that the component doesn't appear in two groups with the same `exclusive_type_identifier` at the same time.  

An example of this could be "availability-zone" as the exclusive_type_identifier and "us-east" and "us-west" as possible values.  A component can be associate with zero or one availability-zones, but never more than that.

## Partition

A partition is a specific kind of exclusive group.  It warrants an explicit type here so that the partitioning model we rely on for security is reflected directly in the data model.