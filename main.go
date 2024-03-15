// Copyright (C) Andrew Harris 2024 - CyberArk, Inc.

package main

import (
	"context"
	"log"

	// CyberArk API
	provider "github.com/aharriscybr/terraform-provider-cyberarkoss/cyberark"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
)

//go:generate terraform fmt -recursive ./examples/
//go:generate go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs generate -provider-name cyberarkoss

var (
	version string = "1.0.0"
)

func main() {
	opts := providerserver.ServeOpts{
		Address: "registry.terraform.io/aharriscybr/cyberarkoss",
		Debug:   false,
	}

	err := providerserver.Serve(context.Background(), provider.New(version), opts)
	if err != nil {
		log.Fatal(err.Error())
	}
}