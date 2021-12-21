data "twcc_project" "testProject" {
    name = "ENTxxxxxx"
    platform = "ceph-taichung-default"
}

resource "twcc_s3_key" "key1" {
    name = "geminitestkey"
    platform = data.twcc_project.testProject.platform
    project = data.twcc_project.testProject.id
}
