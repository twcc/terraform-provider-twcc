data "twcc_project" "testProject" {
    name = "ENT108079"
    platform = "k8s-taichung-default"
}

data "twcc_solution" "solution" {
    name = "TensorFlow"
    project = data.twcc_project.testProject.id
}

resource "twcc_container" "container1" {
    extra_property = {
        flavor = "1 GPU + 04 cores + 090GB memory"
        gpfs01-mount-path = "/mnt"
        gpfs02-mount-path = "/tmp"
        image = "tensorflow-19.08-py3:latest"
        replica = "1"
    }

    name = "geminitestcontainer1"
    platform = data.twcc_project.testProject.platform
    project = data.twcc_project.testProject.id
    solution = data.twcc_solution.solution.id
}
