data "twcc_project" "testProject" {
    name = "ENTxxxxxx"
    platform = "openstack-taichung-default-2"
}

resource "twcc_network" "network1" {
    cidr = "10.0.0.0/24"
    gateway = "10.0.0.254"
    name = "geminitestnet1"
    platform = data.twcc_project.testProject.platform
    project = data.twcc_project.testProject.id
    with_router = true
}

resource "twcc_ike_policy" "ike1" {
    name = "geminitestike1"
    platform = data.twcc_project.testProject.platform
    project = data.twcc_project.testProject.id
}

resource "twcc_ipsec_policy" "ipsec1" {
    name = "geminitestipsec1"
    platform = data.twcc_project.testProject.platform
    project = data.twcc_project.testProject.id
}

resource "twcc_vpn" "vpn1" {
    ike_policy = twcc_ike_policy.ike1.id
    ipsec_policy = twcc_ipsec_policy.ipsec1.id
    name = "geminitestvpn1"
    platform = data.twcc_project.testProject.platform
    private_network = twcc_network.network1.id
}

resource "twcc_vpn_connection" "vc1" {
    peer_address = "10.0.0.254"
    peer_cidrs = ["10.0.0.0/24"]
    platform = data.twcc_project.testProject.platform
    psk = "testgemini"
    vpn = twcc_vpn.vpn1.id
}
