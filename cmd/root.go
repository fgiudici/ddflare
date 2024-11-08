/*
Copyright © 2024 Francesco Giudici <dev@foggy.day>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package cmd

import (
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

// This is called by main.main().
func Execute() {
	app := &cli.App{
		Usage: "update DNS entries via cloudflare APIs",
		Commands: []*cli.Command{
			newGetCommand(),
			newSetCommand(),
			newVersionCommand(),
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
