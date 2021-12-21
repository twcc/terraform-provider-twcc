data "twcc_project" "testProject" {
    name = "ENTxxxxxx"
    platform = "openstack-taichung-default-2"
}


data "twcc_solution" "solution" {
    name = "ubuntu"
    project = data.twcc_project.testProject.id
}


data "twcc_extra_property" "extra_property" {
    platform = data.twcc_project.testProject.platform
    project = data.twcc_project.testProject.id
    solution = data.twcc_solution.solution.id
}
