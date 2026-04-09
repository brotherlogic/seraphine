# Workflow

I start by onboarding a github project onto Seraphine:

seraphine onboard <github_project>

This checks that (a) the seraphine github user has the appropriate access to the project and then (b) adds all the github workflows and permissions necessary to allow seraphine to run on the project. Effectively this step ensures that seraphine is able to manage the project and create agents in order to effect changes.

In addition it adds the webhook infrastructure necessary to support passback handling in the seraphine control plane. Specifically we add a webhook of the form http://seraphine.brotherlogic-backend.com:9009 which we passback to the seraphine control plane.

We additionally store basic project information, serpahine versioning and project state (ONBOARDING, ONBOARDED) in the cluster.

we also have the option of running:

seraphine projects

Which will list all the currently managed projects and their associated state, and seraphine version.

The seraphine version is used to link the current set of configuration matches the current version of seraphine. So

seraphine upgrade <github_project>

will make the necessary changes to keep the project up to date with the current seraphine github requirements (project settings and github workflows to support the review process.