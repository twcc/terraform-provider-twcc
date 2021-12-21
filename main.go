package main

import (
    "terraform-provider-twcc/twcc"

    "github.com/hashicorp/terraform-plugin-sdk/plugin"
)

func main() {
    plugin.Serve(&plugin.ServeOpts{
        ProviderFunc: twcc.Provider})
}

