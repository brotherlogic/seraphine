# Seraphine

Seraphine is a container based CLI app for managing a set of agents designed to work on coding projects.

Essentially once we onboard a github project onto seraphine it monitors for feature proposals, and then works through a feature proposal lifecycle that drives the feature to implementation.

The effective workflow pushes the feature through stages:

1. Triaged

At this point we record the feature within seraphine, we mark the feature proposal as TRIAGED in the prefix of the .md file (within a comment so it doesn't show up in the body). The triaged doc is then pushed to the repo.

1. Implentation Proposed

The next stage is PROPOSED.

We then trigger our planning agent to look at the feature, and write out an implementation plan for the feature. We give the agent free reign to suggest breaking the feature into a number of changes, or allow it to just plan out the whole feature in one go. This initial implementaion plan, and a <feature>-RATIONALE.md is added to give enough context to jump back to a good state. This is then pushed to the repo and we massage the implementation plan through github comments and changes. Once we're satisifed we LGTM the PR and the implementation is added.

So we end up with two files

<feature>-RATIONALE.md
<feature>-IMPLEMENTATION.md

Note that the implementation may be broken down into a number of steps, each of which operates sequentially since we must wait for review before proceeding. 

2. Implementation Phase

The next stage is IMPLEMENTED

We then trigger our coding agent to implement the feature, the coding agent is allowed to delegate the task to other coding agents as it sees fits - these are spun up on demand. Once the code is implemented the PR is pushed and ready for review. The coder adds a file describing what they did to implement the feature:

<feature>-CODE.md

3. Review Phase

There is a gemini reviewer active on the PR as well as the human request. The coding agent from the previous round addresses both the gemini review and any human comments. Once we're satisified with the change the human LGTMs and it's ready to be integrated.

4. Reflection Phase

Once the PR is LGTM'd the reflection agent looks over the created files, moves them to a directory /proposals/<feature>/ and writes a file in the proposal directory summarizing the effort:

<feature>-SUMMARY.md

it also updates the README file to reflect the impact of the changes made.