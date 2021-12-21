data "twcc_project" "testProject" {
    name = "ENTxxxxxx"
    platform = "openstack-taichung-default-2"
}

data "twcc_solution" "testSolution" {
    name = "F5_WAF"
    project = data.twcc_project.testProject.id
    category = "waf"
}

resource "twcc_waf" "waf1" {
    extra_property = {
        availability-zone = "nova"
        flavor = "08_core_040GB_memory_160GB_disk"
        image = "F5-AWAF-Production"
        password = "password"
        private-network = "default_network"
    }

    name = "geminitestwaf1"
    platform = data.twcc_project.testProject.platform
    project = data.twcc_project.testProject.id
    solution = data.twcc_solution.testSolution.id
}
