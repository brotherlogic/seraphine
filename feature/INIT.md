# init

$ seraphine init

The init process does a few things:

1. Asks the server for a list of files that are needed to support
   Seraphine managed projects. These could be github workflows, agent
   skills, a basic devcontainer setup etc. These files are then installed
   in the container and pushed to the default branch

1. Performs updates to the github project to support seraphine. This may
   be either a list of things the user must do to the project, or just
   making changes direcly through the github cli.

1. We set up the seraphine webhook in both the github project and our
   gtihub webhook handler

1. We also create a project proto for the given project and store it
   in our proto store (brotherlogic/pstore).