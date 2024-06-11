terraform {
  required_version = ">= 1.0.0, < 2.0.0"

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 4.0"
    }
    # kubernetes = {
    #   source  = "hashicorp/kubernetes"
    #   version = "~> 2.0"
    # }
  }
}

provider "aws" {
  # Tokyo
  region = "ap-northeast-1"
}

data "aws_vpc" "default" {
  default = true
}

# 有两个子网， elasticache 和 rds 都共用它们，可以保证容灾可用性
# elasticsearch 可能需要单独的子网，因为需要更大的ip地址空间？

resource "aws_subnet" "data_subnet_1" {
  vpc_id            = "data.aws_vpc.default.id"
  cidr_block        = "172.31.48.0/27"        # 32 IPs
  availability_zone = "ap-northeast-1a"

  tags = {
    Name = "data elasticache subnet 1"
  }
}

resource "aws_subnet" "data_subnet_2" {
  vpc_id            = "data.aws_vpc.default.id"
  cidr_block        = "172.31.48.32/27" # Remains 172.31.48.64 - 172.31.48.255 for scaling needs
  availability_zone = "ap-northeast-1b"

  tags = {
    Name = "data elasticache subnet 2"
  }
}

module "security_groups" {
  source                = "./modules/security_groups"
  vpc_id                = data.aws_vpc.default.id
  cache_ingress_cidr_blocks = [aws_vpc.main.cidr_block]
  rds_ingress_cidr_blocks   = [aws_vpc.main.cidr_block] #TODO: 实际上只有应用k8s集群才需要访问缓存服务器，这里暂时开放vpc下所有
}

module "elasticache" {
  source              = "./modules/elasticache"
  cluster_id          = "data_cache"
  engine              = "redis"
  node_type           = "cache.t2.micro"
  num_cache_nodes     = 1
  parameter_group_name = "default.redis7.1"
  engine_version      = "7.1"
  port                = 6379
  security_group_ids  = [module.security_groups.data_cache_security_group_id]
  subnet_group_name   = "data_cache_subnet_group"
  subnet_ids          = [aws_subnet.data_subnet_1.id, aws_subnet.data_subnet_2.id]
}

module "postgres_database" {
  source                = "./modules/postgres_database"
  identifier            = "data-postgres-db"
  instance_class        = "db.t3.micro"
  engine                = "postgres"
  engine_version        = "15.6"
  username              = "admin"
  password              = var.db_password  # Ensure this variable is defined in your variables.tf at root
  subnet_group_name     = "rds-subnet-group"
  subnet_ids            = [aws_subnet.data_subnet_1.id, aws_subnet.data_subnet_2.id]
  vpc_security_group_ids= [module.security_groups.rds_security_group_id]
  allocated_storage     = 20
  storage_type          = "gp2"
  tags                  = {
    Name = "Data-PostgresDB"
  }
}

/* resource "aws_elasticsearch_domain" "es" {
  domain_name           = "my-es-domain"
  elasticsearch_version = "7.10"

  cluster_config {
    instance_type = "t3.small.elasticsearch"
  }

  ebs_options {
    ebs_enabled = true
    volume_size = 10
  }

  vpc_options {
    subnet_ids = [
      aws_subnet.data_subnet_1.id,
      aws_subnet.data_subnet_2.id
    ]
    security_group_ids = [
      aws_security_group.data_security_group.id
    ]
  }

  access_policies = jsonencode({
    Version   = "2012-10-17"
    Statement = [
      {
        Action    = ["es:*"],
        Principal = {
          AWS = ["*"]
        },
        Effect    = "Allow",
        Resource  = "*"
      },
    ]
  })

  tags = {
    Domain = "MyElasticsearchDomain"
  }
} */

module "s3_buckets" {
  source             = "./modules/s3_buckets"
  bucket_name        = "go-tanken-picture-bucket"
  access_point_name  = "picture-bucket-access-point"
  vpc_id             = data.aws_vpc.default.id
  bucket_tags        = {
    Name        = "picture-bucket"
    Environment = "Dev"
  }
}
