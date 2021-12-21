data "twcc_project" "testProject" {
    name = "ENTxxxxxx"
    platform = "openstack-taichung-default-2"
}

resource "twcc_network" "network1" {
    cidr = "10.0.0.0/24"
    gateway = "10.0.0.254"
    name = "geminitestnet10"
    platform = data.twcc_project.testProject.platform
    project = data.twcc_project.testProject.id
    with_router = true
}

resource "twcc_firewall" "firewall1" {
    associate_networks = [twcc_network.network1.id]
    name = "geminitestwall1"
    platform = data.twcc_project.testProject.platform
    project = data.twcc_project.testProject.id
    rules = [
                twcc_firewall_rule.firewall_rule2.id,
                twcc_firewall_rule.firewall_rule1.id
            ]
}


resource "twcc_firewall_rule" "firewall_rule1" {
    name = "geminitestrule1"
    platform = data.twcc_project.testProject.platform
    project = data.twcc_project.testProject.id
    destination_ip_address = "10.0.0.0/24"
    destination_port = "22"
    action = "allow"
}

resource "twcc_firewall_rule" "firewall_rule2" {
    name = "geminitestrule2"
    platform = data.twcc_project.testProject.platform
    project = data.twcc_project.testProject.id
}
