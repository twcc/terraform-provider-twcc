data "twcc_project" "testProject" {
    name = "ENTxxxxxx"
    platform = "openstack-taichung-default-2"
}

resource "twcc_network" "network1" {
    cidr = "10.0.0.0/24"
    gateway = "10.0.0.254"
    name = "geminitestnet11"
    platform = data.twcc_project.testProject.platform
    project = data.twcc_project.testProject.id
    with_router = true
}

resource "twcc_loadbalancer" "lb1" {
    listeners {
        name = "geminilbl1"
        pool_name = "geminilbp1"
        protocol = "HTTP"
        protocol_port = 80
    }
    name = "geminitestlb5"
    platform = data.twcc_project.testProject.platform
    pools {
        method = "ROUND_ROBIN"
        name = "geminilbp1"
        protocol = "HTTP"
        members {
            ip = "10.0.0.1"
        }

        members {
            ip = "10.0.0.2"
            port = 81
        }

        members {
            ip = "10.0.0.3"
            port = 82
            weight = 2
        }
    }
    private_net = twcc_network.network1.id
}
