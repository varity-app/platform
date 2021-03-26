module "network" {
    source = "./network"

    security_group_name = var.security_group_name
}