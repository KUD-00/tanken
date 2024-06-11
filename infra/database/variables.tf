variable "identifier" {
  type        = string
  description = "The instance identifier for the RDS instance."
}

variable "instance_class" {
  type        = string
  description = "The compute and memory capacity of the RDS instance (e.g., 'db.t3.micro')."
}

variable "engine" {
  type        = string
  description = "The name of the database engine to be used for this instance."
}

variable "engine_version" {
  type        = string
  description = "The version number of the database engine to use."
}

variable "username" {
  type        = string
  description = "Username for the RDS database."
}

variable "password" {
  type        = string
  description = "Password for the RDS database."
}

variable "subnet_group_name" {
  type        = string
  description = "The name of the subnet group for the RDS instance."
}

variable "subnet_ids" {
  type        = list(string)
  description = "A list of subnet IDs to be used for the RDS subnet group."
}

variable "vpc_security_group_ids" {
  type        = list(string)
  description = "A list of security group IDs to associate with the RDS instance."
}

variable "allocated_storage" {
  type        = number
  description = "The allocated storage in gigabytes."
}

variable "storage_type" {
  type        = string
  description = "The type of storage to be used for the database instance."
}

variable "tags" {
  type        = map(string)
  description = "A mapping of tags to assign to the resource."
}
