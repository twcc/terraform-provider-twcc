data "twcc_project" "testProject" {
    name = "ENTxxxxxx"
    platform = "openstack-taichung-default-2"
}

resource "twcc_secret" "secret1" {
    name = "geminitttestsecret1"
    payload = filebase64("/PATH/server.p12")
    platform = data.twcc_project.testProject.platform
    project = data.twcc_project.testProject.id
}
