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

resource "twcc_volume" "volume1" {
    name = "geminitestvol1"
    platform = data.twcc_project.testProject.platform
    project = data.twcc_project.testProject.id
    size = 1
}

resource "twcc_volume_attachment" "volume_attachment" {
    platform = data.twcc_project.testProject.platform
    server = twcc_vcs.vcs1.servers[0].id
    volume = twcc_volume.volume1.id
}

resource "twcc_volume_snapshot" "vol_snapshot" {
    name = "geminitestsnap1"
    platform = data.twcc_project.testProject.platform
    volume = twcc_volume.volume1.id
}

resource "twcc_volume" "volume2" {
    name = "geminitestvol2"
    platform = data.twcc_project.testProject.platform
    src_snapshot = twcc_volume_snapshot.vol_snapshot.id
}
