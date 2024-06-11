variable "cluster_id" {
  type        = string
  description = "The identifier for the Elasticache cluster."
}

variable "engine" {
  type        = string
  description = "The engine type of the Elasticache cluster (e.g., 'redis')."
}

variable "node_type" {
  type        = string
  description = "The compute and memory capacity of the nodes in the node group."
}

variable "num_cache_nodes" {
  type        = number
  description = "The number of cache nodes that the cache cluster should have."
}

variable "parameter_group_name" {
  type        = string
  description = "The name of the parameter group to associate with this cache cluster."
}

variable "engine_version" {
  type        = string
  description = "The version number of the cache engine to use."
}

variable "port" {
  type        = number
  description = "The port number on which each of the cache nodes will accept connections."
}

variable "security_group_ids" {
  type        = list(string)
  description = "A list of security group IDs to associate with this cluster."
}

variable "subnet_group_name" {
  type        = string
  description = "The name of the subnet group associated with the Elasticache cluster."
}

variable "subnet_ids" {
  type        = list(string)
  description = "A list of subnet IDs to be used for the Elasticache subnet group."
}
