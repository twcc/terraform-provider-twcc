data "twcc_project" "testProject" {
    name = "ENT108079"
    platform = "openstack-taichung-default-2"
}

data "twcc_solution" "solution" {
    name = "Ubuntu"
    project = data.twcc_project.testProject.id
}

resource "twcc_vcs" "vcs1" {
    extra_property = {
        availability-zone = "nova"
        flavor = "01_vCPU_002GB_MEM_100GB_HDD"
        floating-ip = "nofloating"
        image = "ubuntu1604"
        keypair = "fatkey"
        private-network = "default_network"
        system-volume-type = "local_disk"
    }

    name = "geminitestvcs1"
    platform = data.twcc_project.testProject.platform
    project = data.twcc_project.testProject.id
    solution = data.twcc_solution.solution.id
}

resource "twcc_auto_scaling_policy" "asp1" {
    meter_name = "cpu_util"
    name = "geminitestasp1"
    platform = data.twcc_project.testProject.platform
    project = data.twcc_project.testProject.id
    scale_max_size = 2
    scaledown_threshold = 10
    scaleup_threshold = 50
}

resource "twcc_auto_scaling_relation" "asr1" {
    platform = data.twcc_project.testProject.platform
    server = twcc_vcs.vcs1.servers[0].id
    auto_scaling_policy = twcc_auto_scaling_policy.asp1.id
}

data "twcc_security_group" "vcs1_sg" {
    platform = data.twcc_project.testProject.platform
    vcs = twcc_vcs.vcs1.id
}

resource "twcc_security_group_rule" "vcs1_sg_rule1" {
    platform = data.twcc_project.testProject.platform
    project = data.twcc_project.testProject.id
    security_group = data.twcc_security_group.vcs1_sg.id
    direction = "egress"
    protocol = "udp"
    remote_ip_prefix = "192.168.0.0/16"
    port_range_min = 8000
    port_range_max = 8010
}

resource "twcc_vcs_image" "snapshot1" {
    name = "geminitestserversnap1"
    os = "Linux"
    os_version = "Ubuntu 16.04"
    platform = data.twcc_project.testProject.platform
    server = twcc_vcs.vcs1.servers[0].id
}
